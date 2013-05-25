package gist5286084

import (
)

// Panics on error
func CheckError(err error) {
	if nil != err {
		panic(err)
	}
}