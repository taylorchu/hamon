package hamon

import (
	"strings"
	"time"
)

type ServerState uint8

func (state ServerState) String() string {
	switch state {
	case ServerUp:
		return "UP"
	case ServerDown:
		return "DOWN"
	case ServerMaint:
		return "MAINT"
	default:
		return "UNKNOWN"
	}
}

const (
	ServerUnknown ServerState = iota
	ServerUp
	ServerDown
	ServerMaint
)

type Server struct {
	Name  string
	State ServerState
}

type Group struct {
	Socket       string
	PollInterval time.Duration
	Server       []Server
	StateFunc    func(*Event)
}

type Event struct {
	Time               time.Time
	ServerName         string
	OldState, NewState ServerState
}

// Monitor checks server state changes every PollInterval.
// It calls StateFunc to announce changes.
func (g *Group) Monitor() error {
	for now := range time.Tick(g.PollInterval) {
		table, err := showStat(g.Socket)
		if err != nil {
			return err
		}
		for i := range g.Server {
			s := &g.Server[i]
			for _, row := range table {
				if s.Name != row[1] {
					continue
				}

				newState := ServerUnknown
				switch {
				case strings.HasPrefix(row[17], "UP"):
					newState = ServerUp
				case strings.HasPrefix(row[17], "DOWN"):
					newState = ServerDown
				case strings.HasPrefix(row[17], "MAINT"):
					newState = ServerMaint
				}

				if s.State != newState {
					g.StateFunc(&Event{
						ServerName: s.Name,
						Time:       now,
						OldState:   s.State,
						NewState:   newState,
					})
					s.State = newState
				}
			}
		}
	}
	return nil
}
