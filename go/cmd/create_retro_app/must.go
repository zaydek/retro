package create_retro_app

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}
