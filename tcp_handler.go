package delobdriver

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const requestLimit int8 = 5

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

func (h *tcpHandler) sendRequest(message string, requestCounter int8) (string, error) {
	if requestCounter > requestLimit {
		return "", fmt.Errorf("connection limit has been exceeded.")
	}

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

		auth := h.prepareClientFirstAuthString()

		s_nonce, salt, iterations, errServerFirstRequest := h.getServerFirstMessage(auth)
		if errServerFirstRequest != nil {
			return "", errServerFirstRequest
		}

		auth = h.prepareClientFinalAuthString(auth, salt, s_nonce, iterations)

		proof := h.prepareProof(h.connectionString.password, salt, auth, iterations)

		errVerifierRequest := h.getVerifier(proof)
		if errVerifierRequest != nil {
			return "", errVerifierRequest
		}

		return h.sendRequest(message, requestCounter+1)
	}

	return response.msg, nil
}

func (h *tcpHandler) getVerifier(proof string) error {
	response, errRespParse := h.getResponse(proof)
	if errRespParse != nil {
		return errRespParse
	}

	if response.status == proofVerified {
		return nil
	}

	return fmt.Errorf("cannot authenticate user")
}

func (h *tcpHandler) prepareProof(password, salt, auth string, iterations int) string {
	return h.authManager.calculateProof(password, salt, auth, iterations)
}

func (h *tcpHandler) prepareClientFirstAuthString() string {
	return h.authManager.addClientFirstAuthString(h.connectionString.username, generateNonce())
}

func (h *tcpHandler) prepareClientFinalAuthString(auth, salt string, s_nonce, iterations int) string {
	return h.authManager.addServerFirstAuthString(auth, salt, s_nonce, iterations)
}

func (h *tcpHandler) getServerFirstMessage(auth string) (int, string, int, error) {
	response, errRespParse := h.getResponse(auth)
	if errRespParse != nil {
		return 0, "", 0, errRespParse
	}

	if response.status == fail {
		return 0, "", 0, fmt.Errorf(response.msg)
	}

	s_nonce, salt, iterations, errServerFirstParse := h.authManager.parseServerFirst(response.msg)
	if errServerFirstParse != nil {
		return 0, "", 0, errServerFirstParse
	}

	return s_nonce, salt, iterations, errServerFirstParse
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
