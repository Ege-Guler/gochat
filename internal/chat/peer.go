package chat

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Peer struct {
	LocalAddr  *net.UDPAddr
	RemoteAddr *net.UDPAddr
	Conn       *net.UDPConn
	SendQueue  chan []byte
}

func NewPeer(local string, remote string) (*Peer, error) {

	laddr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		return nil, err
	}

	raddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}

	return &Peer{
		LocalAddr:  laddr,
		RemoteAddr: raddr,
		Conn:       conn,
		SendQueue:  make(chan []byte, 100),
	}, nil
}

func NewPeerUDPAddr(laddr *net.UDPAddr, raddr *net.UDPAddr) (*Peer, error) {
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}

	return &Peer{
		LocalAddr:  laddr,
		RemoteAddr: raddr,
		Conn:       conn,
		SendQueue:  make(chan []byte, 100),
	}, nil

}

func (p *Peer) Start() {
	go p.sendLoop()
	go p.receiveLoop()
	p.inputLoop()
}

func (p *Peer) sendLoop() {
	for msg := range p.SendQueue {
		_, err := p.Conn.WriteToUDP(msg, p.RemoteAddr)
		if err != nil {
			fmt.Println("Send error:", err)
		}
	}
}

func (p *Peer) receiveLoop() {
	buf := make([]byte, 2048)
	for {
		n, addr, err := p.Conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Recieve error:", err)
			continue
		}
		if addr.String() == p.RemoteAddr.String() {
			fmt.Printf("\r< %s\n", string(buf[:n]))
		}
	}
}

func (p *Peer) inputLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			text := scanner.Text()
			p.SendQueue <- []byte(text)
		}
	}
}
