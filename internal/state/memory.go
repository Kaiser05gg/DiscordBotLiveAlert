package state

import "sync"

// user_login -> isLive
var (
	mu     sync.Mutex
	status = map[string]bool{}
)

// WasLive returns previous state
func WasLive(user string) bool {
	mu.Lock()
	defer mu.Unlock()
	return status[user]
}

// SetLive updates current state
func SetLive(user string, isLive bool) {
	mu.Lock()
	defer mu.Unlock()
	status[user] = isLive
}
