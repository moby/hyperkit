package gist5286084

// CheckError panics on error.
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
