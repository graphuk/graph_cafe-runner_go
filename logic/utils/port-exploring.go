package utils

import (
	"net"
	"strconv"
	"time"
)

func GetFirstFreeLocalPorts(portFrom, portTo int) (int, int) {
	res := 0
	found := false
	for i := portFrom; i <= portTo; i++ {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("", strconv.Itoa(i)), time.Second)
		if err == nil {
			conn.Close()
		} else {
			conn, err := net.DialTimeout("tcp", net.JoinHostPort("", strconv.Itoa(i+1)), time.Second)
			if err == nil {
				conn.Close()
			} else {
				res = i
				found = true
				break
			}
		}
	}
	if !found {
		panic(`Free ports not found in the range: ` + strconv.Itoa(portFrom) + ` - ` + strconv.Itoa(portTo))
	}
	return res, res + 1
}
