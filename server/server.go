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
		data := []byte(dns.GetIPAdress(m.Questions[0].Name.String(), database))

		// Questions Header
		var questionAnswer = layers.DNSQuestion{
			Type:  layers.DNSTypeA,
			Class: layers.DNSClassAny,
			Name:  []byte(m.Questions[0].Name.String()),
		}

		var questionAnswerArray []layers.DNSQuestion
		questionAnswerArray = append(questionAnswerArray, questionAnswer)

		// Combine header, question, and record into one response message
		var dnsAnswer layers.DNSResourceRecord
		dnsAnswer.Type = layers.DNSTypeA
		ipadress := net.ParseIP(string(data))
		responseanswer := layers.DNSResourceRecord{
			Name:  []byte(m.Questions[0].Name.String()),
			Type:  layers.DNSTypeA,
			Class: layers.DNSClassIN,
			TTL:   0,
			Data:  []byte(m.Questions[0].Name.String()),
			IP:    ipadress,
		}

		var dnsAnswerArray []layers.DNSResourceRecord
		dnsAnswerArray = append(dnsAnswerArray, responseanswer)

		dnsAnswer.IP = ipadress
		dnsAnswer.Name = []byte(m.Questions[0].Name.String())
		dnsAnswer.Class = layers.DNSClassIN
		dnsAnswer.TTL = 300 // Set TTL to a reasonable value, e.g., 300 seconds

		var replyMess layers.DNS
		replyMess.ID = m.ID
		replyMess.QR = true
		replyMess.OpCode = layers.DNSOpCodeQuery // Use Query opcode for standard query response
		replyMess.AA = true
		replyMess.RD = true
		replyMess.RA = false // Assuming the resolver is recursive
		replyMess.ResponseCode = layers.DNSResponseCodeNoErr
		replyMess.QDCount = 1
		replyMess.ANCount = 1
		replyMess.Questions = questionAnswerArray
		replyMess.Answers = dnsAnswerArray
		replyMess.Additionals = dnsAnswerArray
		replyMess.Authorities = dnsAnswerArray

		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		} // See SerializeOptions for more details.
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
