package cli

// Describes the dev command
type DevCommand struct {
	Port      int
	Sourcemap bool
}

// Describes the build command
type BuildCommand struct {
	Sourcemap bool
}

// Describes the serve command
type ServeCommand struct {
	Port int
}
