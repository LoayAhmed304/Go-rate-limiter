package structures

import "time"

// ///////////////// Route        IP     ClientLogs
var ClientsLogs *map[string]map[string]*[]time.Time

// InitClientsLogs initializes the ClientsLogs map with
// empty slices for each route defined in the configurations.
func InitClientsLogs(routes []string) {
	s := make(map[string]map[string]*[]time.Time)

	for _, route := range routes {
		s[route] = make(map[string]*[]time.Time)
	}

	ClientsLogs = &s
}
