package server

import (
	"dnsserver/dns"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

func Server() {
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
	defer func(ln *net.UDPConn) {
		err := ln.Close()
		if err != nil {
		}
	}(ln)
	for {
		tmp := make([]byte, 1024)
		_, addr, _ := ln.ReadFromUDP(tmp)
		clientAddr := addr
		packet := gopacket.NewPacket(tmp, layers.LayerTypeDNS, gopacket.Default)
		dnsPacket := packet.Layer(layers.LayerTypeDNS)
		tcp, _ := dnsPacket.(*layers.DNS)
		serveDNS(ln, clientAddr, tcp)
	}
}

func serveDNS(u *net.UDPConn, clientAddr net.Addr, request *layers.DNS) {
	database := dns.LoadDatabase()
	replyMess := request
	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	var ip string
	var err error
	// Get IP Adress from JSON Database.
	ip = dns.GetIPAdress(string(request.Questions[0].Name), database)
	a, _, _ := net.ParseCIDR(ip + "/24")
	// SetUP Dns Answer
	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.IP = a
	dnsAnswer.Name = request.Questions[0].Name
	dnsAnswer.Class = layers.DNSClassIN
	// Setup ReplyMess Anwer
	replyMess.QR = true
	replyMess.ANCount = 1
	replyMess.OpCode = layers.DNSOpCodeNotify
	replyMess.AA = true
	replyMess.Answers = append(replyMess.Answers, dnsAnswer)
	replyMess.ResponseCode = layers.DNSResponseCodeNoErr
	// Convert to binary
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
	err = replyMess.SerializeTo(buf, opts)
	if err != nil {
		panic(err)
	}
	// Send Answer to DNS Client
	to, err := u.WriteTo(buf.Bytes(), clientAddr)
	if err != nil {
		return
	}
	_ = to
}
