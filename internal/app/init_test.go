package app

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func readEspFile(path string) espFile {
	espFile := espFile{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, espFile)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return espFile
}

func checkEspFile(c1 espFile, c2 espFile) bool {
	if c1.Backend != c2.Backend {
		return false
	}

	if c1.OrgName != c2.OrgName {
		return false
	}

	if c1.OrgPrefix != c2.OrgPrefix {
		return false
	}

	if c1.AppName != c2.AppName {
		return false
	}

	for i, e := range c1.Envs {
		if e != c2.Envs[i] {
			return false
		}
	}

	return true
}

func TestWriteConfig(t *testing.T) {
	c := Config{}
	tc := Config{}
	tmpdir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	c.Path = filepath.Join(tmpdir, ".espFile.yaml")
	ci := configInput{
		Backend: "ssm",
		OrgName: "foo",
		OrgPrefix: "FOO",
		AppName: "test-service",
		Envs: "dev,test,prod",
	}
	c.UpdateWithInput(ci)
	if err = c.WriteConfig(); err != nil {
		log.Fatal(err)
	}
	actualEsp := readEspFile(c.Path)
	tc.UpdateWithInput(ci)
	testEsp := tc.createEspFile()

	if checkEspFile(actualEsp, testEsp) {
		t.Errorf("The written config didn't match the test input")
	}
}