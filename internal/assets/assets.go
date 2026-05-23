package assets

import (
	"embed"
	"io/fs"
)

//go:embed assets
var embedded embed.FS

var FS fs.FS

func init() {
	var err error
	FS, err = fs.Sub(embedded, "assets")
	if err != nil {
		panic(err)
	}
}
