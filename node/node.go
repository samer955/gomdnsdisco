package node

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"log"
	"net"
)

const discoveryTag = "discoveryRoom"

type Node struct {
	Ip   string
	Host host.Host
	mdns Mdns
}

func InitializeNode(ctx context.Context) *Node {

	n := new(Node)
	n.createLibp2pHost()
	n.getIp()
	go n.findPeers(ctx)

	return n

}

// initialize Node using Libp2p, listening all ip4 address and default tcp port
func (n *Node) createLibp2pHost() {

	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}
	n.Host = host
	log.Printf("Hello LAN, my hosts ID is %s\n", n.Host.ID().ShortString())

}

func (n *Node) getIp() {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		n.Ip = ""
		return
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				n.Ip = ipnet.IP.String()
				return
			}
		}
	}
	n.Ip = ""

}

func (n *Node) findPeers(ctx context.Context) {

	peerChan := n.mdns.initMDNS(n.Host, discoveryTag)
	for {
		peer := <-peerChan
		fmt.Println("Found peer:", peer.ID.ShortString(), ", connecting")

		if err := n.Host.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}
	}
}
