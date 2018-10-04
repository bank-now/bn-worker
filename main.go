package main

import (
	"fmt"
	"github.com/bank-now/bn-common-io/queues/sub"
	"github.com/bank-now/bn-common-model/common/operation"
	"github.com/sirupsen/logrus"
)

func main() {

	c := sub.Config{Topic: "interest-calculation-v1",
		Version: "v1",
		Name:    "worker",
		F:       handle}
	sub.Subscribe(c)

}

func handle(b []byte) {
	i, err := operation.GetInterestOperation(b)
	if err != nil {
		logrus.Errorln("Did not understand: ", err)
	}
	fmt.Println(i)
}
