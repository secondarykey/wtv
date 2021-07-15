package main

import (
	"fmt"
	"os"
	"wtv"

	"golang.org/x/xerrors"
)

func main() {

	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "wtv error:\n%+v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Bye!")
}

func run() error {

	err := wtv.Show()
	if err != nil {
		return xerrors.Errorf("wtv.Show() error: %w", err)
	}

	return nil
}
