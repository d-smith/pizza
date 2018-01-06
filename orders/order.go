package orders

import (
	"net/http"
	"html/template"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gorilla/mux"
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
	vars := mux.Vars(req)
	order,  err := getOrderStatus(vars["orderid"])
	if err != nil {
		respondError(rw, http.StatusInternalServerError, err)
		return
	}

	statusDetail := orderStatusFromTxns(order)

	if err = templates.ExecuteTemplate(rw, "orderstatus.html", statusDetail); err != nil {
		respondError(rw, http.StatusInternalServerError, err)
	}
}

func getOrderStatus(order string)(map[string]string, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":iid": {
				S: aws.String(order),
			},
		},
		KeyConditionExpression: aws.String("instanceId = :iid"),
		TableName:              aws.String(modelInstanceTable),
	}

	qout, err := ddb.Query(input)
	if err != nil {
		return nil, err
	}

	orderTxns := make(map[string]string)
	items := qout.Items
	for _, item := range items {
		orderTxns[*item["state"].S] = *item["txnId"].S
	}

	return orderTxns, nil
}

type orderStepState struct {
	Step string
	Status string
}

func orderStatusFromTxns(order map[string]string) []orderStepState {
	states := []orderStepState{}

	steps := []string{"OrderReceived","AssemblingPizza","CookingPizza","OrderReady"}
	for _, step := range steps {
		txn := order[step]
		switch txn {
		case "":
			states = append(states, orderStepState{
				Step: step,
				Status: "",
			})
		default:
			states = append(states, orderStepState{
				Step: step,
				Status: "is-complete",
			})
		}
	}

	return states

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