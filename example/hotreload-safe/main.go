package main

import (
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
	"log"
	"strconv"
	"sync"
	"time"
)

func main() {
	var slice []string
	var rwMut sync.RWMutex

	go func(slice *[]string) {
		var count int
		for _ = range time.Tick(1 * time.Second) {
			count++
			rwMut.Lock()
			*slice = append(*slice, strconv.Itoa(count))
			rwMut.Unlock()
		}
	}(&slice)

	idx, err := fuzzyfinder.Search(
		func() []string {
			rwMut.Lock()
			defer rwMut.Unlock()

			return slice
		},
		fuzzyfinder.WithHotReload(),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(slice[idx])
}
