package pipe_util

import "gopkg.in/pipe.v2"

// OutputDir is identical to pipe.Output, except it sets the starting dir.
func OutputDir(p pipe.Pipe, dir string) ([]byte, error) {
	outb := &pipe.OutputBuffer{}
	s := pipe.NewState(outb, nil)
	s.Dir = dir
	err := p(s)
	if err == nil {
		err = s.RunTasks()
	}
	return outb.Bytes(), err
}
