package delobdriver

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type AuthenticationManager struct {

	// handshake client <-> data base SCRAM (password)
	// tcp client <-> data base
	//
	// 1. Handshake
	// 2. userData persistance
	// 3. session storage

}

func NewAuthenticationManager() AuthenticationManager {
	return AuthenticationManager{}
}

// 0-1 failure/success
// server responses:
// 9 - not auth -> challenge me!
// 8 - here is challenge data
// 6/7 forbidden/access
// client --> server
// server -> data (salt, random number) -> client
// client -> generate hash -> server
// server generate hash -> compare hashes
// ok -> 7 & save session
// not OK -> forbid - 6

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

func generateNonce() int {
	return rand.Intn(256)
}
