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
	const uniqueDelimiter string = "\x1E\x1F"

	return fmt.Sprintf("%s%s%s\n", r.user, uniqueDelimiter, r.msg)
}
