package dns

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type dnsentry struct {
	Dns string `json:"dns"`
	Ip  string `json:"ip"`
}

type Alldnsentry []dnsentry

var dnsdatabase Alldnsentry

// DNSHeader represents the header of a DNS packet
type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

type DNSQuestion struct {
	QName  string
	QType  uint16
	QClass uint16
}

// DNSRecord represents a DNS resource record
type DNSRecord struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32
	Data  net.IP
}

// Convert the DNSHeader to a byte slice
func (h *DNSHeader) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, h)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Convert the DNSRecord to a byte slice
func (r *DNSRecord) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Convert the Name to the DNS format
	parts := strings.Split(r.Name, ".")
	for _, part := range parts {
		buf.WriteByte(byte(len(part)))
		buf.WriteString(part)
	}
	buf.WriteByte(0) // End of the Name

	err := binary.Write(buf, binary.BigEndian, r.Type)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, r.Class)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, r.TTL)
	if err != nil {
		return nil, err
	}

	// Write the data length (IPv4 address is 4 bytes)
	buf.WriteByte(0)
	buf.WriteByte(4)

	// Write the IP address
	ip := r.Data.To4()
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address")
	}
	_, err = buf.Write(ip)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Convert the DNSQuestion to a byte slice
func (q *DNSQuestion) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Convert the QName to the DNS format
	parts := strings.Split(q.QName, ".")
	for _, part := range parts {
		buf.WriteByte(byte(len(part)))
		buf.WriteString(part)
	}
	buf.WriteByte(0) // End of the QName

	err := binary.Write(buf, binary.BigEndian, q.QType)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, q.QClass)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func AddDnsEntry(dnsdatabase Alldnsentry, dns string, ip string) Alldnsentry {
	dnsdatabase = append(dnsdatabase, dnsentry{
		Dns: dns,
		Ip:  ip,
	})
	return dnsdatabase
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
	fmt.Print(configuration)
	return configuration
}

func GetIPAdress(dns string, database Alldnsentry) string {
	for _, entry := range database {
		println(entry.Dns)
		println(entry.Ip)
		println(dns)
		if entry.Dns == dns {
			return entry.Ip
		}
	}
	return "Not exist!"
}
