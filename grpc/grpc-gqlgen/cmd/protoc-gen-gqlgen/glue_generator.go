package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/rakyll/statik/fs"
	"golang.org/x/xerrors"
)

func glueGenerate(ctx context.Context, b *Builder) ([]*plugin.CodeGeneratorResponse_File, error) {

	fileInfos := b.FileInfos

	tmplBytes, err := ioutil.ReadFile("./tmpls/glue.gotmpl")
	if os.IsNotExist(err) {
		statikFS, err := fs.New()
		if err != nil {
			return nil, xerrors.Errorf("on fs.New: %w", err)
		}
		f, err := statikFS.Open("/glue.gotmpl")
		if err != nil {
			return nil, xerrors.Errorf("on statikFS.Open: %w", err)
		}
		tmplBytes, err = ioutil.ReadAll(f)
		if err != nil {
			return nil, xerrors.Errorf("on ioutil.ReadAll: %w", err)
		}

	} else if err != nil {
		return nil, xerrors.Errorf("on read glue.gotmpl: %w", err)
	}
	tmpl, err := template.
		New("glue").
		Funcs(map[string]interface{}{
			"first": func(ss ...string) string {
				for _, s := range ss {
					if s != "" {
						return s
					}
				}
				return ""
			},
			"goName": func(name string) string {
				return templates.ToGo(name)
			},
			"goNamePrivate": func(name string) string {
				return templates.ToGoPrivate(name)
			},
		}).
		Parse(string(tmplBytes))
	if err != nil {
		return nil, xerrors.Errorf("on parse template: %w", err)
	}

	var files []*plugin.CodeGeneratorResponse_File
	for _, fileInfo := range fileInfos {
		err = fileInfo.Prepare()
		if err != nil {
			return nil, xerrors.Errorf("%s on fileInfo.Prepare: %w", fileInfo.PackageName, err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, fileInfo)
		if err != nil {
			return nil, xerrors.Errorf("%s on tmpl.Execute: %w", fileInfo.PackageName, err)
		}

		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fmt.Sprintf("%s.gql.go", fileInfo.PackageName)),
			Content: proto.String(buf.String()),
		})
	}

	return files, nil
}
