package bur

import "sync"

func ctrlServer(wg *sync.WaitGroup) {
	defer wg.Done()
}
