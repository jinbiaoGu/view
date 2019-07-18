package view

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// Template template struct
type Template struct {
	view             *View
	layout             string
	usingDefaultLayout bool
	funcMap            template.FuncMap
}

// FuncMap get func maps from tmpl
func (tmpl *Template) funcMapMaker(writer http.ResponseWriter,req *http.Request) template.FuncMap {
	var funcMap = template.FuncMap{}

	for key, fc := range tmpl.view.funcMaps {
		funcMap[key] = fc
	}

	if tmpl.view.Config.FuncMapMaker != nil {
		for key, fc := range tmpl.view.Config.FuncMapMaker(tmpl.view, writer, req) {
			funcMap[key] = fc
		}
	}

	for key, fc := range tmpl.funcMap {
		funcMap[key] = fc
	}
	return funcMap
}

// Funcs register Funcs for tmpl
func (tmpl *Template) Funcs(funcMap template.FuncMap) *Template {
	tmpl.funcMap = funcMap
	return tmpl
}

// View view tmpl
func (tmpl *Template) View(templateName string, obj interface{}, writer http.ResponseWriter, request *http.Request) (template.HTML, error) {
	var (
		content []byte
		t       *template.Template
		err     error
		funcMap = tmpl.funcMapMaker(writer,request)
		view  = func(name string, objs ...interface{}) (template.HTML, error) {
			var (
				err           error
				viewObj     interface{}
				viewContent []byte
			)

			if len(objs) == 0 {
				// default obj
				viewObj = obj
			} else {
				// overwrite obj
				for _, o := range objs {
					viewObj = o
					break
				}
			}
			if viewContent, err = tmpl.findTemplate(name); err == nil {
				var partialTemplate *template.Template
				result := bytes.NewBufferString("")
				if partialTemplate, err = template.New(filepath.Base(name)).Funcs(funcMap).Parse(string(viewContent)); err == nil {
					if err = partialTemplate.Execute(result, viewObj); err == nil {
						return template.HTML(result.String()), err
					}
				}
			} else {
				err = fmt.Errorf("failed to find template: %v", name)
			}

			if err != nil {
				fmt.Println(err)
			}

			return "", err
		}
	)

	// funcMaps
	funcMap["view"] = view
	funcMap["yield"] = func() (template.HTML, error) { return view(templateName) }

	layout := tmpl.layout

	usingDefaultLayout := false

	if layout == "" && tmpl.usingDefaultLayout {
		usingDefaultLayout = true
		layout = tmpl.view.DefaultLayout

		log.Println(layout)
	}

	if layout != "" {
		content, err = tmpl.findTemplate(filepath.Join("layouts", layout))
		if err == nil {
			if t, err = template.New("").Funcs(funcMap).Parse(string(content)); err == nil {
				var tpl bytes.Buffer
				if err = t.Execute(&tpl, obj); err == nil {
					return template.HTML(tpl.String()), nil
				}
			}
		} else if usingDefaultLayout {
			err = fmt.Errorf("Failed to view layout: '%v', got error: %v", filepath.Join("layouts", tmpl.layout), err)
			return template.HTML(""), err
		}
	}

	if content, err = tmpl.findTemplate(templateName); err == nil {
		if t, err = template.New("").Funcs(funcMap).Parse(string(content)); err == nil {
			var tpl bytes.Buffer
			if err = t.Execute(&tpl, obj); err == nil {
				return template.HTML(tpl.String()), nil
			}
		}
	} else {
		err = fmt.Errorf("failed to find template: %v", templateName)
	}

	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(""), err
}

// Execute execute tmpl
func (tmpl *Template) Execute(templateName string, obj interface{}, w http.ResponseWriter, req *http.Request) error {
	result, err := tmpl.View(templateName, obj, w, req)
	if err == nil {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html")
		}

		_, err = w.Write([]byte(result))
	}
	return err
}

func (tmpl *Template) findTemplate(name string) ([]byte, error) {
	return tmpl.view.Asset(strings.TrimSpace(name) + tmpl.view.ViewExt)

}
