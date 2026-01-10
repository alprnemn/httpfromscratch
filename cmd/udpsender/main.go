package main

import (
	"bufio"
	"fmt"
	cmn "httpfromscratch/common"
	"net"
	"os"
)

const networkName = "udp"

func main() {

	udpAddr, err := net.ResolveUDPAddr(networkName, "localhost:42069")
	cmn.FailOnError(err, "error resolving udp addr")

	conn, err := net.DialUDP(networkName, nil, udpAddr)
	cmn.FailOnError(err, "error dialin udp connection")
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		line, err := reader.ReadString('\n')
		cmn.FailOnError(err, "error reading line using reader readstring")

		_, err = conn.Write([]byte(line))
		cmn.FailOnError(err, "error writing []byte data")
	}
}
