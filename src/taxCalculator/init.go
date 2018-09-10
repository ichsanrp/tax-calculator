package taxCalculator

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

//TaxCalculator module definition, we store database and all query statement here.
type TaxCalculator struct {
	db                        *sql.DB
	getAllSessionStatement    *sql.Stmt
	createSessionStatement    *sql.Stmt
	deleteSessionStatement    *sql.Stmt
	getitemBySessionStatement *sql.Stmt
	getitemByIDStatement      *sql.Stmt
	addItemStatement          *sql.Stmt
	editItemStatement         *sql.Stmt
	removeItemStatement       *sql.Stmt
}

//Init initialize module, init model and controller
func Init(router *httprouter.Router) (module *TaxCalculator, err error) {
	module = &TaxCalculator{}
	err = module.initDB()
	if err != nil {
		log.Fatalln("[DB CONNECTION]", err)
	}
	err = module.initModel()
	if err != nil {
		log.Fatalln("[MODEL INIT]", err)
	}
	err = module.initController(router)
	if err != nil {
		log.Fatalln("[CONTROLLER INIT]", err)
	}
	return
}

// initDB will create new connection to database via mysql driver
func (m *TaxCalculator) initDB() (err error) {
	m.db, err = sql.Open("mysql", "root:rahasia@tcp(mysql)/")
	return
}
