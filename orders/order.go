package orders

import (
	"net/http"
	"html/template"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

var templates = template.Must(template.ParseFiles("./orders/orders.html", "./orders/orderstatus.html"))
var awsSession = session.Must(session.NewSession())
var modelInstanceTable =os.Getenv("MODEL_INSTANCE_TABLE")
var ddb = dynamodb.New(awsSession)


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

func OrderStatusHandler(rw http.ResponseWriter, req *http.Request) {
	if err := templates.ExecuteTemplate(rw, "orderstatus.html", nil); err != nil {
		respondError(rw, http.StatusInternalServerError, err)
	}
}


func getOrders() ([]string, error) {

	resultsLimit := int64(50)
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]*string{
			"#ST": aws.String("state"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String("OrderReceived"),
			},
		},
		FilterExpression:     aws.String("#ST = :s"),
		TableName:            aws.String(modelInstanceTable),
		Limit: &resultsLimit,
	}

	qout, err := ddb.Scan(input)
	if err != nil {
		return nil, err
	}

	var orders []string

	items := qout.Items
	for _, item := range items {
		orders = append(orders, *item["instanceId"].S)
	}

	return orders, nil
}