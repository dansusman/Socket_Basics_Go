package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	port, encrypt, host, neu := readArgs()
	CONNECT := fmt.Sprintf("%s:%d", host, port)

	nonTlsConn, tlsConn, err := makeConn(encrypt, CONNECT)

	checkError(err)

	var connection net.Conn

	if nonTlsConn != nil {
		// -s flag not set, non TLS communication desired
		connection = nonTlsConn
	} else {
		connection = tlsConn
	}
	defer connection.Close()

	helloMessage := "ex_string HELLO " + neu + "\n"

	writeToServer(connection, helloMessage)

	response, readError := readFromServer(connection)
	// loop until we find a BYE message
	for !strings.Contains(response, "BYE") {
		fmt.Println(response)
		checkError(readError)
		verifyResponse(response)
		count := countOccurrence(response)
		countMessage := "ex_string COUNT " + strconv.Itoa(count) + "\n"
		fmt.Println(countMessage)
		writeToServer(connection, countMessage)
		response, readError = readFromServer(connection)
	}
	fmt.Println(response)
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

func readFromServer(connection net.Conn) (string, error) {
	reader := bufio.NewReader(connection)
	var buff bytes.Buffer
	for {
		// read a single line because we are guaranteed to see a '\n'
		// at the end of a valid server response
		line, isPrefix, readError := reader.ReadLine()
		if readError != nil {
			if readError == io.EOF {
				break
			} else {
				// error is something other than EOF: exit program
				return "", readError
			}
		}

		// write to the buffer
		buff.Write(line)

		// From the docs: "If the line was too long for the buffer
		// then isPrefix is set and the beginning of the line is returned.
		// The rest of the line will be returned from future calls.
		// isPrefix will be false when returning the last fragment
		// of the line."
		// We are sure that the whole message has been read iff isPrefix
		// is false, so break out of the loop
		if !isPrefix {
			break
		}
	}
	return buff.String(), nil
}

func writeToServer(connection net.Conn, data string) {
	_, writeError := connection.Write([]byte(data))
	checkError(writeError)
}

func countOccurrence(response string) int {
	stringArr := strings.Split(response, " ")
	return strings.Count(stringArr[3], stringArr[2])
}

func checkError(err error) {
	if err != nil {
		panic("failed in comm with server; reason: " + err.Error())
	}
}

func verifyResponse(response string) {
	stringArr := strings.Split(response, " ")
	if len(stringArr) < 2 || stringArr[0] != "ex_string" || invalidCommand(stringArr[1], stringArr) {
		panic("Response does not conform to the protocol! " + response)
	}
}

func invalidCommand(command string, array []string) bool {
	switch command {
	case "HELLO":
		return len(array) != 3
	case "FIND":
		return len(array) != 4
	case "COUNT":
		return len(array) != 3
	case "BYE":
		return len(array) != 3
	default:
		return true
	}
}
