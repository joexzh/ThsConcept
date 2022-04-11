package repos

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/joexzh/dbh"
	"golang.org/x/exp/slices"

	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/db"
	"github.com/joexzh/ThsConcept/model"
	"github.com/joexzh/ThsConcept/tmpl"
	"github.com/pkg/errors"
)

type StockMarketRepo struct {
	DB   *sql.DB
	Name string
}

func NewStockMarketRepo(db *sql.DB) *StockMarketRepo {
	return &StockMarketRepo{DB: db, Name: "StockMarketRepo"}
}

func (repo *StockMarketRepo) ZdtListDesc(ctx context.Context, start time.Time, limit int) ([]*model.ZDTHistory, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	list, err := dbh.QueryContext[*model.ZDTHistory](repo.DB, ctx, tmpl.SelectZdt,
		start, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	// reverse
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	return list, nil
}

func (repo *StockMarketRepo) InsertZdtList(ctx context.Context, list []*model.ZDTHistory) (int64, error) {
	if len(list) < 1 {
		return 0, nil
	}

	conn, err := repo.DB.Conn(ctx)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}
	defer conn.Close()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}
	defer tx.Rollback()

	if err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}

	total, err := dbh.BulkInsertContext(tx, ctx, 1000, list...)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}
	if err = tx.Commit(); err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}
	return total, nil
}

func scsToConceptMap(scs []*model.ConceptStockFt) map[string]*model.Concept {
	conceptMap := make(map[string]*model.Concept)
	for _, sc := range scs {
		concept, ok := conceptMap[sc.ConceptId]
		if !ok {
			concept = &model.Concept{
				Id:        sc.ConceptId,
				Name:      sc.ConceptName,
				PlateId:   sc.ConceptPlateId,
				Define:    sc.ConceptDefine,
				UpdatedAt: sc.ConceptUpdatedAt,
				Stocks:    make([]*model.ConceptStock, 0),
			}
			conceptMap[sc.ConceptId] = concept
		}
		stock := &model.ConceptStock{
			StockCode:   sc.StockCode,
			StockName:   sc.StockName,
			Description: sc.Description,
			UpdatedAt:   sc.UpdatedAt,
			ConceptId:   sc.ConceptId,
		}
		concept.Stocks = append(concept.Stocks, stock)
	}
	return conceptMap
}

type UpdateConceptResult struct {
	ConceptConceptInserted int64
	ConceptConceptUpdated  int64
	ConceptConceptDeleted  int64
	ConceptStockInserted   int64
	ConceptStockUpdated    int64
	ConceptStockDeleted    int64
}

// compare newList and oldList, inject insert, update and delete
// Both list should be sorted.
func injectUpdateWhenCompareList[S []E, E model.ComparableConcept[E]](newList S, oldList S,
	insert func(E) error, update func(E) error, delete func(E) error,
	compareSubList func(newItem, oldItem E) error) error {

	for i, j := 0, 0; i < len(newList) && j < len(oldList); {
		if j == len(oldList) {
			for k := range newList[i:] {
				if err := insert(newList[k]); err != nil {
					return err
				}
			}
			break
		}
		if i == len(newList) {
			for k := range oldList[j:] {
				if err := delete(oldList[k]); err != nil {
					return err
				}
			}
			break
		}

		if newList[i].GetId() == oldList[j].GetId() {
			if !newList[i].Cmp(oldList[j]) {
				// concept changed and add to update list
				if err := update(newList[i]); err != nil {
					return err
				}
			}
			if compareSubList != nil {
				if err := compareSubList(newList[i], oldList[j]); err != nil {
					return err
				}
			}
			i++
			j++
		} else if oldList[j].GetId() < newList[i].GetId() {
			// cannot find new newcs, add to delete list
			if err := delete(oldList[j]); err != nil {
				return err
			}
			j++
		} else {
			// cannot find old concept, add to insert list
			if err := insert(newList[i]); err != nil {
				return err
			}
			i++
		}
	}
	return nil
}

func SearchConceptById(id string, cs []*model.Concept) (*model.Concept, bool) {
	i, ok := slices.BinarySearchFunc(cs, &model.Concept{Id: id}, func(a, b *model.Concept) int {
		if a.Id < b.Id {
			return -1
		} else if a.Id > b.Id {
			return 1
		}
		return 0
	})
	if !ok {
		return nil, ok
	}
	return cs[i], ok
}

func (repo *StockMarketRepo) UpdateConcept(ctx context.Context, newcs ...*model.Concept) (UpdateConceptResult, error) {
	result := UpdateConceptResult{}

	conn, err := repo.DB.Conn(ctx)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	defer conn.Close()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	defer tx.Rollback()

	// change cache
	insertedConcepts := make([]*model.Concept, 0)
	updatedConcepts := make([]*model.Concept, 0)
	deletedConcepts := make([]*model.Concept, 0)
	insertedConceptStocks := make([]*model.ConceptStock, 0)
	updatedConceptStocks := make([]*model.ConceptStock, 0)
	deletedConceptStocks := make([]*model.ConceptStock, 0)

	// query old concepts for update
	oldcs, err := dbh.QueryContext[*model.Concept](tx, ctx, tmpl.SelectAllConceptOrderById+" for update")
	if err != nil {
		return result, errors.Wrap(err, repo.Name+"\n"+
			`dbh.QueryContext[*model.Concept](tx, ctx, tmpl.SelectAllConceptOrderById+" for update")`)
	}
	sort.Sort(model.ConceptSortById(newcs))

	appendInsertConceptList := func(concept *model.Concept) error {
		insertedConcepts = append(insertedConcepts, concept)
		insertedConceptStocks = append(insertedConceptStocks, concept.Stocks...)
		return nil
	}
	appendUpdateConcept := func(concept *model.Concept) error {
		updatedConcepts = append(updatedConcepts, concept)
		return nil
	}
	appendDeleteConcept := func(concept *model.Concept) error {
		deletedConcepts = append(deletedConcepts, concept)
		deletedConceptStocks = append(deletedConceptStocks, concept.Stocks...)
		return nil
	}
	conceptStockStmt, err := tx.PrepareContext(ctx, tmpl.SelectConceptStockByConceptIdOrderByCode+" for update")
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	err = injectUpdateWhenCompareList(newcs, oldcs, appendInsertConceptList, appendUpdateConcept, appendDeleteConcept,
		func(newConcept, oldConcept *model.Concept) error {
			// query stocks by conceptId from db
			rows, err := conceptStockStmt.QueryContext(ctx, oldConcept.Id)
			if err != nil {
				return err
			}
			var oldStocks []*model.ConceptStock
			err = dbh.ScanList(rows, &oldStocks)
			rows.Close()
			if err != nil {
				return err
			}

			sort.Sort(model.ConceptStockSortByCode(newConcept.Stocks))

			// compare stocks
			return injectUpdateWhenCompareList(newConcept.Stocks, oldStocks,
				func(stock *model.ConceptStock) error {
					insertedConceptStocks = append(insertedConceptStocks, stock)
					return nil
				},
				func(stock *model.ConceptStock) error {
					updatedConceptStocks = append(updatedConceptStocks, stock)
					return nil
				},
				func(stock *model.ConceptStock) error {
					deletedConceptStocks = append(deletedConceptStocks, stock)
					return nil
				},
				nil)
		},
	)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}

	commands := make([]*model.ConceptFtCommand, 0, len(insertedConcepts)+len(updatedConcepts)+len(deletedConcepts)+
		len(insertedConceptStocks)+len(updatedConceptStocks)+len(deletedConceptStocks))

	// exec changes
	// insert concept_concept
	if len(insertedConcepts) > 0 {
		r, err := dbh.BulkInsertContext(tx, ctx, 1000, insertedConcepts...)
		if err != nil {
			return result, errors.Wrap(err, repo.Name)
		}
		result.ConceptConceptInserted += r

		for i := range insertedConcepts {
			commands = append(commands, &model.ConceptFtCommand{
				Command: model.InsertConcept,
				Obj:     model.ConceptStockFt{ConceptPlateId: i}, // just for remove compile err for i no being used
			})
		}

	}
	// update concept_concept
	if len(updatedConcepts) > 0 {
		updateConceptStmt, err := tx.PrepareContext(ctx, tmpl.UpdateConcept)
		if err != nil {
			return result, errors.Wrap(err, repo.Name)
		}
		for i := range updatedConcepts {
			ra, err := updateConceptStmt.ExecContext(ctx, updatedConcepts[i].Name, updatedConcepts[i].PlateId,
				updatedConcepts[i].Define, updatedConcepts[i].UpdatedAt, updatedConcepts[i].Id)
			if err != nil {
				return result, errors.Wrap(err, repo.Name)
			}
			r, _ := ra.RowsAffected()
			result.ConceptConceptUpdated += r

			commands = append(commands, &model.ConceptFtCommand{
				Command: model.UpdateConcept,
				Obj: model.ConceptStockFt{
					ConceptId:        insertedConcepts[i].Id,
					ConceptName:      insertedConcepts[i].Name,
					ConceptPlateId:   insertedConcepts[i].PlateId,
					ConceptDefine:    insertedConcepts[i].Define,
					ConceptUpdatedAt: insertedConcepts[i].UpdatedAt,
				},
			})
		}
	}
	// delete concept_concept
	if len(deletedConcepts) > 0 {
		ids := make([]string, len(deletedConcepts))
		for i := range deletedConcepts {
			ids[i] = deletedConcepts[i].Id

			commands = append(commands, &model.ConceptFtCommand{
				Command: model.DeleteConcept,
				Obj: model.ConceptStockFt{
					ConceptId: insertedConcepts[i].Id,
				},
			})
		}
		listSql, vals := db.ArgList(ids)
		ra, err := tx.ExecContext(ctx, fmt.Sprintf("delete from concept_concept where id in %s", listSql), vals...)
		if err != nil {
			return result, errors.Wrap(err, repo.Name+"\n"+
				`tx.ExecContext(ctx, fmt.Sprintf("delete from concept_concept where id in %s", listSql), vals...)`)
		}
		r, _ := ra.RowsAffected()
		result.ConceptConceptDeleted += r
	}
	// insert concept_stock
	if len(insertedConceptStocks) > 0 {
		r, err := dbh.BulkInsertContext(tx, ctx, 1000, insertedConceptStocks...)
		if err != nil {
			return result, errors.Wrap(err, repo.Name+"\n"+"dbh.BulkInsertContext(tx, ctx, 1000, insertedConceptStocks...)")
		}
		result.ConceptStockInserted += r

		for i := range insertedConceptStocks {
			concept, ok := SearchConceptById(insertedConceptStocks[i].ConceptId, newcs)
			if ok {
				commands = append(commands, &model.ConceptFtCommand{
					Command: model.InsertConceptStock,
					Obj: model.ConceptStockFt{
						StockCode:        insertedConceptStocks[i].StockCode,
						StockName:        insertedConceptStocks[i].StockName,
						UpdatedAt:        insertedConceptStocks[i].UpdatedAt,
						Description:      insertedConceptStocks[i].Description,
						ConceptId:        concept.Id,
						ConceptName:      concept.Name,
						ConceptPlateId:   concept.PlateId,
						ConceptDefine:    concept.Define,
						ConceptUpdatedAt: concept.UpdatedAt,
					},
				})
			}

		}
	}
	// update concept_stock
	if len(updatedConceptStocks) > 0 {
		updateConceptStockStmt, err := tx.PrepareContext(ctx, tmpl.UpdateConceptStock)
		if err != nil {
			return result, errors.Wrap(err, repo.Name)
		}
		for i := range updatedConceptStocks {
			ra, err := updateConceptStockStmt.ExecContext(ctx, updatedConceptStocks[i].StockName,
				updatedConceptStocks[i].Description, updatedConceptStocks[i].UpdatedAt,
				updatedConceptStocks[i].StockCode, updatedConceptStocks[i].ConceptId)
			if err != nil {
				return result, errors.Wrap(err, repo.Name+"\n"+
					"updateConceptStockStmt.ExecContext")
			}
			r, _ := ra.RowsAffected()
			result.ConceptStockUpdated += r

			commands = append(commands, &model.ConceptFtCommand{
				Command: model.UpdateConceptStock,
				Obj: model.ConceptStockFt{
					StockCode:   updatedConceptStocks[i].StockCode,
					StockName:   updatedConceptStocks[i].StockName,
					UpdatedAt:   updatedConceptStocks[i].UpdatedAt,
					Description: updatedConceptStocks[i].Description,
					ConceptId:   updatedConceptStocks[i].ConceptId,
				},
			})

		}
	}
	// delete concept_stock
	if len(deletedConceptStocks) > 0 {
		deleteConceptStockStmt, err := tx.PrepareContext(ctx, tmpl.DeleteConceptStock)
		if err != nil {
			return result, errors.Wrap(err, repo.Name)
		}
		for i := range deletedConceptStocks {
			ra, err := deleteConceptStockStmt.ExecContext(ctx, deletedConceptStocks[i].StockCode, deletedConceptStocks[i].ConceptId)
			if err != nil {
				return result, errors.Wrap(err, repo.Name+"\n"+
					`deleteConceptStockStmt.ExecContext(ctx, deletedConceptStocks[i].StockCode, deletedConceptStocks[i].ConceptId)`)
			}
			r, _ := ra.RowsAffected()
			result.ConceptStockDeleted += r

			commands = append(commands, &model.ConceptFtCommand{
				Command: model.DeleteConceptStock,
				Obj: model.ConceptStockFt{
					StockCode: deletedConceptStocks[i].StockCode,
					ConceptId: deletedConceptStocks[i].ConceptId,
				},
			})
		}
	}

	// save udpate command for sync
	if err = repo.saveCommands(tx, ctx, commands); err != nil {
		return result, errors.Wrap(err, repo.Name+"\n"+"saveCommands")
	}

	if err = tx.Commit(); err != nil {
		return result, errors.Wrap(err, repo.Name)
	}

	return result, nil
}

func (repo *StockMarketRepo) saveCommands(db dbh.DbInterface, ctx context.Context, commands []*model.ConceptFtCommand) error {
	_, err := dbh.BulkInsertContext(db, ctx, 1000, commands...)
	return err
}

func (repo *StockMarketRepo) ConceptStockFtSync(ctx context.Context) error {
	conn, err := repo.DB.Conn(ctx)
	if err != nil {
		return errors.Wrap(err, repo.Name)
	}
	defer conn.Close()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, repo.Name)
	}
	defer tx.Rollback()

	cs, err := dbh.QueryContext[*model.ConceptFtCommand](tx, ctx, tmpl.SelectAllConceptStockFtSyncOrderById)
	if err != nil {
		return errors.Wrap(err, repo.Name)
	}
	for i := range cs {
		switch cs[i].Command {
		case model.InsertConcept:
			// do nothing
		case model.UpdateConcept:
			ft := &cs[i].Obj
			_, err = tx.ExecContext(ctx, tmpl.UpdateConceptInConceptStockFt,
				ft.ConceptName, ft.ConceptPlateId, ft.ConceptDefine, ft.ConceptUpdatedAt, ft.ConceptId)
			if err != nil {
				return errors.Wrap(err, repo.Name+"\n"+"tmpl.UpdateConceptInConceptStockFt")
			}
		case model.DeleteConcept:
			_, err = tx.ExecContext(ctx, tmpl.DeleteConceptStockFtByConceptId, cs[i].Obj.ConceptId)
			if err != nil {
				return errors.Wrap(err, repo.Name+"\n"+"tmpl.DeleteConceptStockFtByConceptId")
			}
		case model.InsertConceptStock:
			_, err = dbh.InsertContext(tx, ctx, &cs[i].Obj)
			if err != nil {
				return errors.Wrap(err, repo.Name+"\n"+"dbh.InsertContext(tx, ctx, ft)")
			}
		case model.UpdateConceptStock:
			ft := &cs[i].Obj
			_, err = tx.ExecContext(ctx, tmpl.UpdateStockInConceptStockFt, ft.StockName, ft.Description, ft.UpdatedAt,
				ft.StockCode, ft.ConceptId)
			if err != nil {
				return errors.Wrap(err, repo.Name+"\n"+"tmpl.UpdateStockInConceptStockFt")
			}
		case model.DeleteConceptStock:
			_, err = tx.ExecContext(ctx, tmpl.DeleteConceptStockFtByStockCodeConceptId, cs[i].Obj.StockCode, cs[i].Obj.ConceptId)
			if err != nil {
				return errors.Wrap(err, repo.Name+"\n"+"tmpl.DeleteConceptStockFtByStockCodeConceptId")
			}
		}

		_, err = tx.ExecContext(ctx, tmpl.DeleteConceptStockFtSyncById, cs[i].Id)
		if err != nil {
			return errors.Wrap(err, repo.Name+"\n"+"tmpl.DeleteConceptStockFtSyncById")
		}
	}

	return tx.Commit()
}

const conceptLimit = 500

func (repo *StockMarketRepo) QueryConceptFtStockByKw(ctx context.Context, stockKw string, conceptKw string, limit int) (
	[]*model.ConceptStockFt, error) {

	switch {
	case stockKw != "" && conceptKw != "":
		return repo.QueryConceptStockFtByStockConceptKw(ctx, stockKw, conceptKw, limit)
	case stockKw != "":
		return repo.QeuryConceptStockFtByStockKw(ctx, stockKw, limit)
	case conceptKw != "":
		return repo.QueryConceptStockFtByConceptKw(ctx, conceptKw, limit)
	default:
		return repo.QueryConceptStockFtSortUpdatedDesc(ctx, limit)
	}
}

func (repo *StockMarketRepo) QueryConceptStockFtByStockConceptKw(ctx context.Context, stockKw string, conceptKw string, limit int) (
	[]*model.ConceptStockFt, error) {

	if limit < 1 || limit > conceptLimit {
		limit = conceptLimit
	}
	scs := make([]*model.ConceptStockFt, 0)
	scs, err := dbh.QueryContext[*model.ConceptStockFt](repo.DB, ctx, tmpl.SelectConceptStockFtByStockConceptKw,
		stockKw, stockKw, conceptKw, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QeuryConceptStockFtByStockKw(ctx context.Context, stockKw string, limit int) ([]*model.ConceptStockFt, error) {
	if limit < 1 || limit > conceptLimit {
		limit = conceptLimit
	}
	scs := make([]*model.ConceptStockFt, 0)
	scs, err := dbh.QueryContext[*model.ConceptStockFt](repo.DB, ctx, tmpl.SelectConceptStockFtByStockKw,
		stockKw, stockKw, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryConceptStockFtByConceptKw(ctx context.Context, conceptKw string, limit int) ([]*model.ConceptStockFt, error) {
	if limit < 1 || limit > conceptLimit {
		limit = conceptLimit
	}
	scs := make([]*model.ConceptStockFt, 0)
	scs, err := dbh.QueryContext[*model.ConceptStockFt](repo.DB, ctx, tmpl.SelectConceptStockFtByConceptKw,
		conceptKw, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryConceptStockFtSortUpdatedDesc(ctx context.Context, limit int) ([]*model.ConceptStockFt, error) {
	if limit < 1 || limit > conceptLimit {
		limit = conceptLimit
	}
	scs := make([]*model.ConceptStockFt, 0)
	scs, err := dbh.QueryContext[*model.ConceptStockFt](repo.DB, ctx, tmpl.SelectConceptStockFtOrderByUpdateAtDesc,
		limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryConcepts(ctx context.Context, conceptKw string, limit int, includeStock bool) ([]*model.Concept, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	var conceptVal interface{}
	if conceptKw != "" {
		conceptVal = conceptKw
	}
	concepts, err := dbh.QueryContext[*model.Concept](repo.DB, ctx,
		tmpl.SelectConceptByName,
		conceptVal, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	if includeStock {
		for _, concept := range concepts {
			stocks, err := repo.QueryConceptStockByConceptId(ctx, concept.Id)
			if err != nil {
				return nil, errors.Wrap(err, repo.Name)
			}
			concept.Stocks = stocks
		}
	}
	return concepts, nil
}

func (repo *StockMarketRepo) QueryConceptStockByConceptId(ctx context.Context, conceptId string) ([]*model.ConceptStock, error) {
	scs, err := dbh.QueryContext[*model.ConceptStock](repo.DB, ctx, tmpl.SelectConceptStockByConceptIdOrderByUpdatedAt,
		conceptId)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryRealtimeArchive(ctx context.Context, userId int, limit int) ([]*model.RealtimeMessage, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	messages, err := dbh.QueryContext[*model.RealtimeMessage](repo.DB, ctx,
		tmpl.SelectRealtimeByUserId,
		userId, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return messages, nil
}

func (repo *StockMarketRepo) SaveRealtimeArchive(ctx context.Context, message *model.RealtimeMessage) (int64, error) {
	return dbh.InsertContext(repo.DB, ctx, message)
}

func (repo *StockMarketRepo) DeleteRealtimeArchive(ctx context.Context, userId int, seq string) (int64, error) {
	ret, err := repo.DB.ExecContext(ctx, tmpl.DeleteRealtimeByUserIdSeq, userId, seq)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}
	ra, _ := ret.RowsAffected()
	return ra, nil
}

func (repo *StockMarketRepo) filterConceptLinesToInsert(db dbh.DbInterface, ctx context.Context, plateId string,
	newLines []*model.ConceptLine) ([]*model.ConceptLine, error) {
	var lastDate time.Time
	err := db.QueryRowContext(ctx, "select date from concept_line where plate_id=? order by date desc limit 1", plateId).
		Scan(&lastDate)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	sort.Sort(model.ConceptLineSortedByDate(newLines))
	var filteredLines []*model.ConceptLine

	for i := len(newLines) - 1; i >= 0; i-- {
		if newLines[i].Date.After(lastDate) {
			filteredLines = append(filteredLines, newLines[i])
		}
	}
	return filteredLines, nil
}

func (repo *StockMarketRepo) FilterConceptLineMap(ctx context.Context, m map[string][]*model.ConceptLine) ([]*model.ConceptLine, error) {
	if len(m) == 0 {
		return nil, nil
	}

	conn, err := repo.DB.Conn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name+".FilterConceptLineMap")
	}
	defer conn.Close()

	var lines []*model.ConceptLine
	for k := range m {
		l, err := repo.filterConceptLinesToInsert(conn, ctx, k, m[k])
		if err != nil {
			return nil, errors.Wrap(err, repo.Name+".FilterConceptLineMap")
		}
		lines = append(lines, l...)
	}
	return lines, nil
}

func (repo *StockMarketRepo) InsertConceptLines(ctx context.Context, lines []*model.ConceptLine) (int64, error) {
	if len(lines) == 0 {
		return 0, nil
	}

	conn, err := repo.DB.Conn(ctx)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name+".InsertConceptLines")
	}
	defer conn.Close()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name+".InsertConceptLines")
	}
	defer tx.Rollback()

	r, err := dbh.BulkInsertContext(tx, ctx, 1000, lines...)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name+".InsertConceptLines")
	}
	return r, tx.Commit()
}

func (repo *StockMarketRepo) QueryConceptLineByDate(ctx context.Context, date time.Time) ([]*model.ConceptLine, error) {
	return dbh.QueryContext[*model.ConceptLine](repo.DB, ctx, `select * from concept_line
where date=?`, date)
}

func (repo *StockMarketRepo) QueryAllPlateIds(ctx context.Context) ([]string, error) {
	rows, err := repo.DB.QueryContext(ctx, "select plate_id from concept_concept")
	if err != nil {
		return nil, errors.Wrap(err, repo.Name+".QueryAllPlateIds")
	}
	defer rows.Close()

	pIds := make([]string, 0)
	for rows.Next() {
		var pId string
		if err = rows.Scan(&pId); err != nil {
			return nil, errors.Wrap(err, repo.Name+".QueryAllPlateIds")
		}
		pIds = append(pIds, pId)
	}
	return pIds, nil
}

func (repo *StockMarketRepo) ViewConceptLineByDateRange(ctx context.Context, startDate time.Time, endDate time.Time) (
	[]*model.ConceptLineDatePctChgOrderedView, error) {
	const limit = 10
	if endDate.Sub(startDate) >= limit*24*time.Hour {
		switch {
		case startDate.IsZero() && endDate.IsZero():
			now := time.Now()
			endDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, config.ChinaLoc())
			startDate = endDate.Add(-limit * 24 * time.Hour)
		case startDate.IsZero():
			startDate = endDate.Add(-limit * 24 * time.Hour)
		case endDate.IsZero():
			endDate = startDate.Add(limit * 24 * time.Hour)
		default:
			startDate = endDate.Add(-limit * 24 * time.Hour)
		}
	}

	conn, err := repo.DB.Conn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name+".ViewConceptLineByDateRange")
	}
	defer conn.Close()

	view := make([]*model.ConceptLineDatePctChgOrderedView, 0)

	for startDate.Before(endDate) || startDate.Equal(endDate) {
		lines, err := dbh.QueryContext[*model.ConceptLineWithName](conn, ctx,
			`SELECT
			l.*,
			c.NAME AS concept_name 
		FROM
			concept_line AS l
			INNER JOIN concept_concept AS c ON l.plate_id = c.plate_id 
		WHERE
			date = ? 
		ORDER BY
			pct_chg DESC 
			LIMIT ?`,
			startDate, limit*2)
		if err != nil {
			return nil, errors.Wrap(err, repo.Name+".ViewConceptLineByDateRange")
		}
		if len(lines) > 0 {
			view = append(view, &model.ConceptLineDatePctChgOrderedView{
				Date: startDate, Lines: lines})
		}

		startDate = startDate.Add(24 * time.Hour)
	}

	return view, nil
}
