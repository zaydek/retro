package cli

type DevCommand struct {
	Port      int
	Sourcemap bool
}

type BuildCommand struct {
	Sourcemap bool
}

type ServeCommand struct {
	Port int
}
