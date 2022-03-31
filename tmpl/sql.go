package tmpl

// zdt

const SelectZdt = "SELECT * FROM long_short WHERE date >= ? ORDER BY date DESC LIMIT ?"

// zdt end

// concept

const SelectAllSc = `SELECT
s.CODE AS stock_code,
s.NAME AS stock_name,
sc.updated_at,
sc.description,
c.id AS concept_id,
c.NAME AS concept_name,
c.plate_id AS concept_plate_id,
c.define AS concept_define,
c.updated_at AS concept_updated_at 
FROM
concept_stock AS s
INNER JOIN concept_stock_concept AS sc ON sc.stock_code = s.
CODE INNER JOIN concept_concept AS c ON c.id = sc.concept_id
for update`

const SelectScByName = `SELECT
s.CODE AS stock_code,
s.NAME AS stock_name,
sc.updated_at,
sc.description,
c.id AS concept_id,
c.NAME AS concept_name,
c.plate_id AS concept_plate_id,
c.define AS concept_define,
c.updated_at AS concept_updated_at 
FROM
concept_stock AS s
INNER JOIN concept_stock_concept AS sc ON sc.stock_code = s.
CODE INNER JOIN concept_concept AS c ON c.id = sc.concept_id
WHERE
(
	s.CODE = IFNULL(?, s.code)
	OR s.NAME = IFNULL(?, s.name)
) 
and c.NAME = IFNULL(?, c.name) 
ORDER BY
sc.updated_at DESC 
LIMIT ?`

const SelectScByConceptId = `SELECT
s.CODE AS stock_code,
s.NAME AS stock_name,
sc.updated_at,
sc.description,
c.id AS concept_id,
c.NAME AS concept_name,
c.plate_id AS concept_plate_id,
c.define AS concept_define,
c.updated_at AS concept_updated_at 
FROM
concept_stock AS s
INNER JOIN concept_stock_concept AS sc ON sc.stock_code = s.
CODE INNER JOIN concept_concept AS c ON c.id = sc.concept_id
where c.id=?
order by sc.updated_at`

const SelectConceptByName = "select * from concept_concept where name=IFNULL(?, name) order by updated_at desc limit ?"

const SelectAllConceptStock = "SELECT * FROM concept_stock"

const UpdateConceptStock = "UPDATE concept_stock SET name = ? WHERE code = ?"

const InsertConceptStock = "insert INTO concept_stock VALUES (?,?)"

const InsertConcept = "insert INTO concept_concept VALUES (?,?,?,?,?)"

const UpdateConcept = "UPDATE concept_concept SET name = ?, plate_id = ?, define = ?, updated_at = ? WHERE id = ?"

const InsertConceptStockConcept = "insert INTO concept_stock_concept VALUES (?,?,?,?)"

const UpdateConceptStockConcept = "UPDATE concept_stock_concept SET description = ?, updated_at = ? WHERE stock_code = ? AND concept_id = ?"

// concept end

// realtime start

const SelectRealtimeByUserId = "select * from realtime_archive where user_id=? order by seq desc limit ?"

const DeleteRealtimeByUserIdSeq = "delete from realtime_archive where user_id=? and seq=?"

// realtime end
