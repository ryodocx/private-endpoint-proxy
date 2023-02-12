package util

import (
	"io/fs"
	"path"
)

func Dirwalk(f fs.FS, dir string) (paths []string, err error) {
	entries, err := fs.ReadDir(f, dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			nestedPaths, err := Dirwalk(f, path.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			paths = append(paths, nestedPaths...)
			continue
		}
		paths = append(paths, path.Join(dir, entry.Name()))
	}
	return paths, nil
}
