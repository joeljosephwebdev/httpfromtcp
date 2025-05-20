package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to %s", udpAddr.String())
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	// main connection loop
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF { // handle graceful close for ctrl+D
				log.Println("EOF received, exiting.")
				break
			}
			log.Printf("Error reading input: %v", err)
			continue
		}
		if line == "exit\n" { // handle graceful close for exit
			break
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("Error sending data: %v", err)
			continue
		}
	}

}
