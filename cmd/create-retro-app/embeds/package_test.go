package embeds

import (
	"bytes"
	"testing"

	"github.com/zaydek/retro/pkg/expect"
)

func TestJavaScriptTemplate(t *testing.T) {
	want := `{
	"name": "app-name",
	"scripts": {
		"dev": "retro dev",
		"build": "retro build",
		"serve": "retro serve"
	},
	"dependencies": {
		"@zaydek/retro": "^1.33.7",
		"react": "^17.0.1",
		"react-dom": "^17.0.1"
	},
	"devDependencies": {
		"esbuild": "^0.8.46"
	}
}
`

	dot := PackageDot{
		APP_NAME:      "app-name",
		RETRO_VERSION: "1.33.7",
	}
	var buf bytes.Buffer
	if err := JavaScriptPackageTemplate.Execute(&buf, dot); err != nil {
		t.Fatal(err)
	}
	expect.DeepEqual(t, buf.String(), want)
}

func TestTypeScriptTemplate(t *testing.T) {
	want := `{
	"name": "app-name",
	"scripts": {
		"dev": "retro dev",
		"build": "retro build",
		"serve": "retro serve"
	},
	"dependencies": {
		"@zaydek/retro": "^1.33.7",
		"react": "^17.0.1",
		"react-dom": "^17.0.1"
	},
	"devDependencies": {
		"@types/react": "^17.0.0",
		"@types/react-dom": "^17.0.0",
		"esbuild": "^0.8.46"
	}
}
`

	dot := PackageDot{
		APP_NAME:      "app-name",
		RETRO_VERSION: "1.33.7",
	}
	var buf bytes.Buffer
	if err := TypeScriptPackageTemplate.Execute(&buf, dot); err != nil {
		t.Fatal(err)
	}
	expect.DeepEqual(t, buf.String(), want)
}
