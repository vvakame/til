package metago_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/vvakame/til/go/metago"
)

const basePath = "github.com/vvakame/til/go/metago"
const fixturePath = "./internal/testbed"

var update = flag.Bool("update", false, "update gen files")

func TestProcessor(t *testing.T) {

	fileInfos, err := ioutil.ReadDir(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			continue
		}

		fileInfo := fileInfo
		packagePath := path.Join(basePath, fixturePath, fileInfo.Name())

		t.Run(packagePath, func(t *testing.T) {
			p, err := metago.NewProcessor()
			if err != nil {
				t.Fatal(err)
			}
			result, err := p.Process(&metago.Config{
				TargetPackages: []string{
					packagePath,
				},
			})
			if result != nil {
				// テスト継続可

			} else if err != nil {
				if result != nil {
					for _, pErr := range result.CompileErrors {
						t.Error(pErr)
					}
				}
				t.Fatal(err)
			}

			for _, fileResult := range result.Results {
				fileResult := fileResult
				baseFileName := strings.TrimSuffix(filepath.Base(fileResult.BaseFilePath), ".go")

				t.Run(baseFileName, func(t *testing.T) {
					t.Run("code", func(t *testing.T) {
						if fileResult.GeneratedCode == "" {
							t.Skip()
						}

						actualCode := fileResult.GeneratedCode
						generatedFilePath := path.Join(fixturePath, fileInfo.Name(), fmt.Sprintf("%s_gen.go", baseFileName))
						b, err := ioutil.ReadFile(generatedFilePath)
						if *update || os.IsNotExist(err) {
							t.Logf("update %s", generatedFilePath)
							err = ioutil.WriteFile(generatedFilePath, []byte(actualCode), 0666)
							if err != nil {
								t.Fatal(err)
							}
							b = []byte(actualCode)
						} else if err != nil {
							t.Fatal(err)
						}

						if actualCode == string(b) {
							return
						}

						diff := difflib.UnifiedDiff{
							A:       difflib.SplitLines(actualCode),
							B:       difflib.SplitLines(string(b)),
							Context: 5,
						}
						d, err := difflib.GetUnifiedDiffString(diff)
						if err != nil {
							t.Fatal(err)
						}
						t.Error(d)
					})
					t.Run("diagnostic", func(t *testing.T) {
						actualBytes, err := json.MarshalIndent(fileResult.Errors, "", "  ")
						if err != nil {
							t.Fatal(err)
						}
						diagnosticFilePath := path.Join(fixturePath, fileInfo.Name(), fmt.Sprintf("%s_diagnostic.json", baseFileName))
						b, err := ioutil.ReadFile(diagnosticFilePath)
						if *update || os.IsNotExist(err) {
							t.Logf("update %s", diagnosticFilePath)
							err = ioutil.WriteFile(diagnosticFilePath, actualBytes, 0666)
							if err != nil {
								t.Fatal(err)
							}
							b = actualBytes
						} else if err != nil {
							t.Fatal(err)
						}

						if string(actualBytes) == string(b) {
							return
						}

						diff := difflib.UnifiedDiff{
							A:       difflib.SplitLines(string(b)),
							B:       difflib.SplitLines(string(actualBytes)),
							Context: 5,
						}
						d, err := difflib.GetUnifiedDiffString(diff)
						if err != nil {
							t.Fatal(err)
						}
						t.Error(d)
					})
				})
			}
		})
	}
}
