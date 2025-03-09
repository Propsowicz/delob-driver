package delobdriver

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const clientKeySalt string = "Client Key"
const serverKeySalt string = "Server Key"

type AuthenticationManager struct {
}

func NewAuthenticationManager() AuthenticationManager {
	return AuthenticationManager{}
}

func (a *AuthenticationManager) addServerFirstAuthString(auth, salt string, s_nonce, iterations int) string {
	auth = a.addToAuth(auth, fmt.Sprintf("s_nonce=%d,", s_nonce))
	auth = a.addToAuth(auth, fmt.Sprintf("salt=%s,", salt))
	auth = a.addToAuth(auth, fmt.Sprintf("iterations=%d", iterations))
	return auth
}

func (a *AuthenticationManager) addClientFirstAuthString(user string, c_nonce int) string {
	var auth string
	auth = a.addToAuth(auth, fmt.Sprintf("user=%s,", user))
	auth = a.addToAuth(auth, fmt.Sprintf("c_nonce=%d,", c_nonce))
	return auth
}

func (a *AuthenticationManager) addToAuth(auth string, s interface{}) string {
	var toAdd string
	switch v := s.(type) {
	case string:
		toAdd = v
	case int, int8:
		toAdd = fmt.Sprintf("%d", v)
	}

	return auth + toAdd
}

func (a *AuthenticationManager) parseServerFirst(s string) (int, string, int, error) {
	parts := strings.Split(s, ",")

	const snoncePrefix string = "s_nonce="
	const saltPrefix string = "salt="
	const iterationsPrefix string = "iterations="
	if parts[2][0:len(snoncePrefix)] == snoncePrefix && parts[3][0:len(saltPrefix)] == saltPrefix && parts[4][0:len(iterationsPrefix)] == iterationsPrefix {
		s_nonce, err := strconv.Atoi(parts[2][len(snoncePrefix):])
		if err != nil {
			return 0, "", 0, err
		}
		iterations, err := strconv.Atoi(parts[4][len(iterationsPrefix):])
		if err != nil {
			return 0, "", 0, err
		}
		return s_nonce, parts[3][len(saltPrefix):], iterations, nil
	}

	return 0, "", 0, fmt.Errorf("cannot parse client first message")
}

func (a *AuthenticationManager) calculateProof(password, salt, auth string, iterations int) string {
	saltByte, err := hex.DecodeString(salt)
	if err != nil {
		panic(err.Error())
	}
	hashedPassword := calculateHashedPassword(password, saltByte, iterations)
	clientKey := computeHmacHash(hashedPassword, []byte(clientKeySalt))
	storedKey := computeSha256Hash(clientKey)

	clientSignature := computeHmacHash(storedKey, []byte(auth))

	return hex.EncodeToString(xorBytes(clientKey, clientSignature))
}
