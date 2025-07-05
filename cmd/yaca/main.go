package main

import (
	"yaca/client"
	"yaca/pkg/utils"
)

func main() {
  utils.LoadEnv()

  args := utils.ParseArgs()

  client := client.GetSingletonClient()

  // zoneID, err = client.GetZoneIDByName(args.ZoneName)

  // page, err := client.Zones.List(context.TODO(), zones.ZoneListParams{
  //   Name: cloudflare.F(args.ZoneName),
  // })
	// if err != nil {
	// 	panic(err.Error())
	// }

	// fmt.Printf("%+v\n", page.Result[0].ID)

  // GetZoneID()
  // DoesDomainExistOnZone()
  // if true
  // UpdateDomain()
  // else
  // CreateDomain()
}
