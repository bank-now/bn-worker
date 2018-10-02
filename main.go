package main

import (
	"fmt"
	"github.com/bank-now/bn-common-io/queues/sub"
)

func main() {

	c := sub.Config{Topic: "interest-calculation-v1",
		Version: "v1",
		Name:    "worker",
		F:       handle}
	sub.Subscribe(c)

}

func handle(b []byte) {
	fmt.Println(string(b))
}
