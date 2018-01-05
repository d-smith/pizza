package main

import (
	"github.com/d-smith/pizza/orders"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
)



func main() {
	r := mux.NewRouter()
	r.HandleFunc("/orders", orders.OrdersHandler)
	r.HandleFunc("/orderstatus/{orderid}", orders.OrderStatusHandler)
	http.Handle("/", r)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe(fmt.Sprintf(":%d", 8778), nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}