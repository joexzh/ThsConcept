package tmpl

// zdt

const SelectZdt = "SELECT * FROM long_short WHERE date >= ? ORDER BY date DESC LIMIT ?"

// zdt end

// concept

const SelectConceptStockByConceptId = "SELECT * FROM concept_stock WHERE concept_id = ? ORDER BY updated_at"

const SelectFromConceptStockView = `SELECT
	stock_code,
	stock_name,
	updated_at,
	description,
	concept_id,
	concept_name,
	concept_plate_id,
	concept_define,
	concept_updated_at
FROM concept_stock_full`

const SelectConceptStockViewByStockConceptKw = SelectFromConceptStockView + `
WHERE
	(
	stock_code = ?
	OR stock_name = ?)
AND MATCH ( description, concept_name, concept_define ) against ( ? ) 
LIMIT ?`

const SelectConceptStockViewByStockKw = SelectFromConceptStockView + `
WHERE
	stock_code = ?
	OR stock_name = ?
LIMIT ?`

const SelectConceptStockViewByConceptKw = SelectFromConceptStockView + `
WHERE
	MATCH ( description, concept_name, concept_define ) against ( ? ) 
LIMIT ?`

const SelectConceptStockViewByUpdateAtDesc = SelectFromConceptStockView + " order by updated_at desc limit ?"

const SelectConceptByName = "select * from concept_concept where name=IFNULL(?, name) order by updated_at desc limit ?"

const SelectAllConceptStock = "SELECT * FROM concept_stock"

const UpdateConceptStock = "UPDATE concept_stock SET stock_name=?, description=?, updated_at=? WHERE stock_code=? and concept_id=?"

const InsertConceptStock = "insert INTO concept_stock VALUES (?,?,?,?,?)"

const InsertConcept = "insert INTO concept_concept VALUES (?,?,?,?,?)"

const UpdateConcept = "UPDATE concept_concept SET name = ?, plate_id = ?, define = ?, updated_at = ? WHERE id = ?"

const SelectDistinctConceptStockByKw = `SELECT DISTINCT
stock_code,
stock_name 
FROM
concept_stock 
WHERE
MATCH ( stock_code, stock_name ) against ( ? IN boolean MODE ) 
AND stock_name <> ''`

// concept end

// realtime start

const SelectRealtimeByUserId = "select * from realtime_archive where user_id=? order by seq desc limit ?"

const DeleteRealtimeByUserIdSeq = "delete from realtime_archive where user_id=? and seq=?"

// realtime end
