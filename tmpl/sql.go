package tmpl

// zdt

const SelectZdt = "SELECT * FROM long_short WHERE date >= ? ORDER BY date DESC LIMIT ?"

// zdt end

// concept

const SelectConceptStockByConceptId = "SELECT * FROM concept_stock WHERE concept_id = ? ORDER BY updated_at"

const SelectConceptStockViewBody = `SELECT
s.stock_code,
s.stock_name,
s.updated_at,
s.description,
c.id AS concept_id,
c.NAME AS concept_name,
c.plate_id AS concept_plate_id,
c.define AS concept_define,
c.updated_at AS concept_updated_at`

const SelectConceptStockViewFrom = `
FROM
concept_stock AS s
INNER JOIN concept_concept AS c ON c.id = s.concept_id`

const SelectAllConceptStockView = SelectConceptStockViewBody + SelectConceptStockViewFrom

const SelectConceptStockViewByStockConceptKw = SelectConceptStockViewBody + SelectConceptStockViewFrom + `
WHERE
	MATCH ( s.stock_code, s.stock_name ) against ( ? ) 
	AND (
		MATCH ( c.NAME, c.define ) against ( ? ) 
	OR MATCH ( s.description ) against ( ? )) 
ORDER BY
	(MATCH ( s.description ) against ( ? ) + MATCH ( c.NAME, c.define ) against ( ? )) DESC,
	s.updated_at desc
LIMIT ?`

const SelectConceptStockViewByStockKw = SelectConceptStockViewBody + SelectConceptStockViewFrom + `
WHERE MATCH ( s.stock_code, s.stock_name ) against ( ? )
order by MATCH ( s.stock_code, s.stock_name ) against ( ? ) desc,
s.updated_at desc
LIMIT ?`

const SelectConceptStockViewByConceptKw = SelectConceptStockViewBody + SelectConceptStockViewFrom + `
where 
	MATCH ( c.NAME, c.define ) against ( ? ) 
	OR MATCH ( s.description ) against ( ? )
	ORDER BY
	MATCH ( s.description ) against ( ? ) DESC,
	MATCH ( c.NAME, c.define ) against ( ? ) desc, 
	s.updated_at desc
LIMIT ?`

const SelectConceptStockViewByUpdateAtDesc = SelectConceptStockViewBody + SelectConceptStockViewFrom + " order by s.updated_at desc limit ?"

const SelectConceptByName = "select * from concept_concept where name=IFNULL(?, name) order by updated_at desc limit ?"

const SelectAllConceptStock = "SELECT * FROM concept_stock"

const UpdateConceptStock = "UPDATE concept_stock SET stock_name=?, description=?, updated_at=? WHERE stock_code=? and concept_id=?"

const InsertConceptStock = "insert INTO concept_stock VALUES (?,?,?,?,?)"

const InsertConcept = "insert INTO concept_concept VALUES (?,?,?,?,?)"

const UpdateConcept = "UPDATE concept_concept SET name = ?, plate_id = ?, define = ?, updated_at = ? WHERE id = ?"

// concept end

// realtime start

const SelectRealtimeByUserId = "select * from realtime_archive where user_id=? order by seq desc limit ?"

const DeleteRealtimeByUserIdSeq = "delete from realtime_archive where user_id=? and seq=?"

// realtime end
