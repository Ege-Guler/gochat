package chat

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Peer struct {
	RemoteAddr *net.UDPAddr
	Conn       *net.UDPConn
	SendQueue  chan []byte
}

func NewPeer(conn *net.UDPConn, raddr *net.UDPAddr) (*Peer, error) {

	return &Peer{
		RemoteAddr: raddr,
		Conn:       conn,
		SendQueue:  make(chan []byte, 100),
	}, nil
}

func (p *Peer) PunchHole() {

	msg := []byte("powpow")
	_, err := p.Conn.WriteToUDP(msg, p.RemoteAddr)
	if err != nil {
		fmt.Println("Punch failed", err)
	}
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

// !TODO ip base filtering
func (p *Peer) receiveLoop() {
	buf := make([]byte, 2048)
	for {
		n, _, err := p.Conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Recieve error:", err)
			continue
		}
		fmt.Printf("\r< %s\n", string(buf[:n]))

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
