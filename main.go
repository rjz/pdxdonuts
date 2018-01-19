package main

import (
	"flag"
	"os"
)

var (
	optCli  = flag.Bool("cli", false, "Invoke via CLI")
	optPort = flag.String("port", ":8345", "HTTP port")
)

func main() {
	flag.Parse()

	if *optCli {
		cli()
		os.Exit(0)
	}
	serve(*optPort)
}
