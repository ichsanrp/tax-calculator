// Model is where you should keep your data model, the algorithms.
package taxCalculator

import (
	"context"
	"errors"
	"log"
	"time"
)

var errParameterNotSatisfied = errors.New("Required Parameter Missing")

// initModel is model initialization, it will try to prepare table on database and preparing sql query in statement for later use
func (m *TaxCalculator) initModel() (err error) {
	err = m.initTable()
	err = m.initStatement()
	return
}

// createSession will create new session in database and returning this session id for later use
func (m *TaxCalculator) createSession() (session *Session, err error) {
	session = &Session{
		CreateTime: time.Now(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	tx, err := m.db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	_, err = tx.Stmt(m.createSessionStatement).Exec(session.CreateTime)
	err = tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&session.ID)
	return
}

// addItem is function to create an item in database and doing tax calculation before item saved to database
func (m *TaxCalculator) addItem(item *Item) (newitem *Item, err error) {
	if item.SessionID == 0 {
		err = errParameterNotSatisfied
		return
	}
	newitem = &Item{
		Name:      item.Name,
		Price:     item.Price,
		TaxCodeID: item.TaxCodeID,
		SessionID: item.SessionID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	tx, err := m.db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// get last weight
	err = tx.QueryRow("SELECT max(weight) FROM tax_calculator.item WHERE session_id = ?", item.SessionID).Scan(newitem.Weight)

	// get tax
	newitem.Tax = calculateTax(newitem.Price, newitem.TaxCodeID)
	_, err = tx.Stmt(m.addItemStatement).Exec(newitem.Weight, newitem.Name, newitem.Price, newitem.Tax, newitem.TaxCodeID, item.SessionID)
	err = tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&newitem.ID)

	return
}

// deleteItem is function to delete item in database
func (m *TaxCalculator) deleteItem(item *Item) (err error) {
	if item.ID == 0 {
		return errParameterNotSatisfied
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	tx, err := m.db.BeginTx(ctx, nil)
	defer func() {
		log.Println(err)
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	_, err = tx.Stmt(m.removeItemStatement).Exec(item.ID)
	return
}

// getItem is used to create Item object from database based on itemId as parameters
func (m *TaxCalculator) getItem(itemId int) (item *Item, err error) {
	if itemId == 0 {
		return nil, errParameterNotSatisfied
	}
	item = &Item{
		ID: itemId,
	}
	err = m.getitemByIDStatement.QueryRow(itemId).Scan(
		&item.Name,
		&item.Price,
		&item.Tax,
		&item.TaxCodeID,
		&item.Weight,
		&item.SessionID,
	)
	return
}

// updateItem is function to create an item in database and doing tax calculation before item saved to database
func (m *TaxCalculator) updateItem(item *Item) (err error) {

	if item.ID == 0 {
		return errParameterNotSatisfied
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	tx, err := m.db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	// get tax
	item.Tax = calculateTax(item.Price, item.TaxCodeID)
	_, err = tx.Stmt(m.editItemStatement).Exec(item.Name, item.Price, item.Tax, item.TaxCodeID, item.ID)
	return
}

// getItemsBySession is function to get all item based on session information given
func (m *TaxCalculator) getItemsBySession(sessionID int) (items []*Item, err error) {
	if sessionID == 0 {
		return nil, errParameterNotSatisfied
	}

	items = make([]*Item, 0)
	rows, err := m.getitemBySessionStatement.Query(sessionID)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		item := &Item{}
		err = rows.Scan(&item.ID, &item.Name, &item.Price, &item.Tax, &item.TaxCodeID)
		if err != nil {
			log.Println("[getItemsBySession]", err)
			continue
		}
		item.SessionID = sessionID
		items = append(items, item)
	}

	return
}

// getAllSession is function to get latest sessions created.
func (m *TaxCalculator) getAllSession(page, perPage int) (sessions []*Session, err error) {
	if page == 0 {
		page = 1
	}

	if perPage == 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	sessions = make([]*Session, 0)
	rows, err := m.getAllSessionStatement.Query(offset, perPage)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var t string
		session := &Session{}
		err = rows.Scan(&session.ID, &t)
		ctime, _ := time.Parse("2006-01-02 15:04:05", t)
		session.CreateTime = ctime
		if err != nil {
			log.Println("[getAllSession]", err)
			continue
		}
		sessions = append(sessions, session)
	}

	return
}
