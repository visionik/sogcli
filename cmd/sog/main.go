// sog - Standard Open (mail) CLI
//
// IMAP4/SMTP alternative to gog/mog for standard mail protocols.
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/visionik/sog/internal/cli"
)

var version = "dev"

func main() {
	var root cli.Root
	ctx := kong.Parse(&root,
		kong.Name("sog"),
		kong.Description("Standard Open (mail) CLI — IMAP4/SMTP for everyone"),
		kong.UsageOnError(),
		kong.Vars{
			"version": version,
		},
	)

	err := ctx.Run(&root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
