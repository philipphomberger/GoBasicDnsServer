package dns

import (
	"github.com/google/gopacket/layers"
	"net"
)

// SetUP Dns Answer
func getDnsAnswer(a net.IP, request *layers.DNS) layers.DNSResourceRecord {
	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.IP = a
	dnsAnswer.Name = request.Questions[0].Name
	dnsAnswer.Class = layers.DNSClassIN
	return dnsAnswer
}

func ReplyDnsAnswer(a net.IP, replyMess *layers.DNS) *layers.DNS {
	replyMess.QR = true
	replyMess.ANCount = 1
	replyMess.OpCode = layers.DNSOpCodeNotify
	replyMess.AA = true
	replyMess.Answers = append(replyMess.Answers, getDnsAnswer(a, replyMess))
	replyMess.ResponseCode = layers.DNSResponseCodeNoErr
	return replyMess
}

func LoadDatabase() []dnsentry {
	file, _ := os.Open("database.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	var configuration Alldnsentry
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}

func GetIPAdress(dns string, database Alldnsentry) string {
	for _, entry := range database {
		if entry.Dns == dns {
			return entry.Ip
		}
	}
	return "Not exist!"
}
