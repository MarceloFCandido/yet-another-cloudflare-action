package models

type Args struct {
	Domain string `arg:"required,-d,--domain" name:"Domain" help: "Domain name to be created/updated"`
	Proxy bool `arg:"-p,--proxy" name:"Proxy" help: "Whether to enable Cloudflare proxy for the domain name"`
	Target string `arg:"required,-t,--target" name:"Target" help: "Target/IP address the domain name should point to"`
	Ttl int `arg:"-l,--ttl" name:"TTL" help: "Time-to-live for the domain name" default:"5"`
	Type string `arg:"-y,--type" name:"Type" help: "Type of the domain name to be created/updated" default:"A"`
	ZoneName string `arg:"required,-z,--zone-name" name:"ZoneName" help: "Zone name of the domain name"`
}
