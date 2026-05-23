package widget

import (
	"fmt"
	"io"
	"os"

	"hack.moontide.ink/pingfisher/muffled/internal/assets"
)

func Get(path string) ([]byte, error) {
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("unable to read widget file %q: %w", path, err)
		}

		return data, nil
	}

	const dw = "widget.html"

	file, err := assets.FS.Open(dw)
	if err != nil {
		return nil, fmt.Errorf("open embedded widget %q: %w", dw, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read embedded widget %q: %w", dw, err)
	}

	return data, nil
}
