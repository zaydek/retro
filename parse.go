package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	n, err := strconv.Atoi(strings.ReplaceAll("10_000", "_", ""))
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}
