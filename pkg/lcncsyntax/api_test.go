package lcncsysntax

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestApiParsing(t *testing.T) {
	tests := []struct {
		src string
	}{
		{src: "./../../examples"},
	}
	for _, test := range tests {
		files, err := os.ReadDir(test.src)
		if err != nil {
			t.Error(err)
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".yaml") {
				t.Log(f.Name())
				b, err := os.ReadFile(filepath.Join(test.src, f.Name()))
				if err != nil {
					t.Error(err)
				}

				cfg := &LcncConfig{}
				if err := yaml.Unmarshal(b, cfg); err != nil {
					t.Error(err)
				}
				t.Logf("\nfor: %s\n", cfg.For.Gvk)
				for _, v := range cfg.Vars {
					t.Logf("variable: %s, query: %v, for: %v \n", v.VariableName, v.LcncQuery, v.For)
				}
				for _, r := range cfg.Resources {
					output := []string{}
					for _, o := range r.Output {
						output = append(output, o.LcncGvk.Gvk)
					}
					t.Logf("resource image: %v, input: %v, output: %v \n", r.LcncImage.ImageName, r.Input, output)
				}
			}
		}
	}
}
