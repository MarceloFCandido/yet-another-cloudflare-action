package main

import (
	"yaca/client"
	"yaca/pkg/utils"
)

func main() {
  utils.LoadEnv()

  args := utils.ParseArgs()

  client := client.GetClient()

  // GetZoneID()
  // DoesDomainExistOnZone()
  // if true
  // UpdateDomain()
  // else
  // CreateDomain()
}
