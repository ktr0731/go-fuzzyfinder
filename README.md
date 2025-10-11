# go-fuzzyfinder

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ktr0731/go-fuzzyfinder)](https://pkg.go.dev/github.com/ktr0731/go-fuzzyfinder)
[![GitHub Actions](https://github.com/ktr0731/go-fuzzyfinder/workflows/main/badge.svg)](https://github.com/ktr0731/go-fuzzyfinder/actions)
[![codecov](https://codecov.io/gh/ktr0731/go-fuzzyfinder/branch/master/graph/badge.svg?token=RvpSTKDJGO)](https://codecov.io/gh/ktr0731/go-fuzzyfinder)

`go-fuzzyfinder` is a Go library that provides fuzzy-finding with an fzf-like terminal user interface.

![](https://user-images.githubusercontent.com/12953836/52424222-e1edc900-2b3c-11e9-8158-8e193844252a.png)

## Installation
``` bash
$ go get github.com/ktr0731/go-fuzzyfinder
```

## Usage
`go-fuzzyfinder` provides two functions, `Find` and `FindMulti`.
`FindMulti` can select multiple lines. It is similar to `fzf -m`.

This is [an example](//github.com/ktr0731/go-fuzzyfinder/blob/master/example/track/main.go) of `FindMulti`.

``` go
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
    idx, err := fuzzyfinder.FindMulti(
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
        }))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("selected: %v\n", idx)
}
```

The execution result prints selected item's indexes.

### Preselecting items

You can preselect items using the `WithPreselected` option. It works in both `Find` and `FindMulti`.

``` go
// Single selection mode
// The cursor will be positioned on the first item that matches the predicate
idx, err := fuzzyfinder.Find(
    tracks,
    func(i int) string {
        return tracks[i].Name
    },
    fuzzyfinder.WithPreselected(func(i int) bool {
        return tracks[i].Name == "bar"
    }),
)

// Multi selection mode
// All items that match the predicate will be selected initially
idxs, err := fuzzyfinder.FindMulti(
    tracks,
    func(i int) string {
        return tracks[i].Name
    },
    fuzzyfinder.WithPreselected(func(i int) bool {
        return tracks[i].Artist == "artist2"
    }),
)
```

### Customizing appearance

You can customize the fuzzy finder's dimensions and border using several options:

``` go
idx, err := fuzzyfinder.Find(
    tracks,
    func(i int) string {
        return tracks[i].Name
    },
    fuzzyfinder.WithHeight(10),              // Limit height to 10 lines
    fuzzyfinder.WithWidth(80),               // Limit width to 80 columns
    fuzzyfinder.WithBorder(),                // Enable border
    fuzzyfinder.WithBorderChars([]rune{'╭', '╮', '╰', '╯', '─', '│'}), // Custom border style
)
```

**Available options:**
- `WithHeight(int)` - Set maximum height in lines. Box will be positioned at bottom of terminal (like fzf).
- `WithWidth(int)` - Set maximum width in columns. Box will be centered horizontally.
- `WithBorder()` - Enable border around the finder. Border is enabled by default.
- `WithBorderChars([]rune)` - Customize border characters `[topLeft, topRight, bottomLeft, bottomRight, horizontal, vertical]`.

**Note:** When `WithHeight()` is used, the finder box appears at the **bottom** of the terminal (similar to fzf), not centered. This leaves the top of your terminal free for output or context.

### Text editing and mouse support

The fuzzy finder includes modern text editing capabilities similar to fzf:

**Keyboard shortcuts:**
- `Ctrl+V` - Paste text from clipboard
- `Ctrl+X` - Cut selected text to clipboard
- `Ctrl+C` - Copy selected text to clipboard
- `Ctrl+Z` - Undo last text edit
- `Ctrl+D` or `Ctrl+Q` or `Esc` - Quit/abort the finder
- `Backspace` or `Delete` - Delete selected text (if any) or single character

**Mouse support:**
- **Left-click on prompt** - Position cursor at click location in search box
- **Left-click + drag on prompt** - Select text by dragging in search box
- **Double-click on prompt** - Select word under cursor in search box
- **Right-click on prompt** - Paste text from clipboard into search box
- **Left-click on item** - Select that item in the list
- **Mouse wheel up/down** - Scroll through the item list

Selected text is highlighted with inverted colors for visibility.

## Development

### Building the example application

To build the example application with the latest changes, run:

```bash
go build -o /tmp/cli example/cli/main.go
```

### Running the example application

The fuzzy finder requires an interactive terminal. When running these commands, ensure you are in a terminal where you can interact with the UI (type, use arrow keys, etc.). The error `open /dev/tty: no such device or address` indicates that the program is not running in an interactive terminal.

- **With border:**
  ```bash
  seq 1 10 | /tmp/cli --border
  ```
- **With height limit (e.g., 5 items):**
  ```bash
  seq 1 100 | /tmp/cli --height 5
  ```
- **With border and height limit:**
  ```bash
  seq 1 100 | /tmp/cli --border --height 5
  ```
- **With custom border characters (e.g., circular edges):**
  ```bash
  seq 1 10 | /tmp/cli --border --border-chars "╭╮╰╯─│"
  ```
  - **With custom border characters (e.g., long square edges):**
  ```bash
  seq 1 10 | /tmp/cli --border --border-chars "▛▜▙▟─│"
  ```

### Running tests

To run the project's tests, run:

```bash
go test ./...
```

### Manual testing

Test text editing and mouse support:

```bash
seq 1 100 | /tmp/cli --height 10 --border
```

**Keyboard:**
- Type to search, Ctrl+Z to undo
- Ctrl+C to copy, Ctrl+V to paste, Ctrl+X to cut
- Ctrl+D or Ctrl+Q to quit

**Mouse:**
- Click prompt to position cursor, drag to select
- Double-click prompt to select word
- Right-click prompt to paste
- Click items to select, wheel to scroll

## Motivation
Fuzzy-finder command-line tools such that
[fzf](https://github.com/junegunn/fzf), [fzy](https://github.com/jhawthorn/fzy), or [skim](https://github.com/lotabout/skim)
are very powerful to find out specified lines interactively.
However, there are limits to deal with fuzzy-finder's features in several cases.

First, it is hard to distinguish between two or more entities that have the same text.
In the example of ktr0731/itunes-cli, it is possible to conflict tracks such that same track names, but different artists.
To avoid such conflicts, we have to display the artist names with each track name.
It seems like the problem has been solved, but it still has the problem.
It is possible to conflict in case of same track names, same artists, but other albums, which each track belongs to.
This problem is difficult to solve because pipes and filters are row-based mechanisms, there are no ways to hold references that point list entities.

The second issue occurs in the case of incorporating a fuzzy-finder as one of the main features in a command-line tool such that [enhancd](https://github.com/b4b4r07/enhancd) or [itunes-cli](https://github.com/ktr0731/itunes-cli).
Usually, these tools require that it has been installed one fuzzy-finder as a precondition.
In addition, to deal with the fuzzy-finder, an environment variable configuration such that `export TOOL_NAME_FINDER=fzf` is required.
It is a bother and complicated.

`go-fuzzyfinder` resolves above issues.
Dealing with the first issue, `go-fuzzyfinder` provides the preview-window feature (See an example in [Usage](#usage)).
Also, by using `go-fuzzyfinder`, built tools don't require any fuzzy-finders.

## See Also
- [Fuzzy-finder as a Go library](https://medium.com/@ktr0731/fuzzy-finder-as-a-go-library-590b7458200f)
- [(Japanese) fzf ライクな fuzzy-finder を提供する Go ライブラリを書いた](https://syfm.hatenablog.com/entry/2019/02/09/120000)
