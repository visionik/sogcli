// sog - Standards Ops Gadget
//
// Open-standards CLI for mail, calendar, contacts, tasks, and files.
// IMAP/SMTP/CalDAV/CardDAV/WebDAV alternative to gog (Google) and mog (Microsoft).
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/visionik/sogcli/internal/cli"
)

var version = "dev"

func main() {
	var root cli.Root
	ctx := kong.Parse(&root,
		kong.Name("sog"),
		kong.Description("Standards Ops Gadget — IMAP/SMTP/CalDAV/CardDAV/WebDAV CLI"),
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
