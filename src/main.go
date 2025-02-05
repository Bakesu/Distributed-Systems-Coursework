package main

import (
	"bufio"
	"dissy2021-group1/Hand-in9/peer"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please provide IP adress to connect to \n >")
	ip, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	fmt.Print("Please provide port number for the connection \n >")
	port, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	peer := peer.MkPeer(ip, port)
	fmt.Println(reflect.TypeOf(peer))
	for {
		fmt.Print(">")
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		if strings.TrimSpace(msg) == "quit" {
			return
		}
		// envelope := peer.MkEnvelope()
		// peer.AddEnvelopeToEnvelopeChannel(envelope)
		time.Sleep(500 * time.Millisecond)
	}
}
