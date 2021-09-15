package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	port, encrypt, host, neu := readArgs()
	CONNECT := fmt.Sprintf("%s:%d", host, port)

	nonTlsConn, tlsConn, err := makeConn(encrypt, CONNECT)

	if err != nil {
		// an error has occurred
		panic("failed to connect; reason: " + err.Error())
	}

	var connection net.Conn

	if nonTlsConn != nil {
		// -s flag not set, non TLS communication desired
		connection = nonTlsConn
	} else {
		connection = tlsConn
	}

	helloResponse := clientHello(connection, neu)
	fmt.Println(helloResponse)

	count := countOccurrence(helloResponse)
	countCall := fmt.Sprintf("%s %d %s", "ex_string COUNT", count, "\n")
	_, writeError := connection.Write([]byte(countCall))
	fmt.Println(count)
	if writeError != nil {
		// an error has occurred in writing COUNT
		panic("failed to write to server; reason: " + writeError.Error())
	}
	status, readError := bufio.NewReader(connection).ReadString('\n')

	if readError != nil {
		// an error has occurred in reading HELLO response
		panic("failed to read server response; reason: " + readError.Error())
	}
	fmt.Println("after: ", status)
}

// Reads in command-line inputs given to client program.
// Optional flags:
// -p: Specifies a port to listen to at the given hostname
// -s: Specifies TLS encryption true/false
// Required arguments:
// hostname: name of the server (either a DNS name or an IP address in dotted notation)
// NEU ID: a valid Northeastern ID
func readArgs() (int, bool, string, string) {
	args := os.Args
	if len(args) < 3 {
		panic("Please provide the hostname and your NEU ID!")
	}

	portPtr := flag.Int("p", 0, "port number")
	tlsPtr := flag.Bool("s", false, "TLS encryption")
	flag.Parse()

	if *portPtr == 0 {
		if *tlsPtr {
			// port not specified, TLS encrpytion is desired
			*portPtr = 27994
		} else {
			*portPtr = 27993
		}
	}

	args = flag.Args()

	hostname := args[0]
	neuId := args[1]

	return *portPtr, *tlsPtr, hostname, neuId

}

func makeConn(encrypt bool, CONNECT string) (net.Conn, *tls.Conn, error) {
	if encrypt {
		conn, err := tls.Dial("tcp", CONNECT, &tls.Config{})
		return nil, conn, err
	}
	connection, err := net.Dial("tcp", CONNECT)
	return connection, nil, err
}

func clientHello(connection net.Conn, neu string) string {
	_, writeError := connection.Write([]byte("ex_string HELLO " + neu + "\n"))
	if writeError != nil {
		// an error has occurred in writing HELLO
		panic("failed to write to server; reason: " + writeError.Error())
	}

	readBuf := make([]byte, 8192)
	n, readError := connection.Read(readBuf)
	// status, readError := bufio.NewReader(connection).ReadString('\n')

	if readError != nil {
		// an error has occurred in reading HELLO response
		panic("failed to read server response; reason: " + readError.Error())
	}
	return string(readBuf[:n])
	// return status
}

func countOccurrence(response string) int {
	stringArr := strings.Split(response, " ")
	return strings.Count(stringArr[3], stringArr[2])
}
