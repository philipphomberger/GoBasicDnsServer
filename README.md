In this Project, I build for learning purposes a DNS Server with GoLang.
Set up a Socket Server and convert binary data back to protocol. 
Search for the Host in the JSON Database and return the IP Address with the DNS Protocol as binary.
Build and Push Docker Image.

Using google/gopacket package to not reimplement the DNS Protocol from Scretsh.

Learning:
- How DNS Works under the Hood.
- How to read DNS Protocol from binary.
- How to write a DNS Protocol answer to the requestor.
- Deeper understanding of How to set up a UDP Server in GoLang and how to send answers to the client.

Backround Links:
- https://pkg.go.dev/github.com/google/gopacket/layers#DNS
- https://blog.bytebytego.com/p/how-does-the-domain-name-system-dns
