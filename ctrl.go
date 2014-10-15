package bur

import "sync"

func ctrlServer(config *Config, wg *sync.WaitGroup) {
	defer wg.Done()
}
