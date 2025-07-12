package config

import (
	"fmt"
	"net"
	"os"

	"github.com/pion/stun"
)

type Config struct {
	LocalAddr  string
	RemoteAddr string
}

/*
Parse RemoteAddr and LocalAddr

if LocalAddr is not present get local ip

!TODO improve checks & better parsing
*/
func ParseArgs() *Config {

	cfg := &Config{}

	if len(os.Args) <= 1 {
		fmt.Println("unsufficient args")
		os.Exit(1)
	} else if len(os.Args) == 3 {
		cfg.LocalAddr = os.Args[1]

	} else if len(os.Args) == 2 {
		cfg.LocalAddr = getLocalIp()
	}

	cfg.RemoteAddr = os.Args[len(os.Args)-1]

	return cfg
}

func ParseRemoteAddr(raddr string) *net.UDPAddr {
	udpAddr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		panic(err)
	}
	return udpAddr
}

func getLocalIp() string {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	fmt.Println("localip", conn.LocalAddr().String())

	defer conn.Close()

	return conn.LocalAddr().String()
}

func GetLocalUDPAddr() *net.UDPAddr {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	fmt.Println("localip", conn.LocalAddr().String())

	defer conn.Close()

	udpAddr := conn.LocalAddr().(*net.UDPAddr)

	return udpAddr
}

func GetPublicIp() *net.UDPAddr {
	conn, _ := net.Dial("udp", "stun.l.google.com:19302")
	defer conn.Close()

	c, err := stun.NewClient(conn)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	var addr stun.XORMappedAddress
	c.Do(stun.MustBuild(stun.TransactionID, stun.BindingRequest), func(e stun.Event) {
		if err := addr.GetFrom(e.Message); err != nil {
			panic(err)
		}
	})
	return &net.UDPAddr{
		IP:   addr.IP,
		Port: addr.Port,
	}
}
