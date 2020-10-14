package cmd

import "testing"

func TestGetPathWithFullPath(t *testing.T) {
	testPath := "/corpa/dev/foo_app/"
	path := getPath([]string{testPath})
	if path != testPath {
		t.Errorf("want: %s | got %s", testPath, path)
	}
}

func TestGetPathEnvVarName(t *testing.T) {
	envVar := "TEST_VAR"
	path := getPath([]string{envVar})
	if path != envVar {
		t.Errorf("want: %s | got %s", envVar, path)
	}
}
