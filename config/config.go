package config

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pion/stun"
)

func GetPublicIPAndConn() (*net.UDPConn, *net.UDPAddr, error) {
	stunServerAddr := "stun.l.google.com:19302"

	// unconnected UDP socket to listen for messages
	// listen on all network ifaces and a random port(0)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, nil, fmt.Errorf("listening on udp port failed: %w", err)
	}

	// Manually perform the STUN request.
	raddr, err := net.ResolveUDPAddr("udp", stunServerAddr)
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("resolving stun server address failed: %w", err)
	}

	// Build the STUN Binding Request message.
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	if _, err := conn.WriteToUDP(message.Raw, raddr); err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("sending stun request failed: %w", err)
	}

	// buffer to receive the response.
	buf := make([]byte, 1024)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("reading stun response failed: %w", err)
	}

	conn.SetReadDeadline(time.Time{})

	var response stun.Message
	response.Raw = buf[:n]
	if err := response.Decode(); err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("decoding stun response failed: %w", err)
	}

	var publicAddr stun.XORMappedAddress
	if err := publicAddr.GetFrom(&response); err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("getting public address from response failed: %w", err)
	}

	resolvedAddr := &net.UDPAddr{
		IP:   publicAddr.IP,
		Port: publicAddr.Port,
	}

	return conn, resolvedAddr, nil
}

func ParseRemoteAddr(raddr string) *net.UDPAddr {
	udpAddr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		fmt.Printf("Invalid remote address format '%s'. Please use 'ip:port'. Error: %v\n", raddr, err)
		os.Exit(1)
	}
	return udpAddr
}
