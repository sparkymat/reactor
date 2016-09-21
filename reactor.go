package reactor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/kirillrdy/nadeshiko/html"
	"github.com/sparkymat/reactor/property"
	"github.com/sparkymat/webdsl/css"
)

type fileMapping struct {
	filePath string
	webPath  string
}

type reactor struct {
	name              string
	javascriptFolders []fileMapping
	cssFolders        []fileMapping
	customJavascriptLinks []string
  customCssLinks    []string
	props             map[string]Property
}

func New(name string) reactor {
	r := reactor{name: name}
	r.props = make(map[string]Property)
	return r
}

func (r *reactor) SetStringProperty(name string, value string) {
	property := Property{propertyType: property.String, value: value}
	r.props[name] = property
}

func (r *reactor) SetIntegerProperty(name string, value int64) {
	property := Property{propertyType: property.Integer, value: fmt.Sprintf("%v", value)}
	r.props[name] = property
}

func (r *reactor) SetFloatProperty(name string, value float64) {
	property := Property{propertyType: property.Float, value: fmt.Sprintf("%v", value)}
	r.props[name] = property
}

func (r *reactor) SetObjectProperty(name string, value interface{}) {
	jsonValue, _ := json.Marshal(value)

	// FIXME: Ignoring error for now
	property := Property{propertyType: property.Object, value: strings.Replace(string(jsonValue), "\"", "\\\"\\", -1)}
	r.props[name] = property
}

func (r *reactor) MapJavascriptFolder(file string, web string) {
	r.javascriptFolders = append(r.javascriptFolders, fileMapping{filePath: file, webPath: web})
}

func (r *reactor) MapCssFolder(file string, web string) {
	r.cssFolders = append(r.cssFolders, fileMapping{filePath: file, webPath: web})
}

func (r *reactor) AddCustomJavascriptLink(link string) {
	r.customJavascriptLinks = append(r.customJavascriptLinks, link)
}

func (r *reactor) AddCustomCssLink(link string) {
	r.customCssLinks = append(r.customCssLinks, link)
}

func (r reactor) Html() html.Node {
	headNodes := []html.Node{}

	headNodes = append(headNodes, html.Title().Text(r.name))

	// List javascript files
	for _, jsMap := range r.javascriptFolders {
		files, err := ioutil.ReadDir(jsMap.filePath)
		if err == nil {
			for _, jsFile := range files {
				if filepath.Ext(jsFile.Name()) == ".js" {
					finalPath := fmt.Sprintf("/%v/%v", jsMap.webPath, jsFile.Name())
					headNodes = append(headNodes, html.Script().Type("text/javascript").Src(finalPath))
				}
			}
		}
	}

	for _, jsLink := range r.customJavascriptLinks {
		headNodes = append(headNodes, html.Script().Type("text/javascript").Src(jsLink))
	}

	// List css files
	for _, cssMap := range r.cssFolders {
		files, err := ioutil.ReadDir(cssMap.filePath)

		if err == nil {
			for _, cssFile := range files {
				if filepath.Ext(cssFile.Name()) == ".css" {
					finalPath := fmt.Sprintf("/%v/%v", cssMap.webPath, cssFile.Name())
					headNodes = append(headNodes, html.Link().Rel("stylesheet").Href(finalPath))
				}
			}
		}
	}

	for _, cssLink := range r.customCssLinks {
		headNodes = append(headNodes, html.Link().Rel("stylesheet").Href(cssLink))
	}

	propsList := []string{}

	if r.props != nil {
		for key, value := range r.props {
			propsList = append(propsList, fmt.Sprintf("%v: %v", key, value.String()))
		}
	}

	propsString := fmt.Sprintf("{%v}", strings.Join(propsList, ","))

	return html.Html().Children(
		html.Head().Children(headNodes...),
		html.Body().Children(
			html.Div().Id(css.Id("app-container")),
			html.Script().TextUnsafe(fmt.Sprintf("ReactDOM.render(React.createElement(%v, %v), document.getElementById('app-container'));", r.name, propsString)),
		),
	)
}
