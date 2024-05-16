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
	var err error
	a, _, _ := net.ParseCIDR(dns.GetIPAdress(string(request.Questions[0].Name), dns.LoadDatabase()) + "/24")
	fmt.Println(a.String())
	if a != nil {
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
		err = dns.ReplyDnsAnswer(a, request).SerializeTo(buf, opts)
		if err != nil {
			panic(err)
		}
		// Send Answer to DNS Client
		to, err := u.WriteTo(buf.Bytes(), clientAddr)
		if err != nil {
			return
		}
		_ = to
	} else {
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{}
		err = dns.ReplyDnsAnswerNotFound(a, request).SerializeTo(buf, opts)
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
}
