// Model is where you should keep your data model, the algorithms.
package taxCalculator

import (
	"database/sql"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	_ "github.com/proullon/ramsql/driver"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var calculateTaxModule *TaxCalculator
var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	var err error
	var db *sql.DB
	db, mock, err = sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	//prepare statement mock
	mock.ExpectPrepare("SELECT id, created_time FROM tax_calculator.session")
	mock.ExpectPrepare("INSERT INTO tax_calculator.session")
	mock.ExpectPrepare("DELETE FROM tax_calculator.session")
	mock.ExpectPrepare("SELECT id, name, price, tax, tax_code FROM tax_calculator.item")
	mock.ExpectPrepare("SELECT name, price, tax, tax_code, weight, session_id FROM tax_calculator.item")
	// mock.ExpectPrepare("SELECT max(weight) FROM tax_calculator.item WHERE session_id = ?")
	mock.ExpectPrepare("INSERT INTO tax_calculator.item")
	mock.ExpectPrepare("UPDATE tax_calculator.item")
	mock.ExpectPrepare("DELETE FROM tax_calculator.item")

	//creating module
	calculateTaxModule = &TaxCalculator{}
	calculateTaxModule.db = db
	calculateTaxModule.initStatement()

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestTaxCalculator_createSession(t *testing.T) {
	tests := []struct {
		name        string
		wantSession *Session
		wantErr     bool
	}{
		{
			name: "test1",
			wantSession: &Session{
				CreateTime: time.Now(),
				ID:         1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO tax_calculator.session").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery("SELECT LAST_INSERT_ID()").WillReturnRows(rows)
			if tt.wantErr {
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}

			gotSession, err := calculateTaxModule.createSession()
			if (err != nil) != tt.wantErr {
				t.Errorf("TaxCalculator.createSession() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			if gotSession.ID != tt.wantSession.ID {
				t.Errorf("TaxCalculator.createSession() = %v, want %v", gotSession, tt.wantSession)
			}
		})
	}
}

func TestTaxCalculator_addItem(t *testing.T) {
	type args struct {
		item *Item
	}
	tests := []struct {
		name        string
		args        args
		wantNewitem *Item
		wantErr     bool
	}{
		{
			name: "pizza",
			wantNewitem: &Item{
				ID:        1,
				Name:      "pizza",
				Price:     100,
				TaxCodeID: 1,
				Tax:       10,
			},
			args: args{
				item: &Item{
					ID:        1,
					Name:      "pizza",
					Price:     100,
					TaxCodeID: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "malboro",
			wantNewitem: &Item{
				ID:        1,
				Name:      "malboro",
				Price:     100,
				TaxCodeID: 2,
				Tax:       22,
			},
			args: args{
				item: &Item{
					ID:        1,
					Name:      "malboro",
					Price:     100,
					TaxCodeID: 2,
				},
			},
			wantErr: false,
		},
		{
			name: "movie ticket",
			wantNewitem: &Item{
				ID:        1,
				Name:      "movie ticket",
				Price:     100,
				TaxCodeID: 3,
				Tax:       0,
			},
			args: args{
				item: &Item{
					ID:        1,
					Name:      "movie ticket",
					Price:     100,
					TaxCodeID: 3,
				},
			},
			wantErr: false,
		},
		{
			name: "club ticket",
			wantNewitem: &Item{
				ID:        1,
				Name:      "movie ticket",
				Price:     200,
				TaxCodeID: 3,
				Tax:       1,
			},
			args: args{
				item: &Item{
					ID:        1,
					Name:      "movie ticket",
					Price:     200,
					TaxCodeID: 3,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO tax_calculator.item").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery("SELECT LAST_INSERT_ID()").WillReturnRows(rows)
			if tt.wantErr {
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}

			gotNewitem, err := calculateTaxModule.addItem(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaxCalculator.addItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewitem, tt.wantNewitem) {
				t.Errorf("TaxCalculator.addItem() = %v, want %v", gotNewitem, tt.wantNewitem)
			}
		})
	}
}

func TestTaxCalculator_deleteItem(t *testing.T) {

	type args struct {
		item *Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "club ticket",
			args: args{
				item: &Item{
					ID:        1,
					Name:      "movie ticket",
					Price:     200,
					TaxCodeID: 3,
				},
			},
			wantErr: false,
		},
		{
			name: "club ticket",
			args: args{
				item: &Item{
					Name:      "movie ticket",
					Price:     200,
					TaxCodeID: 3,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec("DELETE FROM tax_calculator.item").WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.wantErr {
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}
			if err := calculateTaxModule.deleteItem(tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("TaxCalculator.deleteItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaxCalculator_updateItem(t *testing.T) {

	type args struct {
		item *Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "club ticket",
			args: args{
				item: &Item{
					ID:        1,
					Name:      "movie ticket",
					Price:     200,
					TaxCodeID: 3,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE tax_calculator.item").WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.wantErr {
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}
			if err := calculateTaxModule.updateItem(tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("TaxCalculator.updateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
