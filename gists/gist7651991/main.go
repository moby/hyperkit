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

// GoReduceLinesFromReader spawns numWorkers goroutines and reduces each line from reader with reduceFunc.
//
//	source := strings.NewReader(`1
//	2
//	three
//	four
//	etc.
//	`)
//
//	reduceFunc := func(in string) interface{} {
//		time.Sleep(2 * time.Second)
//		return "Hello: " + in
//	}
//
//	outChan := gist7651991.GoReduceLinesFromReader(source, 4, reduceFunc)
//
//	for out := range outChan {
//		fmt.Println(out)
//	}
//
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
					if in, ok := <-inChan; ok {
						if out := reduceFunc(in); out != nil {
							outChan <- out
						}
					} else {
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

func GoReduceLinesFromSlice(inSlice []string, numWorkers int, reduceFunc func(string) interface{}) <-chan interface{} {
	inChan := make(chan interface{})
	go func() { // This needs to happen in the background because sending input will be blocked on reading output.
		for _, in := range inSlice {
			inChan <- in
		}
		close(inChan)
	}()
	reduceFuncWrapper := func(in interface{}) interface{} { return reduceFunc(in.(string)) }
	outChan := GoReduce(inChan, numWorkers, reduceFuncWrapper)

	return outChan
}

// Caller is expected to close inChan after sending all input to it. Sending input should be done in a background goroutine,
// because sending input will be blocked on reading output.
func GoReduce(inChan <-chan interface{}, numWorkers int, reduceFunc func(interface{}) interface{}) <-chan interface{} {
	outChan := make(chan interface{})

	go func() {
		var wg sync.WaitGroup

		// TODO: See if I can create goroutines alongside with the work, up to a max number, rather than all in advance
		// Create numWorkers goroutines
		for worker := 0; worker < numWorkers; worker++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					in, ok := <-inChan
					if !ok {
						return
					}

					out := reduceFunc(in)
					if out != nil {
						outChan <- out
					}
				}
			}()
		}

		wg.Wait()
		close(outChan)
	}()

	return outChan
}

// Caller is expected to close inChan after sending all input to it. Sending input should be done in a background goroutine,
// because sending input will be blocked on reading output.
// Order is guaranteed to be preserved.
//
// TODO: Implement this.
func GoReducePreservingOrder(inChan <-chan interface{}, numWorkers int, reduceFunc func(interface{}) interface{}) <-chan interface{} {
	panic("not implemented")
}
