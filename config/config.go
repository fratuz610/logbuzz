package config

import (
	hosieweb "github.com/hoisie/web"
	"logbuzz/data"
	"sync"
)

var currentConfig data.Config = data.Config{HttpPort: 1212, MemoryLimit: 20 * 1024 * 1024}
var configMutex sync.RWMutex

func GetConfig() data.Config {
	configMutex.RLock()
	defer configMutex.RUnlock()

	return currentConfig
}

func UpdateHttpPort(httpPort int) {
	configMutex.Lock()
	defer configMutex.Unlock()
	if currentConfig.HttpPort == httpPort {
		return
	}
	currentConfig.HttpPort = httpPort

	// we stop the server causing it to restart
	hosieweb.Close()

}

func UpdateMemoryLimit(memoryLimit uint64) {
	configMutex.Lock()
	defer configMutex.Unlock()

	if currentConfig.MemoryLimit == memoryLimit {
		return
	}
	currentConfig.MemoryLimit = memoryLimit

	// nothing to do
}
