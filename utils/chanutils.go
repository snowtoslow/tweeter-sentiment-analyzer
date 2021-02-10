package utils

import (
	"sync"
)

//cesto spizdil from https://blog.golang.org/pipelines because I'm pull stack developer
func MergeWait(cs ...chan string) chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan string) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
