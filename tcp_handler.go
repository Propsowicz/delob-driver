package delobdriver

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type tcpHandler struct {
	authManager      AuthenticationManager
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
		authManager:      NewAuthenticationManager(),
		connectionString: connectionString,
		protocolVersion:  "00", // TODO
		conn:             conn,
		reader:           bufio.NewReader(conn),
		writer:           bufio.NewWriter(conn),
	}, nil
}

func (h *tcpHandler) sendRequest(message string) (string, error) {
	response, errRespParse := h.getResponse(message)
	if errRespParse != nil {
		return "", errRespParse
	}

	switch response.status {
	case fail:
		return response.msg, fmt.Errorf(response.msg)
	case success:
		return response.msg, nil
	case authChallenge:

		s_nonce, salt, iterations, errServerFirstParse, auth := h.getClientFirstMessage()
		if errServerFirstParse != nil {
			return "", errServerFirstParse
		}

		fmt.Println(s_nonce)
		fmt.Println(salt)
		fmt.Println(iterations)
		fmt.Println(auth)

		// verify proof -> re-send initial msg
	}

	return response.msg, nil
}

func (h *tcpHandler) getClientFirstMessage() (int, string, int, error, string) {
	auth := h.authManager.addClientFirstAuthString(h.connectionString.username, generateNonce())

	response, errRespParse := h.getResponse(auth)
	if errRespParse != nil {
		return 0, "", 0, errRespParse, ""
	}

	s_nonce, salt, iterations, errServerFirstParse := h.authManager.parseServerFirst(response.msg)
	if errServerFirstParse != nil {
		return 0, "", 0, errServerFirstParse, ""
	}

	return s_nonce, salt, iterations, errServerFirstParse, h.authManager.addServerFirstAuthString(auth, salt, s_nonce, iterations)
}

func (h *tcpHandler) getResponse(requestMessage string) (response, error) {
	request := newRequest(h.connectionString.username, requestMessage)
	_, err := h.writer.WriteString(request.toString())
	if err != nil {
		return response{}, err
	}
	if err := h.writer.Flush(); err != nil {
		return response{}, err
	}

	streamResponse, err := h.reader.ReadString('\n')
	if err != nil {
		return response{}, err
	}

	rawResponse := strings.TrimSpace(streamResponse)

	return newResponse(rawResponse)
}
