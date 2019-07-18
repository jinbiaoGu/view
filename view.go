// Package view support to view templates by your control.
package view

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// DefaultLayout default layout name
const DefaultLayout = "application"

// DefaultViewPath default view path
const DefaultViewPath = "views"

//DefaultView file name File extension
const DefaultViewExt = ".html"

// Config view config
type Config struct {
	ViewExt string
	ViewPaths       []string
	DefaultLayout   string
	FuncMapMaker    func(view *View, writer http.ResponseWriter, request *http.Request) template.FuncMap
}

// View the view struct.
type View struct {
	*Config
	funcMaps template.FuncMap
}

// New initalize the View struct.
func New(config *Config, viewPaths ...string) *View {
	if config == nil {
		config = &Config{}
	}

	if config.DefaultLayout == "" {
		config.DefaultLayout = DefaultLayout
	}

	if config.ViewExt == "" {
		config.ViewExt = DefaultViewExt
	}
	config.ViewPaths = append(append(config.ViewPaths, viewPaths...), DefaultViewPath)

	view := &View{funcMaps: map[string]interface{}{}, Config: config}

	for _, viewPath := range config.ViewPaths {
		view.RegisterViewPath(viewPath)
	}
	return view
}

// RegisterViewPath register view path
func (view *View) RegisterViewPath(paths ...string) {
	tmpFilePath := ""
	for _, pth := range paths {

		if !filepath.IsAbs(pth) {
			if isExistingDir(filepath.Join(getAppRoot(), "vendor", pth)) {
				tmpFilePath =  filepath.Join(getAppRoot(), "vendor", pth)
			}

			for _, goPath := range GOPATH() {
				if p := filepath.Join(goPath, "src", pth); isExistingDir(p) {
					tmpFilePath = p
				}
			}

			if absPath, err := filepath.Abs(pth); err == nil && isExistingDir(absPath) {
				tmpFilePath = absPath
			}

			if p := filepath.Join(getAppRoot(), pth); isExistingDir(p) {
				tmpFilePath = p
			}

		}

	}

	if tmpFilePath != "" {
		view.ViewPaths = append(view.ViewPaths, tmpFilePath)
	}
}


// Layout set layout for template.
func (view *View) Layout(name string) *Template {
	return &Template{view: view, layout: name}
}

// Funcs set helper functions for template with default "application" layout.
func (view *View) Funcs(funcMap template.FuncMap) *Template {
	tmpl := &Template{view: view, usingDefaultLayout: true}
	return tmpl.Funcs(funcMap)
}

// Execute view template with default "application" layout.
func (view *View) Execute(name string, context interface{}, writer http.ResponseWriter, request *http.Request) error {
	tmpl := &Template{view: view, usingDefaultLayout: true}
	return tmpl.Execute(name, context, writer, request)
}

// RegisterFuncMap register FuncMap for view.
func (view *View) RegisterFuncMap(name string, fc interface{}) {
	if view.funcMaps == nil {
		view.funcMaps = template.FuncMap{}
	}
	view.funcMaps[name] = fc
}

// Asset get content from
func (view *View) Asset(name string) ([]byte, error) {
	for _, pth := range view.ViewPaths {
		if _, err := os.Stat(filepath.Join(pth, name)); err == nil {
			return ioutil.ReadFile(filepath.Join(pth, name))
		}
	}
	return []byte{}, fmt.Errorf("%v not found", name)
}
