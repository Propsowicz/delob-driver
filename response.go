package delobdriver

import (
	"fmt"
	"strconv"
	"strings"
)

type response struct {
	version string
	status  status
	msg     string
}

type status int8

const (
	fail                     status = 0
	success                  status = 1
	user_verified            status = 7
	user_not_verified        status = 8
	authentication_challenge status = 9
)

func parseStatus(s string) status {
	num, err := strconv.Atoi(s)
	if err != nil {
		return fail
	}

	switch num {
	case 0:
		return fail
	case 1:
		return success
	case 7:
		return user_verified
	case 8:
		return user_not_verified
	case 9:
		return authentication_challenge

	default:
		return fail
	}
}

func newResponse(s string) (response, error) {
	p := response{}
	if len(s) < 3 {
		return p, fmt.Errorf("wrong response format.")
	}
	p.version = s[0:2]

	p.status = parseStatus(s[2:3])
	p.msg = strings.TrimSpace(s)[3:]

	return p, nil
}
