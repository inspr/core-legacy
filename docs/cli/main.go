package main

import (
	"log"
	"os"

	"github.com/inspr/inspr/cmd/insprctl/cli"
	"github.com/spf13/cobra/doc"
)

var version string

func main() {
	cmd := cli.NewInsprCommand(os.Stdout, os.Stderr, version)
	header := &doc.GenManHeader{
		Title: "Inspr CLI",
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	err = doc.GenManTree(cmd, header, path)
	if err != nil {
		log.Fatal(err)
	}

	err = doc.GenMarkdownTree(cmd, path)
	if err != nil {
		log.Fatal(err)
	}
}
