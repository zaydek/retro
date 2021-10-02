package embeds

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}
