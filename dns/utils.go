package dnsclient

import (
	"encoding/json"
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
	"os"
	"os/exec"
	"strings"
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

func ReplyDnsAnswerNotFound(a net.IP, replyMess *layers.DNS) *layers.DNS {
	replyMess.QR = true
	replyMess.ANCount = 1
	replyMess.OpCode = layers.DNSOpCodeNotify
	replyMess.AA = true
	replyMess.Answers = append(replyMess.Answers, getDnsAnswer(a, replyMess))
	replyMess.ResponseCode = layers.DNSResponseCodeNXDomain
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
	var google string
	google = GetIPFromGoogle(dns)
	if google != "" {
		return google
	} else {
		return ""
	}
}

func GetIPFromGoogle(domain string) string {
	var dnsList []string
	dnsServer := "8.8.8.8"
	cmd := exec.Command("nslookup", domain, dnsServer)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Address: ") {
			fmt.Println(line)
			dnsList = append(dnsList, line)
		}
	}
	var ip string
	if len(dnsList) == 0 {
		ip = ""
		return ip
	}
	ip = dnsList[0]
	return ip[9:len(ip)]
}
