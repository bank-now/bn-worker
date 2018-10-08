package main

import (
	"encoding/json"
	"fmt"
	"github.com/bank-now/bn-common-io/queues/pub"
	"github.com/bank-now/bn-common-io/queues/sub"
	"github.com/bank-now/bn-common-model/common/model"
	"github.com/bank-now/bn-common-model/common/operation"
	"github.com/bank-now/bn-worker/controller"
	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
	"log"
	"math"
	"time"
)

const (
	Name    = "worker"
	Version = "v1"
	Address = "192.168.88.24:4150"
	Action  = "workItem"
)

var (
	fullName = fmt.Sprint(Name, "-", Version, "-", Action)
	producer *nsq.Producer
)

func main() {
	var err error

	pubConfig := pub.Config{Topic: operation.WriteOperationV1Topic,
		Name:    Name,
		Version: Version,
		Address: Address}
	producer, err = pub.Setup(pubConfig)
	if err != nil {
		log.Fatal(err)
	}

	c := sub.Config{
		Topic:   operation.InterestOperationV2Topic,
		Name:    Name,
		Version: Version,
		Address: Address,
		F:       handle}
	sub.Subscribe(c)

}

func handle(b []byte) {
	i, err := operation.GetInterestOperation(b)
	if err != nil {
		panic(err)
	}

	transactions, err := getOrderedTransactions(i.Account)
	if err != nil {
		//TODO: dead-letter queue!
	}

	item := doInterestCalculation(i.Account, transactions)
	writeTransaction(item)

}

func getOrderedTransactions(account string) (transactions []model.Transaction, err error) {
	trxSlice, err := controller.GetTransactionsByAccountId(account)
	if err != nil {
		return
	}
	transactions = model.OrderTransactions(*trxSlice)
	return

}

func doInterestCalculation(account string, transactions []model.Transaction) (transaction model.Transaction) {
	var balance float64 = 0
	for _, trx := range transactions {
		balance += trx.Amount
	}
	interest := calculateInterest(balance, 1)
	u, _ := uuid.NewRandom()
	transaction = model.Transaction{
		ID:         u.String(),
		Amount:     interest,
		AccountID:  account,
		SystemCode: "INTEREST for Day",
		Timestamp:  time.Now()}

	return
}

func writeTransaction(item model.Transaction) (err error) {
	intTrxB, _ := json.Marshal(item)
	write := operation.WriteOperationV1{
		Table:  model.TransactionTable,
		Method: "POST",
		Item:   intTrxB}

	writeB, err := json.Marshal(write)
	if err != nil {
		return
	}
	producer.Publish(operation.WriteOperationV1Topic, writeB)
	return

}

/*
	WARNING: This is not accurate.
*/
func calculateInterest(balance float64, days float64) float64 {
	var i, n, inner, power float64
	i = 0.1    //10% per annum
	n = 365.25 //Days per year
	inner = i/n + 1
	power = 1 //Super wrong
	brackets := math.Pow(inner, power)
	return balance*brackets - balance

}
