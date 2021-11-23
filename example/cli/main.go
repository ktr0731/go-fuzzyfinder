package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	slice := strings.Split(string(b), "\n")

	idxs, err := fuzzyfinder.FindMulti(
		slice,
		func(i int) string {
			return slice[i]
		})
	if err != nil {
		log.Fatal(err)
	}
	for _, idx := range idxs {
		fmt.Println(slice[idx])
	}
}
