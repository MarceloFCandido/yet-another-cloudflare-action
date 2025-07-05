package models

type Args struct {
	Record string `arg:"required,-r,--record" name:"Record" help: "Record name to be created/updated"`
	Proxy bool `arg:"-p,--proxy" name:"Proxy" help: "Whether to enable Cloudflare proxy for the record name"`
	Target string `arg:"required,-t,--target" name:"Target" help: "Target/IP address the record name should point to"`
	Ttl int `arg:"-l,--ttl" name:"TTL" help: "Time-to-live for the record name" default:"5"`
	Type string `arg:"-y,--type" name:"Type" help: "Type of the record name to be created/updated" default:"A"`
	ZoneName string `arg:"required,-z,--zone-name" name:"ZoneName" help: "Zone name of the record name"`
}
