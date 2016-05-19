package hamon

import (
	"encoding/csv"
	"net"
)

func showStat(socket string) ([][]string, error) {
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte("show stat\n"))
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(conn)
	r.Comment = '#'
	return r.ReadAll()
}
