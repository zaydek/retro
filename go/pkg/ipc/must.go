package ipc

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}