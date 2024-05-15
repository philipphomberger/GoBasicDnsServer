package server

import (
	"dnsserver/dns"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

func Server() {
	database := dns.LoadDatabase()
	fmt.Print(database)
	s, err := net.ResolveUDPAddr("udp", ":8090")
	if err != nil {
		fmt.Println(err)
		return
	}

	ln, err := net.ListenUDP("udp", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()
	buffer := make([]byte, 512)
	// Accept incoming connections and handle them
	for {
		_, addr, err := ln.ReadFromUDP(buffer)
		var m dnsmessage.Message
		err = m.Unpack(buffer)
		fmt.Println(m)
		fmt.Println("-> ", m.Questions[0].Name.String())

		data := []byte(dns.GetIPAdress(m.Questions[0].Name.String(), database))
		fmt.Printf("data: %s\n", string(data))

		// Combine header, question, and record into one response message
		var dnsAnswer layers.DNSResourceRecord
		dnsAnswer.Type = layers.DNSTypeA
		var ip string
		ip = m.Questions[0].Name.String()
		a, _, _ := net.ParseCIDR(ip + "/24")
		dnsAnswer.Type = layers.DNSTypeA
		dnsAnswer.IP = a
		dnsAnswer.Name = []byte(m.Questions[0].Name.String())
		fmt.Println(m.Questions[0].Name)
		dnsAnswer.Class = layers.DNSClassIN
		var replyMess layers.DNS
		replyMess.QR = true
		replyMess.ANCount = 1
		replyMess.OpCode = layers.DNSOpCodeNotify
		replyMess.AA = true
		replyMess.Answers = append(replyMess.Answers, dnsAnswer)
		replyMess.ID = m.ID
		replyMess.ResponseCode = layers.DNSResponseCodeNoErr
		fmt.Print(replyMess.Answers[0].Name)
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
		err = replyMess.SerializeTo(buf, opts)
		if err != nil {
			panic(err)
		}
		_, err = ln.WriteTo(buf.Bytes(), addr)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}
