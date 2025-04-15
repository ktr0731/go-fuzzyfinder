package main

import (
	"fmt"
	"log"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

type Track struct {
	Name      string
	AlbumName string
	Artist    string
}

var tracks = []Track{
	{"foo", "album1", "artist1"},
	{"bar", "album1", "artist1"},
	{"foo", "album2", "artist1"},
	{"baz", "album2", "artist2"},
	{"baz", "album3", "artist2"},
}

func main() {
	singleExample()
	multiExample()
}

func singleExample() {
	idx, err := fuzzyfinder.Find(
		tracks,
		func(i int) string {
			return tracks[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Track: %s (%s)\nAlbum: %s",
				tracks[i].Name,
				tracks[i].Artist,
				tracks[i].AlbumName)
		}),
		fuzzyfinder.WithPreselected(func(i int) bool {
			return i == 1
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Selected: %s\n", tracks[idx].Name)
}

func multiExample() {
	idxs, err := fuzzyfinder.FindMulti(
		tracks,
		func(i int) string {
			return tracks[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Track: %s (%s)\nAlbum: %s",
				tracks[i].Name,
				tracks[i].Artist,
				tracks[i].AlbumName)
		}),
		fuzzyfinder.WithPreselected(func(i int) bool {
			return tracks[i].Artist == "artist2"
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Selected:")
	for _, idx := range idxs {
		fmt.Printf("- %v (%s)\n", idx, tracks[idx].Name)
	}
}
