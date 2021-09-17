package main

import (
	"bufio"
	// "bytes"
	"crypto/tls"
	"flag"
	"fmt"
	// "io"
	"net"
	"os"
	"strconv"
	"strings"
)

// Main method to run the TCP client program.
// Reads inputs from command line, initiates a connection
// with the server hostname:port specified in command line
// arguments, sends encrypted or non-encrypted HELLO message to server,
// reads response, returns COUNT messages until server responds with BYE.
func main() {

	// parse command line arguments
	port, encrypt, host, neu := readArgs()

	// generate encrypted or non-encrypted connection
	CONNECT := fmt.Sprintf("%s:%d", host, port)
	nonTlsConn, tlsConn, err := makeConn(encrypt, CONNECT)
	checkError(err)

	// cast connection to TLS or non-TLS (*tls.Conn inherits from net.Conn)
	var connection net.Conn
	if nonTlsConn != nil {
		// -s flag not set, non TLS communication desired
		connection = nonTlsConn
	} else {
		connection = tlsConn
	}

	// wait until function returns to close connection with TCP server
	defer connection.Close()

	// write HELLO to server
	helloMessage := "ex_string HELLO " + neu + "\n"
	writeToServer(connection, helloMessage)

	// grab first FIND message from server response
	response, readError := readFromServer(connection)

	// loop until we find a BYE message
	for strings.Split(response, " ")[1] != "BYE" {
		checkError(readError)
		verifyResponse(response)
		count := countOccurrence(response)
		countMessage := "ex_string COUNT " + strconv.Itoa(count) + "\n"
		writeToServer(connection, countMessage)
		response, readError = readFromServer(connection)
	}
	// print secret flag to console
	fmt.Println(strings.Split(response, " ")[2])
}

// Reads in command-line inputs given to client program.
// Optional flags:
// -p: Specifies a port to listen to at the given hostname
// -s: Specifies TLS encryption true/false
// Required arguments:
// hostname: name of the server (either a DNS name or an IP address in dotted notation)
// NEU ID: a valid Northeastern ID
// Returns the port number, TLS on/off boolean, hostname for server, and NEU ID
// in that order.
func readArgs() (int, bool, string, string) {
	args := os.Args
	if len(args) < 3 {
		panic("Please provide the hostname and your NEU ID!")
	}

	// default value for port = 0 to be used later based on -s flag
	portPtr := flag.Int("p", 0, "port number")
	tlsPtr := flag.Bool("s", false, "TLS encryption")

	flag.Parse()

	// set the true default value for port number based on appearance of -s in
	// command line arguments
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

// Makes connection to the server at the given CONNECT
// string (formatted as: hostname:port); if encrypt is true,
// make a TLS connection, else use non-TLS.
func makeConn(encrypt bool, CONNECT string) (net.Conn, *tls.Conn, error) {
	if encrypt {
		conn, err := tls.Dial("tcp", CONNECT, &tls.Config{})
		return nil, conn, err
	}
	connection, err := net.Dial("tcp", CONNECT)
	return connection, nil, err
}

// Read the response from the given server connection.
func readFromServer(connection net.Conn) (string, error) {

	// initialize a new Reader so ReadString() method can be used
	reader := bufio.NewReader(connection)

	// read response until a newline char is found (end of message
	// according to our protocol)
	line, readError := reader.ReadString('\n')

	// ensure no errors occurred during reading
	checkError(readError)

	// chop off the newline at the end of line (from the docs:
	// "returning a string containing the data up to and including the delimiter)"
	readLine := line[:len(line)-1]
	return readLine, nil
}

// Write the given data message to the given server connection.
func writeToServer(connection net.Conn, data string) {
	_, writeError := connection.Write([]byte(data))
	checkError(writeError)
}

// Counts and returns the number of occurrences of the ASCII symbol
// at index = 2 in the random string at index = 3 of the given string.
// e.g. countOccurrence(ex_string FIND v asdlkfjbvkjvks) --> 2
func countOccurrence(response string) int {
	stringArr := strings.Split(response, " ")
	return strings.Count(stringArr[3], stringArr[2])
}

// Checks if the given error exists and panics if it does. Else, do
// nothing since no error occurred.
func checkError(err error) {
	if err != nil {
		panic("failed in communication with server; reason: " + err.Error())
	}
}

// Verifies that the given string meets the requirements for a valid response
// from the server, i.e. the number of arguments is correct, the string starts
// with "ex_string", and the command-specific criteria are met.
func verifyResponse(response string) {
	stringArr := strings.Split(response, " ")
	if len(stringArr) < 2 || stringArr[0] != "ex_string" || invalidCommand(stringArr[1], stringArr) {
		panic("Response does not conform to the protocol! " + response)
	}
}

// Verifies that the specific rules for commands supported by the TCP server
// are met, i.e. the HELLO, FIND, COUNT, and BYE commands have the correct
// amount of args accompanying them when such a message is received from the server.
func invalidCommand(command string, array []string) bool {
	switch command {
	case "BYE":
		return len(array) != 3
	case "FIND":
		return len(array) != 4
	case "HELLO":
		return len(array) != 3
	case "COUNT":
		return len(array) != 3
	default:
		return true
	}
}
