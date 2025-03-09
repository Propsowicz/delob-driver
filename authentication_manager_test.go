package delobdriver

import (
	"fmt"
	"testing"
)

func Test_IfProofIsCalculatedCorrectly(t *testing.T) {
	authManager := NewAuthenticationManager()
	expectedHardCodedVaLUE := "fca9e49b0c0baf2cd191bd6d4c8575042e5c07b038579acf896f6e02e5867d11"

	result := authManager.calculateProof("myUsername", "f2c147cd7cdfbfad39ce965dc7288467", "someRandomAuthString", 5)

	if result != expectedHardCodedVaLUE {
		t.Errorf("proof is not correctly calcualted.")
	}
}

func Test_IfParseServerFirstWorks(t *testing.T) {
	authManager := NewAuthenticationManager()
	s_nonceMock := 15
	saltMock := "saltRnd"
	iterationsMock := 12

	s_nonce, salt, iterations, errServerFirstParse := authManager.parseServerFirst(fmt.Sprintf("user=user,c_nonce=1,s_nonce=%d,salt=%s,iterations=%d",
		s_nonceMock, saltMock, iterationsMock))

	if errServerFirstParse != nil {
		t.Errorf("should not return error.")
	}
	if s_nonce != s_nonceMock {
		t.Errorf("wrong s_nonce")
	}
	if salt != saltMock {
		t.Errorf("wrong salt")
	}
	if iterations != iterationsMock {
		t.Errorf("wrong iteration value")
	}
}

func Test_IfIterationWithWrongFormatReturnError(t *testing.T) {
	authManager := NewAuthenticationManager()
	s_nonceMock := 15
	saltMock := "saltRnd"

	_, _, _, errServerFirstParse := authManager.parseServerFirst(fmt.Sprintf("user=user,c_nonce=1,s_nonce=%d,salt=%s,iterations=ARGH!@#",
		s_nonceMock, saltMock))

	if errServerFirstParse == nil {
		t.Errorf("should return error.")
	}
}

func Test_IfFillingAuthStringIsWorkginCorrectly(t *testing.T) {
	authManager := NewAuthenticationManager()
	user := "username"
	c_nonce := 12
	s_nonce := 2
	salt := "qweasd"
	iterations := 5
	expected := fmt.Sprintf(fmt.Sprintf("user=%s,c_nonce=%d,s_nonce=%d,salt=%s,iterations=%d",
		user, c_nonce, s_nonce, salt, iterations))

	auth := authManager.addClientFirstAuthString(user, c_nonce)
	auth = authManager.addServerFirstAuthString(auth, salt, s_nonce, iterations)

	if auth != expected {
		t.Errorf("auth is not filled with data properly.")
	}
}
