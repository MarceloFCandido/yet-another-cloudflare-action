package utils

import (
	"os"
	"yaca/models"

	"github.com/alexflint/go-arg"
)

var ParseArgs = parseArgs

func parseArgs() (models.Args) {
	var args models.Args
	p, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		// This error should not happen with empty config and valid args struct
		panic(err)
	}
	err = p.Parse(os.Args[1:])
	if err != nil {
		// Only exit if there's a real parsing error, not just empty args
		if err == arg.ErrHelp {
			p.WriteHelp(os.Stderr)
			os.Exit(0) // Exit with 0 for help
		}
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	return args
}
