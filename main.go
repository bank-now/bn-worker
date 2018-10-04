package main

import (
	"fmt"
	"github.com/bank-now/bn-common-io/queues/sub"
	"github.com/bank-now/bn-common-model/common/model"
	"github.com/bank-now/bn-common-model/common/operation"
	"github.com/bank-now/bn-worker/controller"
	"github.com/google/uuid"
	"math"
	"time"
)

const (
	name    = "worker"
	version = "v1"
)

func main() {
	c := sub.Config{Topic: operation.InterestOperationV1Topic,
		Name:    name,
		Version: version,
		F:       handle}
	sub.Subscribe(c)

}

func handle(b []byte) {
	i, err := operation.GetInterestOperation(b)
	if err != nil {
		panic(err)
	}
	fmt.Println("Fetching account:", i.Account)

	trxSlice, err := controller.GetTransactionsByAccountId(i.Account)
	if err != nil {
		//TODO: dead-letter queue!
		return
	}
	transactions := model.OrderTransactions(*trxSlice)
	var balance float64
	balance = 0
	for _, trx := range transactions {
		balance += trx.Amount
	}
	fmt.Println("Balance is: ", balance)
	interest := calculateInterest(balance, 1)
	fmt.Println("Interest is: ", interest)

	u, err := uuid.NewRandom()

	interestTransaction := model.Transaction{
		ID:         u.String(),
		Amount:     interest,
		AccountID:  i.Account,
		SystemCode: "INTEREST for Day",
		Timestamp:  time.Now()}

	fmt.Println(interestTransaction)

}

/*
	WARNING: This is not accurate.
*/
func calculateInterest(balance float64, days float64) float64 {
	var i, n, inner, power float64
	i = 0.1    //10% per annum
	n = 365.25 //Days per year
	inner = i/n + 1
	power = 1
	brackets := math.Pow(inner, power)
	return balance*brackets - balance

}
