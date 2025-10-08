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

var (
	multi  = pflag.BoolP("multi", "m", false, "multi-select")
	height = pflag.IntP("height", "h", 0, "height of the finder")
	border = pflag.BoolP("border", "b", false, "draw a border around the finder")
	borderCharsStr = pflag.String("border-chars", "", "custom border characters (6 runes: tl, tr, bl, br, h, v)")
)

func main() {
	pflag.Parse()

	if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Println("please use pipe")
		return
	}
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	slice := strings.Split(string(b), "\n")

	var opts []fuzzyfinder.Option
	if *height > 0 {
		opts = append(opts, fuzzyfinder.WithHeight(*height))
	}
	if *border {
		opts = append(opts, fuzzyfinder.WithBorder())
	}
	if *borderCharsStr != "" {
		runes := []rune(*borderCharsStr)
		if len(runes) == 6 {
			opts = append(opts, fuzzyfinder.WithBorderChars(runes))
		} else {
			log.Println("warning: --border-chars expects 6 runes, ignoring.")
		}
	}

	var idxs []int
	if *multi {
		idxs, err = fuzzyfinder.FindMulti(
			slice,
			func(i int) string {
				return slice[i]
			},
			opts...,
		)
	} else {
		var idx int
		idx, err = fuzzyfinder.Find(
			slice,
			func(i int) string {
				return slice[i]
			},
			opts...,
		)
		idxs = []int{idx}
	}

	if err != nil {
		log.Fatal(err)
	}
	for _, idx := range idxs {
		fmt.Println(slice[idx])
	}
}
