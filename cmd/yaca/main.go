package main

import (
	"fmt"
	"os"

	"yaca/client"
	"yaca/models"
	"yaca/pkg/utils"
)

var (
	utilsLoadEnv                = utils.LoadEnv
	utilsParseArgs              = utils.ParseArgs
	utilsValidateArgs           = utils.ValidateArgs
	utilsPanicOnError           = utils.PanicOnError
	clientGetZoneIDByName       = client.GetZoneIDByName
	clientDoesRecordExistOnZone = client.DoesRecordExistOnZone
	clientUpdateRecordOnZone    = client.UpdateRecordOnZone
	clientCreateRecordOnZone    = client.CreateRecordOnZone
	clientDeleteRecordOnZone    = client.DeleteRecordOnZone
)

func run() int {
	utilsLoadEnv()

	args := utilsParseArgs()
	err := utilsValidateArgs(&args)
	utilsPanicOnError(err)

	zoneID, err := clientGetZoneIDByName(args.ZoneName)
	utilsPanicOnError(err)

	fmt.Printf("Zone ID: %+v\n", zoneID)

	recordID, err := clientDoesRecordExistOnZone(zoneID, args.Record)
	utilsPanicOnError(err)

	record := models.Record{
		Record: args.Record,
		Proxy:  args.Proxy,
		Target: args.Target,
		Ttl:    args.Ttl,
		Type:   args.Type,
	}

	if recordID != "" {
		fmt.Printf("Record %s exists on zone %s.\n", args.Record, args.ZoneName)

		if !args.Delete {
			success, err := clientUpdateRecordOnZone(zoneID, recordID, record)
			utilsPanicOnError(err)

			if success {
				fmt.Printf("Record %s updated successfully on zone %s.\n", args.Record, args.ZoneName)
				return 0
			}
		} else {
			success, err := clientDeleteRecordOnZone(zoneID, recordID, record)
			utilsPanicOnError(err)

			if success {
				fmt.Printf("Record %s deleted successfully from zone %s.\n", args.Record, args.ZoneName)
				return 0
			}
		}
	} else {
		fmt.Printf("Record %s does not exist on zone %s.\n", args.Record, args.ZoneName)

		success, err := clientCreateRecordOnZone(zoneID, record)
		utilsPanicOnError(err)

		if success {
			fmt.Printf("Record %s created successfully on zone %s.\n", args.Record, args.ZoneName)
			return 0
		}
	}

	return 1
}

func main() {
	os.Exit(run())
}
