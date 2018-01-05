package orders

import (
	"net/http"
	"html/template"
	"encoding/json"
)

var templates = template.Must(template.ParseFiles("./orders/orders.html"))

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func respondError(w http.ResponseWriter, status int, err error) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func OrdersHandler(rw http.ResponseWriter, req *http.Request) {
	orders, err := getOrders()
	if err != nil {
		respondError(rw, http.StatusInternalServerError, err)
		return
	}

	if err := templates.ExecuteTemplate(rw, "orders.html", orders); err != nil {
		respondError(rw, http.StatusInternalServerError, err)
	}
}


func getOrders() ([]string, error) {
	return []string{"xxxxx-yyyy-zzzzzzzz"}, nil
}