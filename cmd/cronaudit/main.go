package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/cronaudit/internal/formatter"
	"github.com/yourorg/cronaudit/internal/loader"
	"github.com/yourorg/cronaudit/internal/reporter"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "cronaudit: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		format      = flag.String("format", "text", "output format: text, markdown, json")
		invalidOnly = flag.Bool("invalid", false, "show only invalid entries")
		validOnly   = flag.Bool("valid", false, "show only valid entries")
	)
	flag.Parse()

	if *invalidOnly && *validOnly {
		return fmt.Errorf("flags -invalid and -valid are mutually exclusive")
	}

	var src *loader.Source
	var err error

	if flag.NArg() > 0 {
		src, err = loader.FromFile(flag.Arg(0))
	} else {
		src, err = loader.FromStdin()
	}
	if err != nil {
		return err
	}
	defer src.Close()

	lines, err := loader.Lines(src)
	if err != nil {
		return err
	}

	rpt := reporter.Build(lines, src.Name)

	switch {
	case *invalidOnly:
		rpt = reporter.FilterInvalid(rpt)
	case *validOnly:
		rpt = reporter.FilterValid(rpt)
	}

	return formatter.Render(os.Stdout, rpt, *format)
}
