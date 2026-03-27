package main

import (
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/nizarmah/nitpicker/patch"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: nitpicker <patch-url>\n")
		os.Exit(1)
	}
	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching patch: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "error fetching patch: HTTP %d\n", resp.StatusCode)
		os.Exit(1)
	}

	stats, err := patch.Parse(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing patch: %v\n", err)
		os.Exit(1)
	}

	if len(stats) == 0 {
		fmt.Println("no files changed")
		return
	}

	totalAdded := 0
	totalDeleted := 0

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _, s := range stats {
		fmt.Fprintf(w, "%s\t+%d\t-%d\n", s.Path, s.Added, s.Deleted)
		totalAdded += s.Added
		totalDeleted += s.Deleted
	}
	w.Flush()

	fmt.Printf("\n%d files changed, %d insertions(+), %d deletions(-)\n", len(stats), totalAdded, totalDeleted)
}
