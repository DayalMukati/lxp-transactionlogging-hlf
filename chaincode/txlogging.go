package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionLog provides functions for logging and querying transactions
type TransactionLog struct {
	contractapi.Contract
}

// Transaction represents a transaction between two users
type Transaction struct {
	TransactionID string  `json:"transactionID"`
	Sender        string  `json:"sender"`
	Receiver      string  `json:"receiver"`
	Amount        float64 `json:"amount"`
}

// LogTransaction logs a new transaction
func (t *TransactionLog) LogTransaction(ctx contractapi.TransactionContextInterface, transactionID string, sender string, receiver string, amount float64) error {
	transaction := Transaction{
		TransactionID: transactionID,
		Sender:        sender,
		Receiver:      receiver,
		Amount:        amount,
	}

	transactionAsBytes, _ := json.Marshal(transaction)

	return ctx.GetStub().PutState(transactionID, transactionAsBytes)
}

// QueryTransaction retrieves a specific transaction by ID
func (t *TransactionLog) QueryTransaction(ctx contractapi.TransactionContextInterface, transactionID string) (*Transaction, error) {
	transactionAsBytes, err := ctx.GetStub().GetState(transactionID)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if transactionAsBytes == nil {
		return nil, fmt.Errorf("transaction %s does not exist", transactionID)
	}

	transaction := new(Transaction)
	_ = json.Unmarshal(transactionAsBytes, transaction)

	return transaction, nil
}

// QueryAllTransactions returns all logged transactions
func (t *TransactionLog) QueryAllTransactions(ctx contractapi.TransactionContextInterface) ([]Transaction, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transactions []Transaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var transaction Transaction
	 _ = json.Unmarshal(queryResponse.Value, &transaction)
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(TransactionLog))

	if err != nil {
		fmt.Printf("Error creating transaction log chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting transaction log chaincode: %s", err.Error())
	}
}
