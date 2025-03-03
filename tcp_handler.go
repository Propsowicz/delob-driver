package delobdriver

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type tcpHandler struct {
	connectionString connectionString
	protocolVersion  string
	conn             net.Conn
	reader           *bufio.Reader
	writer           *bufio.Writer
}

func newTcpHandler(rawConnectionString string) (*tcpHandler, error) {
	connectionString, errConStr := parseConnectionString(rawConnectionString)
	if errConStr != nil {
		return nil, errConStr
	}

	conn, err := net.Dial("tcp", connectionString.adress)
	if err != nil {
		return nil, err
	}

	return &tcpHandler{
		connectionString: connectionString,
		protocolVersion:  "00", // TODO
		conn:             conn,
		reader:           bufio.NewReader(conn),
		writer:           bufio.NewWriter(conn),
	}, nil
}

func (h *tcpHandler) sendMessage(message string) (string, error) {
	_, err := h.writer.WriteString(message + " \n")
	if err != nil {
		return "", err
	}
	if err := h.writer.Flush(); err != nil {
		return "", err
	}

	response, err := h.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	rawResponse := strings.TrimSpace(response)
	protocolVersion := rawResponse[0:2]
	executionState := rawResponse[2:3]

	if protocolVersion != h.protocolVersion {
		return "", fmt.Errorf("wrong protocol version")
	}
	if executionState == "0" {
		return "", fmt.Errorf("expression execution was not successful: %s", strings.TrimSpace(response)[3:])
	}

	return strings.TrimSpace(response)[3:], nil
}

type connectionString struct {
	server   string
	port     string
	adress   string
	username string
	password string
}

func parseConnectionString(rawConnectionString string) (connectionString, error) {
	const defaultPort string = "5678"
	const serverKey string = "server"
	const portKey string = "port"
	const uidKey string = "uid"
	const pwdKey string = "pwd"

	tokens := strings.Split(rawConnectionString, ";")
	connectionString := connectionString{}

	for i := range tokens {
		tokenKeyValue := strings.Split(tokens[i], "=")
		switch strings.ToLower(tokenKeyValue[0]) {
		case serverKey:
			connectionString.server = tokenKeyValue[1]
		case portKey:
			connectionString.port = tokenKeyValue[1]
		case uidKey:
			connectionString.username = tokenKeyValue[1]
		case pwdKey:
			connectionString.password = tokenKeyValue[1]
		}
	}
	if connectionString.port == "" {
		connectionString.port = defaultPort
	}

	if err := validateConnectionStringElement(connectionString.server, serverKey); err != nil {
		return connectionString, err
	}
	if err := validateConnectionStringElement(connectionString.username, uidKey); err != nil {
		return connectionString, err
	}
	if err := validateConnectionStringElement(connectionString.password, pwdKey); err != nil {
		return connectionString, err
	}

	connectionString.adress = connectionString.server + ":" + connectionString.port

	return connectionString, nil
}

func validateConnectionStringElement(element, key string) error {
	if element == "" {
		return fmt.Errorf("cannot find %s element in connection string", key)
	}
	return nil
}
