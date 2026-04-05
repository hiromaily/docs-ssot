package main

import "github.com/hiromaily/docs-ssot/internal/cli"

// appVersion is set at build time via -ldflags "-X main.appVersion=<version>".
var appVersion = "dev"

func main() {
	cli.Execute(appVersion)
}
