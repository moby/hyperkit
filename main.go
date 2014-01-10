package gist7651991

import (
	"bufio"
	"io"
	"sync"
)

func ProcessLinesFromReader(r io.Reader, processFunc func(string)) {
	br := bufio.NewReader(r)
	for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
		processFunc(line[:len(line)-1]) // Trim last newline
	}
}

func GoReduceLinesFromReader(r io.Reader, numWorkers int, reduceFunc func(string) interface{}) <-chan interface{} {
	outChan := make(chan interface{})

	go func() {
		inChan := make(chan string)
		var wg sync.WaitGroup

		// TODO: See if I can create goroutines alongside with the work, up to a max number, rather than all in advance
		// Create numWorkers goroutines
		for worker := 0; worker < numWorkers; worker++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					switch in, ok := <-inChan; {
					case ok:
						if out := reduceFunc(in); out != nil {
							outChan <- out
						}
					case !ok:
						return
					}
				}
			}()
		}

		ProcessLinesFromReader(r, func(in string) { inChan <- in })
		close(inChan)
		wg.Wait()
		close(outChan)
	}()

	return outChan
}
