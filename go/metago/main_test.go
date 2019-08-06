package metago_test

import (
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
			p, err := metago.NewProcessor(&metago.Config{
				TargetPackages: []string{
					packagePath,
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			result, err := p.Process()
			if err != nil {
				if result != nil {
					for _, pErr := range result.CompileErrors {
						t.Error(pErr)
					}
				}
				t.Fatal(err)
			}

			for _, fileResult := range result.Results {
				fileResult := fileResult
				baseFileName := strings.TrimSuffix(filepath.Base(fileResult.FilePath), ".go")

				t.Run(baseFileName, func(t *testing.T) {
					expectedCode := fileResult.GeneratedCode

					generatedFilePath := path.Join(fixturePath, fileInfo.Name(), fmt.Sprintf("%s_gen.go", baseFileName))
					b, err := ioutil.ReadFile(generatedFilePath)
					if *update || os.IsNotExist(err) {
						t.Logf("update %s", generatedFilePath)
						err = ioutil.WriteFile(generatedFilePath, []byte(expectedCode), 0666)
						if err != nil {
							t.Fatal(err)
						}
						b = []byte(expectedCode)
					} else if err != nil {
						t.Fatal(err)
					}

					if expectedCode == string(b) {
						return
					}

					diff := difflib.UnifiedDiff{
						A:       difflib.SplitLines(expectedCode),
						B:       difflib.SplitLines(string(b)),
						Context: 5,
					}
					d, err := difflib.GetUnifiedDiffString(diff)
					if err != nil {
						t.Fatal(err)
					}
					t.Error(d)
				})
			}
		})
	}
}
