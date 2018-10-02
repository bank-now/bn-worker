package main

import (
	"fmt"
	"github.com/bank-now/bn-worker/io"
)

func main() {
	io.Setup(handle)

}

func handle(b []byte) {
	fmt.Println(string(b))
}
