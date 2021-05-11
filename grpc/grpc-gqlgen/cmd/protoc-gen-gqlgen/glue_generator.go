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
	g := &glueGenerator{b: b}
	return g.Generate(ctx, b.GenerateFileInfos)
}

type glueGenerator struct {
	b *Builder
}

func (g *glueGenerator) Generate(ctx context.Context, fileInfos []*FileInfo) ([]*plugin.CodeGeneratorResponse_File, error) {

	headerTmpl, err := g.headerTemplate()
	if err != nil {
		return nil, err
	}
	bodyTmpl, err := g.bodyTemplate()
	if err != nil {
		return nil, err
	}

	var files []*plugin.CodeGeneratorResponse_File
	for _, fileInfo := range fileInfos {
		err = fileInfo.Prepare()
		if err != nil {
			return nil, xerrors.Errorf("%s on fileInfo.Prepare: %w", fileInfo.PackageName, err)
		}

		var file *plugin.CodeGeneratorResponse_File
		fileName := fmt.Sprintf("%s.gql.go", fileInfo.PackageName)
		for _, f := range files {
			if f.GetName() == fileName {
				file = f
				break
			}
		}
		if file == nil {
			file = &plugin.CodeGeneratorResponse_File{
				Name: proto.String(fileName),
			}
			files = append(files, file)

			var buf bytes.Buffer
			err = headerTmpl.Execute(&buf, fileInfo)
			if err != nil {
				return nil, xerrors.Errorf("%s on headerTmpl.Execute: %w", fileInfo.PackageName, err)
			}

			file.Content = proto.String(buf.String())
		}

		var buf bytes.Buffer
		err = bodyTmpl.Execute(&buf, fileInfo)
		if err != nil {
			return nil, xerrors.Errorf("%s on bodyTmpl.Execute: %w", fileInfo.PackageName, err)
		}

		file.Content = proto.String(fmt.Sprintf("%s\n\n\n// ------ %s ------\n\n\n%s", file.GetContent(), fileInfo.Proto.GetFullyQualifiedName(), buf.String()))
	}

	return files, nil
}

func (g *glueGenerator) headerTemplate() (*template.Template, error) {
	tmplBytes, err := ioutil.ReadFile("./tmpls/header.gotmpl")
	if os.IsNotExist(err) {
		statikFS, err := fs.New()
		if err != nil {
			return nil, xerrors.Errorf("on fs.New: %w", err)
		}
		f, err := statikFS.Open("/header.gotmpl")
		if err != nil {
			return nil, xerrors.Errorf("on statikFS.Open: %w", err)
		}
		tmplBytes, err = ioutil.ReadAll(f)
		if err != nil {
			return nil, xerrors.Errorf("on ioutil.ReadAll: %w", err)
		}

	} else if err != nil {
		return nil, xerrors.Errorf("on read header.gotmpl: %w", err)
	}
	tmpl, err := template.
		New("header").
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

	return tmpl, nil
}

func (g *glueGenerator) bodyTemplate() (*template.Template, error) {
	tmplBytes, err := ioutil.ReadFile("./tmpls/body.gotmpl")
	if os.IsNotExist(err) {
		statikFS, err := fs.New()
		if err != nil {
			return nil, xerrors.Errorf("on fs.New: %w", err)
		}
		f, err := statikFS.Open("/body.gotmpl")
		if err != nil {
			return nil, xerrors.Errorf("on statikFS.Open: %w", err)
		}
		tmplBytes, err = ioutil.ReadAll(f)
		if err != nil {
			return nil, xerrors.Errorf("on ioutil.ReadAll: %w", err)
		}

	} else if err != nil {
		return nil, xerrors.Errorf("on read body.gotmpl: %w", err)
	}
	tmpl, err := template.
		New("body").
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

	return tmpl, nil
}
