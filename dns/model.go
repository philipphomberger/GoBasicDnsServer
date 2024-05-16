package dns

type dnsentry struct {
	Dns string `json:"dns"`
	Ip  string `json:"ip"`
}

type Alldnsentry []dnsentry

var dnsdatabase Alldnsentry
