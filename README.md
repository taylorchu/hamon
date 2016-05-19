# hamon

Monitor haproxy server state changes.

## Example

```go
	var (
		counter   int
		now       time.Time
		printTime bool
	)
	group := hamon.Group{
		Socket: "/var/run/haproxy.sock",
		Server: []hamon.Server{
			{Name: "kombucha2"},
			{Name: "kombucha3"},
			{Name: "kombucha4"},
			{Name: "kombucha5"},
			{Name: "kombucha6"},
		},
		PollInterval: 10 * time.Millisecond,
		StateFunc: func(e *hamon.Event) {
			if e.Time.After(now) {
				now = e.Time
				counter++
				printTime = false
			}

			if !printTime {
				printTime = true
				fmt.Fprintf(os.Stderr, "-- %d --\n", counter)
			}

			fmt.Fprintf(os.Stderr, "%s goes from %s to %s\n", e.ServerName, e.OldState, e.NewState)
		},
	}
	group.Monitor()
```
