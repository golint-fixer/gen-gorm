package util

// HandleErr handles errors in a consistent way
func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}
