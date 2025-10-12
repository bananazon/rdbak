package globals

import (
	"fmt"
	"os/user"
	"sync"
)

var (
	hexKey  string = "2a89f4811dae393b45e9d902388783b39a15ddd63a2cbf19b90f6a3e3dd9a06b"
	homeDir string
	mu      sync.RWMutex
)

func GetHexKey() (x string) {
	mu.Lock()
	x = hexKey
	mu.Unlock()
	return x
}

func SetHomeDirectory() (err error) {
	mu.Lock()
	userObj, err := user.Current()
	if err != nil {
		mu.Unlock()
		return fmt.Errorf("failed to determine the path of your home directory: %s", err)
	}
	homeDir = userObj.HomeDir
	mu.Unlock()
	return nil
}

func GetHomeDirectory() (x string) {
	mu.Lock()
	x = homeDir
	mu.Unlock()
	return x
}
