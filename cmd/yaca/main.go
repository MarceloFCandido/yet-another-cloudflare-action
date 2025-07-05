package main

import (
	"fmt"
	"yaca/client"
	"yaca/pkg/utils"
)

func main() {
  utils.LoadEnv()

  args := utils.ParseArgs()

  zoneID, err := client.GetZoneIDByName(args.ZoneName)
  if err != nil {
    panic(err)
  }

	fmt.Printf("Zone ID: %+v\n", zoneID)

  // DoesDomainExistOnZone()
  // if true
  // UpdateDomain()
  // else
  // CreateDomain()
}
