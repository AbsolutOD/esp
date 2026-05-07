package app

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func readEspFile(path string) espFile {
	espFile := espFile{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &espFile)
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

	if !checkEspFile(actualEsp, testEsp) {
		t.Errorf("The written config didn't match the test input")
	}
}

func TestUpdateWithInput(t *testing.T) {
	tests := []struct {
		name     string
		envs     string
		wantEnvs []string
	}{
		{
			name:     "no whitespace",
			envs:     "dev,test,prod",
			wantEnvs: []string{"dev", "test", "prod"},
		},
		{
			name:     "single trailing space after each comma",
			envs:     "dev, test, prod",
			wantEnvs: []string{"dev", "test", "prod"},
		},
		{
			name:     "multiple trailing spaces after each comma",
			envs:     "dev,   test,  prod",
			wantEnvs: []string{"dev", "test", "prod"},
		},
		{
			name:     "single env",
			envs:     "dev",
			wantEnvs: []string{"dev"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := Config{}
			c.UpdateWithInput(configInput{
				Backend:   "ssm",
				OrgName:   "acme",
				OrgPrefix: "ACME",
				AppName:   "billing",
				Envs:      tc.envs,
			})
			if c.Backend != "ssm" || c.OrgName != "acme" || c.OrgPrefix != "ACME" || c.AppName != "billing" {
				t.Errorf("scalar fields not copied: got %+v", c)
			}
			if !reflect.DeepEqual(c.Envs, tc.wantEnvs) {
				t.Errorf("Envs = %#v, want %#v", c.Envs, tc.wantEnvs)
			}
		})
	}
}

func TestCreateEspFile(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want espFile
	}{
		{
			name: "fully populated",
			cfg: Config{
				Backend:   "ssm",
				OrgName:   "acme",
				OrgPrefix: "ACME",
				AppName:   "billing",
				Envs:      []string{"dev", "test", "prod"},
			},
			want: espFile{
				Backend:   "ssm",
				OrgName:   "acme",
				OrgPrefix: "ACME",
				AppName:   "billing",
				Envs:      []string{"dev", "test", "prod"},
			},
		},
		{
			name: "zero value Config produces zero value espFile",
			cfg:  Config{},
			want: espFile{},
		},
		{
			name: "non-espFile fields on Config are not copied",
			cfg: Config{
				Backend:      "ssm",
				OrgName:      "acme",
				AppName:      "billing",
				Envs:         []string{"dev"},
				IsEspProject: true,
				Env:          "dev",
				Path:         "/tmp/somewhere",
				Filename:     ".espFile",
			},
			want: espFile{
				Backend: "ssm",
				OrgName: "acme",
				AppName: "billing",
				Envs:    []string{"dev"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.cfg.createEspFile()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("createEspFile() = %#v, want %#v", got, tc.want)
			}
		})
	}
}