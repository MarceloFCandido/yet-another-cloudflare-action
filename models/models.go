package models

type Args struct {
	Delete	 bool    `arg:"-d,--delete" name:"Delete" help: "Whether to delete the record name" default:"false"`
	Record   string  `arg:"required,-r,--record" name:"Record" help: "Record name to be created/updated"`
	Proxy    bool    `arg:"-p,--proxy" name:"Proxy" help: "Whether to enable Cloudflare proxy for the record name" default:"false"`
	Target   string  `arg:"-t,--target" name:"Target" help: "Target/IP address the record name should point to"`
	Ttl      float64 `arg:"-l,--ttl" name:"TTL" help: "Time-to-live for the record name" default:"3600"`
	Type     string  `arg:"-y,--type" name:"Type" help: "Type of the record name to be created/updated"`
	ZoneName string  `arg:"required,-z,--zone-name" name:"ZoneName" help: "Zone name of the record name"`
}

type Record struct {
	Record string  `json:"name"`
	Proxy  bool    `json:"proxied"`
	Target string  `json:"content"`
	Ttl    float64 `json:"ttl"`
	Type   string  `json:"type"`
}

type RecordData struct {
	ZoneID   string
	RecordID string
	Record   Record
}
