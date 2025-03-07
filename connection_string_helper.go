package delobdriver

import (
	"fmt"
	"strings"
)

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
