package cmd

import "testing"

func TestGetPathWithFullPath(t testing.T) {
	testPath := "/corpa/dev/foo_app/"
	paths := []string{testPath}
	path := getPath(paths)
	if path != testPath {
		t.Errorf("want: %s | got %s", testPath, path)
	}
}
