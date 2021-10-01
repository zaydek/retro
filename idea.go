package main

import (
	"fmt"
	"time"

	"github.com/zaydek/retro/go/pkg/watch"
)

func main() {
	for result := range watch.Directory("lol", 100*time.Millisecond) {
		if result.Err != nil {
			panic(fmt.Errorf("watch.Directory: %w", result.Err))
		}
		fmt.Println(result)
	}
}
