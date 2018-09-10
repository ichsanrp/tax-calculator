package taxCalculator

import (
	"log"
)

func (m *TaxCalculator) initStatement() (err error) {
	err = m.prepareStatement()
	if err != nil {
		log.Fatalln("[DATABASE INIT] failed to prepare statement", err)
	}
	return
}

func (m *TaxCalculator) initTable() (err error) {
	err = m.prepareTable()
	if err != nil {
		log.Fatalln("[DATABASE INIT] failed to prepare table", err)
	}
	return
}

//  prepareTable initiating database table setting to mysql, creating table if table not exist
func (m *TaxCalculator) prepareTable() (err error) {

	// creating table session
	q := `
	CREATE TABLE IF NOT EXISTS tax_calculator.session (
		id INT NOT NULL AUTO_INCREMENT,
		created_time DATETIME NULL,
		PRIMARY KEY (id))
	`

	_, err = m.db.Exec(q)

	// creating table item
	q = `
	CREATE TABLE IF NOT EXISTS tax_calculator.item (
		id INT NOT NULL AUTO_INCREMENT,
		weight SMALLINT NOT NULL,
		name VARCHAR(128) NOT NULL,
		price FLOAT NOT NULL,
		tax FLOAT NULL,
		tax_code INT NOT NULL,
		session_id INT NOT NULL,
		PRIMARY KEY (id),
		INDEX fk_item_session1_idx (session_id ASC),
		CONSTRAINT fk_item_session1
		  FOREIGN KEY (session_id)
		  REFERENCES tax_calculator.session (id)
		  ON DELETE NO ACTION
		  ON UPDATE NO ACTION)
	`

	_, err = m.db.Exec(q)

	return
}

func (m *TaxCalculator) prepareStatement() (err error) {
	var q string

	// table session Operation
	q = `SELECT id, created_time FROM tax_calculator.session ORDER BY created_time DESC LIMIT ?, ?`
	m.getAllSessionStatement, err = m.db.Prepare(q)
	if err != nil {
		log.Fatalln("[STATEMENT ERR] getAllSessionStatement ", err)
	}

	q = `INSERT INTO tax_calculator.session (created_time) VALUES (?);`
	m.createSessionStatement, err = m.db.Prepare(q)
	if err != nil {
		log.Fatalln("[STATEMENT ERR] createSessionStatement ", err)
	}

	q = `DELETE FROM tax_calculator.session where id = ?`
	m.deleteSessionStatement, err = m.db.Prepare(q)
	if err != nil {
		log.Fatalln("[STATEMENT ERR] deleteSessionStatement ", err)
	}

	q = `SELECT id, name, price, tax, tax_code FROM tax_calculator.item 
		 WHERE session_id = ?
		 ORDER BY weight ASC`
	m.getitemBySessionStatement, err = m.db.Prepare(q)
	if err != nil {
		log.Fatalln("[STATEMENT ERR] getitemBySessionStatement ", err)
	}

	q = `SELECT name, price, tax, tax_code, weight, session_id FROM tax_calculator.item 
		 WHERE id = ?`
	m.getitemByIDStatement, err = m.db.Prepare(q)
	if err != nil {
		log.Fatalln("[STATEMENT ERR] getitemByIdStatement ", err)
	}

	// table item Operation
	q = `INSERT INTO tax_calculator.item (weight, name, price, tax, tax_code, session_id) 
		 VALUES (?, ?, ?, ?, ?, ?)`
	m.addItemStatement, err = m.db.Prepare(q)
	if err != nil {
		log.Fatalln("[STATEMENT ERR] addItemStatement ", err)
	}

	q = `UPDATE tax_calculator.item SET name = ?, price = ?, tax = ?, tax_code = ? 
		 WHERE id = ?`
	m.editItemStatement, err = m.db.Prepare(q)
	if err != nil {
		log.Fatalln("[STATEMENT ERR] editItemStatement ", err)
	}

	q = `DELETE FROM tax_calculator.item where id = ?`
	m.removeItemStatement, err = m.db.Prepare(q)
	if err != nil {
		log.Fatalln("[STATEMENT ERR] removeItemStatement ", err)
	}

	return
}
