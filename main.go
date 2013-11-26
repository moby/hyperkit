package gist7651991

import (
	"bufio"
	"io"

	"sync"

	. "gist.github.com/5892738.git"
)

func ProcessLinesFromReader(r io.Reader, processFunc func(string)) {
	br := bufio.NewReader(r)
	for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
		processFunc(MustTrimLastNewline(line))
	}
}

func GoReduceLinesFromReader(r io.Reader, numWorkers int, reduceFunc func(string) (string, bool)) <-chan string {
	outChan := make(chan string)

	go func() {
		inChan := make(chan string)
		var wg sync.WaitGroup

		for worker := 0; worker < numWorkers; worker++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					switch in, ok := <-inChan; {
					case ok:
						if out, ok := reduceFunc(in); ok {
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
