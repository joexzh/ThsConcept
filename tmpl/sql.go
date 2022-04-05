package tmpl

// zdt

const SelectZdt = "SELECT * FROM long_short WHERE date >= ? ORDER BY date DESC LIMIT ?"

// concept

const SelectConceptById = "SELECT * FROM concept_concept WHERE id=?"

const SelectAllConceptOrderById = "SELECT * FROM concept_concept order by id"

const SelectConceptByName = "select * from concept_concept where name=IFNULL(?, name) order by updated_at desc limit ?"

const InsertConcept = "insert INTO concept_concept VALUES (?,?,?,?,?)"

const UpdateConcept = "UPDATE concept_concept SET name = ?, plate_id = ?, define = ?, updated_at = ? WHERE id = ?"

// concept_stock

const SelectConceptStockByConceptIdOrderByUpdatedAt = "SELECT * FROM concept_stock WHERE concept_id = ? ORDER BY updated_at"

const SelectConceptStockByConceptIdOrderByCode = "SELECT * FROM concept_stock WHERE concept_id = ? ORDER BY stock_code"

const SelectConceptStockByStockCodeConceptId = "SELECT * FROM concept_stock where stock_code=? and concept_id=?"

const UpdateConceptStock = "UPDATE concept_stock SET stock_name=?, description=?, updated_at=? WHERE stock_code=? and concept_id=?"

const InsertConceptStock = "insert INTO concept_stock VALUES (?,?,?,?,?)"

const DeleteConceptStock = "DELETE FROM concept_stock WHERE stock_code=? and concept_id=?"

// concept_concept join concept_stock

const SelectConceptJoinStockByConceptIdStockCode = `SELECT
s.stock_code,
s.stock_name,
s.updated_at,
s.description,
c.id AS concept_id,
c.NAME AS concept_name,
c.plate_id AS concept_plate_id,
c.define AS concept_define,
c.updated_at AS concept_updated_at 
FROM
concept_stock AS s
INNER JOIN concept_concept AS c ON c.id = s.concept_id
WHERE s.stock_code=? and c.id=?`

// const concept_stock_ft

const SelectFromConceptStockFt = `SELECT
	stock_code,
	stock_name,
	updated_at,
	description,
	concept_id,
	concept_name,
	concept_plate_id,
	concept_define,
	concept_updated_at
FROM concept_stock_ft`

const SelectConceptStockFtByStockConceptKw = SelectFromConceptStockFt + `
WHERE
	(
	stock_code = ?
	OR stock_name = ?)
AND MATCH ( description, concept_name, concept_define ) against ( ? ) 
LIMIT ?`

const SelectConceptStockFtByStockKw = SelectFromConceptStockFt + `
WHERE
	stock_code = ?
	OR stock_name = ?
LIMIT ?`

const SelectConceptStockFtByConceptKw = SelectFromConceptStockFt + `
WHERE
	MATCH ( description, concept_name, concept_define ) against ( ? ) 
LIMIT ?`

const SelectConceptStockFtByUpdateAtDesc = SelectFromConceptStockFt + " order by updated_at desc limit ?"

const UpdateConceptInConceptStockFt = `update concept_stock_ft set concept_name=?,concept_plate_id=?,concept_define=?,concept_updated_at=?
where concept_id = ?`

const UpdateStockInConceptStockFt = `update concept_stock_ft set stock_name=?,description=?,updated_at=?
where stock_code = ? and concept_id = ?`

const DeleteConceptStockFtByConceptId = "delete from concept_stock_ft where concept_id=?"

const DeleteConceptStockFtByStockCodeConceptId = "delete from concept_stock_ft where stock_code=? and concept_id=?"

// concept_stock_ft_sync

const SelectAllConceptStockFtSync = "SELECT * FROM concept_stock_ft_sync"

const DeleteConceptStockFtSyncById = "delete from concept_stock_ft_sync where id=?"

// realtime start

const SelectRealtimeByUserId = "select * from realtime_archive where user_id=? order by seq desc limit ?"

const DeleteRealtimeByUserIdSeq = "delete from realtime_archive where user_id=? and seq=?"
