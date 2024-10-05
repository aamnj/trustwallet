package services

import (
	"amanj/trustwallet/ethparser"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JSONResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type subscribeRequest struct {
	Address string
}

type Services struct {
	ethParser ethparser.Parser
}

func NewServices(parser ethparser.Parser) *Services {
	return &Services{
		ethParser: parser,
	}
}

func (svc *Services) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	block := svc.ethParser.GetCurrentBlock()

	JSONResponseHandler(w, &JSONResponse{
		Status: http.StatusOK,
		Data: map[string]interface{}{
			"block": block,
		},
	})
}

func (svc *Services) Subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		JSONErrorHandler(w, "Unable to read request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var subscribeReq *subscribeRequest
	if err := json.Unmarshal(body, &subscribeReq); err != nil {
		JSONErrorHandler(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TODO: Add validation for allowed address type
	if subscribeReq.Address == "" {
		JSONErrorHandler(w, "Invalid Address", http.StatusBadRequest)
		return
	}

	if !svc.ethParser.Subscribe(subscribeReq.Address) {
		JSONErrorHandler(w, fmt.Sprintf("Already subscribed to address %v", subscribeReq.Address), http.StatusBadRequest)
		return
	}

	JSONResponseHandler(w, &JSONResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Address %v subscribed", subscribeReq.Address),
	})
}

func (svc *Services) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		JSONErrorHandler(w, "Please provide address to fetch transactions", http.StatusBadRequest)
		return
	}

	txns := svc.ethParser.GetTransactions(address)
	if len(txns) == 0 {
		JSONErrorHandler(w, "No transactions found for this address", http.StatusBadRequest)
		return
	}

	JSONResponseHandler(w, &JSONResponse{
		Status: http.StatusOK,
		Data:   txns,
	})
}

func JSONResponseHandler(w http.ResponseWriter, response *JSONResponse) {
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code
	w.WriteHeader(http.StatusOK)

	// Convert the response object to JSON and send it
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func JSONErrorHandler(w http.ResponseWriter, err string, status int) {
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code
	w.WriteHeader(status)

	// Convert the response object to JSON and send it
	err2 := json.NewEncoder(w).Encode(map[string]interface{}{"status": status, "error": err})
	if err2 != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}
