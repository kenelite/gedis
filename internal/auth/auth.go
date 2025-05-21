package auth

import (
	"github.com/kenelite/gedis/internal/config"
	"sync"
)

var (
	users = make(map[string]string)
	mu    sync.RWMutex
)

func LoadUsersFromConfig(path string) error {
	cfg, err := config.Load(path)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	userSection := cfg.GetSection("users")
	if userSection == nil {
		return nil // no users section found
	}

	for k, v := range userSection {
		users[k] = v
	}

	return nil
}

func CheckUser(username, password string) bool {
	mu.RLock()
	defer mu.RUnlock()

	if pass, ok := users[username]; ok {
		return pass == password
	}

	if defPass, ok := users["default"]; ok {
		return defPass == password
	}

	return false
}
