package delobdriver

import "fmt"

type request struct {
	user string
	msg  string
}

func newRequest(user, msg string) request {
	return request{
		user: user,
		msg:  msg,
	}
}

func (r *request) toString() string {
	return fmt.Sprintf("%s\r\n%s\n", r.user, r.msg)
}
