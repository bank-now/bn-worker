package main

import (
	"fmt"
	"github.com/bank-now/bn-common-io/queues/sub"
	"github.com/bank-now/bn-common-model/common/operation"
)

func main() {

	c := sub.Config{Topic: operation.InterestOperationV1Topic,
		Name:    "worker",
		Version: "v1",
		F:       handle}
	sub.Subscribe(c)

}

func handle(b []byte) {
	i, err := operation.GetInterestOperation(b)
	if err != nil {
		panic(err)
	}
	fmt.Println(i.Account)
}
