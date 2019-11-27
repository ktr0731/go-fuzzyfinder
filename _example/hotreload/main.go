package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	isatty "github.com/mattn/go-isatty"
	"github.com/spf13/pflag"
)

var multi = pflag.BoolP("multi", "m", false, "multi-select")

func main() {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Println("please use pipe")
		return
	}
	var slice []string
	go func(slice *[]string) {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			*slice = append(*slice, s.Text())
		}
	}(&slice)

	idxs, err := fuzzyfinder.FindMulti(
		&slice,
		func(i int) string {
			return slice[i]
		},
		fuzzyfinder.WithHotReload(),
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, idx := range idxs {
		fmt.Println(slice[idx])
	}
}
