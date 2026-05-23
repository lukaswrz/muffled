package gonfig

import (
	"errors"
	"fmt"
	"os"
)

type UnmarshalFunc[T any] func([]byte, T) error

type FinalizeFunc[T any] func(T) error

func ReadConfig[T any](path string, searchPaths []string, c *T, unmarshal UnmarshalFunc[*T], finalize FinalizeFunc[*T]) (string, error) {
	var err error

	path, err = FindConfig(path, searchPaths)
	if err != nil {
		return "", err
	}

	return path, ReadFoundConfig(path, c, unmarshal, finalize)
}

func ReadFoundConfig[T any](path string, c *T, unmarshal UnmarshalFunc[*T], finalize FinalizeFunc[*T]) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read configuration file %q: %w", path, err)
	}

	err = unmarshal(content, c)
	if err != nil {
		return fmt.Errorf("unable to unmarshal configuration file %q: %w", path, err)
	}

	return finalize(c)
}

func FindConfig(path string, paths []string) (string, error) {
	var err error

	if path == "" {
		for _, p := range paths {
			_, err = os.Stat(p)
			if err != nil {
				continue
			}

			path = p
		}

		if path == "" {
			return "", errors.New("could not locate configuration file")
		}
	} else {
		_, err = os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("could not stat configuration file %q: %w", path, err)
		}
	}

	return path, nil
}
