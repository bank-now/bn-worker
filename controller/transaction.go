package controller

import (
	"fmt"
	"github.com/bank-now/bn-common-io/rest"
	"github.com/bank-now/bn-common-model/common/model"
)

const (
	grest = "http://192.168.88.24:3001/"
)

func GetTransactionsByAccountId(account string) (trxs *[]model.Transaction, err error) {
	url := fmt.Sprint(grest, model.TransactionTable, "?id=eq.", account)
	jBytes, err := rest.Get(url)
	if err != nil {
		return
	}
	return model.GetTransactions(jBytes)

}
