package utils

import (
	"yaca/models"

	"github.com/alexflint/go-arg"
)

func ParseArgs() (models.Args) {
	var args models.Args

	arg.MustParse(&args)
	
	return args
}
