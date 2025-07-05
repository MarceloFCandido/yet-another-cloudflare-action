package main

import (
	"fmt"
	"yaca/client"
	"yaca/models"
	"yaca/pkg/utils"
)

func main() {
	utils.LoadEnv()

	args := utils.ParseArgs()

	zoneID, err := client.GetZoneIDByName(args.ZoneName)
	utils.PanicOnError(err)

	fmt.Printf("Zone ID: %+v\n", zoneID)

	recordID, err := client.DoesRecordExistOnZone(zoneID, args.Record)
	utils.PanicOnError(err)

  record := models.Record{
    Record: args.Record,
    Proxy:  args.Proxy,
    Target: args.Target,
    Ttl:    args.Ttl,
    Type:   args.Type,
  }

	if recordID != "" {
		fmt.Printf("Record %s already exists on zone %s.\n", args.Record, args.ZoneName)

    success, err := client.UpdateRecordOnZone(zoneID, recordID, record)
    utils.PanicOnError(err)

    if success {
      fmt.Printf("Record %s updated successfully on zone %s.\n", args.Record, args.ZoneName)
    }
	} else {
		fmt.Printf("Record %s does not exist on zone %s.\n", args.Record, args.ZoneName)

		success, err := client.CreateRecordOnZone(zoneID, record)
    utils.PanicOnError(err)

    if success {
      fmt.Printf("Record %s created successfully on zone %s.\n", args.Record, args.ZoneName)
    }
	}
}
