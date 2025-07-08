package utils

import (
	"fmt"

	"yaca/models"
)

func ValidateArgs(args *models.Args) error {
	if args.Record == "" {
		return fmt.Errorf("record is required")
	}
	if args.ZoneName == "" {
		return fmt.Errorf("zone is required")
	}

	if args.Delete {
		if args.Target != "" || args.Type != "" || args.Proxy || args.Ttl != 0 {
			return fmt.Errorf("all the arguments, except for record and zone name, must be empty when delete is true")
		}
		return nil
	}

	if args.Target == "" {
		return fmt.Errorf("target is required")
	}
	if args.Type == "" {
		return fmt.Errorf("type is required")
	}

	return nil
}
