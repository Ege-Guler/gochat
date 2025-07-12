package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Ege-Guler/gochat/config"
	"github.com/Ege-Guler/gochat/internal/chat"
)

func main() {

	conn, publicAddr, err := config.GetPublicIPAndConn()
	if err != nil {
		log.Fatalf("Failed to get public IP: %v", err)
	}

	fmt.Println("Your Public IP: ", publicAddr.String())
	fmt.Println("Your Local IP: ", conn.LocalAddr().String())
	fmt.Print("Remote addr: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	remoteAddr := config.ParseRemoteAddr(input)

	p, err := chat.NewPeer(conn, remoteAddr)

	if err != nil {

		log.Fatalf("failed to start %v", err)
	}

	p.PunchHole()

	time.Sleep(500 * time.Millisecond)

	p.Start()
}
