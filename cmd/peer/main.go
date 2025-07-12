package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Ege-Guler/gochat/config"
	"github.com/Ege-Guler/gochat/internal/chat"
)

func main() {

	localAddr := config.GetPublicIp()

	laddr := config.GetLocalUDPAddr()

	fmt.Println("Your Public Ip: ", localAddr)
	fmt.Print("Remote addr: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	remoteAddr := config.ParseRemoteAddr(input)

	p, err := chat.NewPeerUDPAddr(laddr, remoteAddr)

	//cfg := config.ParseArgs()
	//p, err := chat.NewPeer(cfg.LocalAddr, cfg.RemoteAddr)

	if err != nil {

		log.Fatal("failed to start", err)
	}

	p.Start()
}
