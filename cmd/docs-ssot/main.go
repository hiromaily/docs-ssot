package main

// appVersion is set at build time via -ldflags "-X main.appVersion=<version>".
var appVersion = "dev"

func main() {
	Execute()
}
