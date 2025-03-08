package delobdriver

import (
	"crypto/hmac"
	"crypto/sha256"
	"math/rand"

	"golang.org/x/crypto/pbkdf2"
)

func computeHmacHash(arg_1, arg_2 []byte) []byte {
	mac := hmac.New(sha256.New, arg_1)
	mac.Write([]byte(arg_2))
	return mac.Sum(nil)
}

func computeSha256Hash(arg_1 []byte) []byte {
	hash := sha256.Sum256(arg_1)
	return hash[:]
}

func xorBytes(k, j []byte) []byte {
	if len(k) != len(j) {
		panic("byte slices must be of equal length")
	}
	result := make([]byte, len(k))
	for i := range k {
		result[i] = k[i] ^ j[i]
	}
	return result
}

func calculateHashedPassword(password, salt string, iterations int) []byte {
	const keyLength int = 32

	hashedPassword := pbkdf2.Key([]byte(password), []byte(salt), iterations, keyLength, sha256.New)
	return hashedPassword
}

func generateNonce() int {
	return rand.Intn(256)
}
