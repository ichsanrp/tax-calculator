package taxCalculator

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (m *TaxCalculator) initController(router *httprouter.Router) (err error) {

	router.GET("/session", m.getAllSessionHandler)
	router.GET("/session/:id", m.getSessionItemHandler)
	router.POST("/session", m.createSessionHandler)

	router.POST("/item", m.addItemHandler)
	router.DELETE("/item", m.removeItemHandler)
	router.PUT("/item", m.updateItemHandler)
	router.OPTIONS("/item", m.optionsHandler)

	return
}

func (m *TaxCalculator) optionsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers := w.Header()
	headers.Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}

func (m *TaxCalculator) getAllSessionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := &DefaultResponse{
		Header: &DefaultResponseHeader{},
	}

	page, _ := strconv.Atoi(r.FormValue("page"))
	perPage, _ := strconv.Atoi(r.FormValue("per_page"))

	sessions, err := m.getAllSession(page, perPage)
	if err != nil {
		response.Header.StatusCode = http.StatusInternalServerError
		response.Header.Error = err.Error()
	} else {
		response.Header.StatusCode = http.StatusOK
	}

	response.Data = sessions
	b, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "Application/json")
	w.Write(b)
}

func (m *TaxCalculator) getSessionItemHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := &DefaultResponse{
		Header: &DefaultResponseHeader{},
	}

	sessionID, _ := strconv.Atoi(ps.ByName("id"))

	items, err := m.getItemsBySession(sessionID)
	if err != nil {
		response.Header.StatusCode = http.StatusInternalServerError
		response.Header.Error = err.Error()
	} else {
		response.Header.StatusCode = http.StatusOK
	}

	response.Data = items
	b, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "Application/json")
	w.Write(b)
}

func (m *TaxCalculator) createSessionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := &DefaultResponse{
		Header: &DefaultResponseHeader{},
	}

	session, err := m.createSession()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Header.StatusCode = http.StatusInternalServerError
		response.Header.Error = err.Error()
		return
	}

	response.Data = session
	response.Header.StatusCode = http.StatusOK
	b, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "Application/json")
	w.Write(b)
}

func (m *TaxCalculator) addItemHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := &DefaultResponse{
		Header: &DefaultResponseHeader{},
	}

	item := parseItemFromRequest(r)
	newItem, err := m.addItem(item)
	if err != nil {
		response.Header.StatusCode = http.StatusInternalServerError
		response.Header.Error = err.Error()
	} else {
		response.Header.StatusCode = http.StatusOK
	}

	response.Data = newItem
	b, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "Application/json")
	w.Write(b)
}

func (m *TaxCalculator) removeItemHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := &DefaultResponse{
		Header: &DefaultResponseHeader{},
	}

	item := parseItemFromRequest(r)
	err := m.deleteItem(item)
	if err != nil {
		response.Header.StatusCode = http.StatusInternalServerError
		response.Header.Error = err.Error()
	} else {
		response.Header.StatusCode = http.StatusOK
	}

	response.Data = item
	b, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "Application/json")
	w.Write(b)
}

func (m *TaxCalculator) updateItemHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := &DefaultResponse{
		Header: &DefaultResponseHeader{},
	}

	item := parseItemFromRequest(r)
	err := m.updateItem(item)
	if err != nil {
		response.Header.StatusCode = http.StatusInternalServerError
		response.Header.Error = err.Error()
	} else {
		response.Header.StatusCode = http.StatusOK
	}

	response.Data = item
	b, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "Application/json")
	w.Write(b)
}
