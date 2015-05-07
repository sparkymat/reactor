package reactor

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/kirillrdy/nadeshiko/html"
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
	props             map[string]string
}

func New(name string) reactor {
	r := reactor{name: name}
	return r
}

func (r *reactor) SetInitialProperties(props map[string]string) {
	// Keep ref or copy?
	r.props = props
}

func (r *reactor) MapJavascriptFolder(file string, web string) {
	r.javascriptFolders = append(r.javascriptFolders, fileMapping{filePath: file, webPath: web})
}

func (r *reactor) MapCssFolder(file string, web string) {
	r.cssFolders = append(r.cssFolders, fileMapping{filePath: file, webPath: web})
}

func (r reactor) Html() html.Node {
	headNodes := []html.Node{}

	headNodes = append(headNodes, html.Title().Text(r.name))
	headNodes = append(headNodes, html.Script().Type("text/javascript").Src("https://code.jquery.com/jquery-2.1.4.min.js"))
	headNodes = append(headNodes, html.Script().Type("text/javascript").Src("https://cdnjs.cloudflare.com/ajax/libs/lodash.js/3.8.0/lodash.min.js"))
	headNodes = append(headNodes, html.Script().Type("text/javascript").Src("https://cdnjs.cloudflare.com/ajax/libs/react/0.13.2/react.min.js"))

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

	propsList := []string{}

	if r.props != nil {
		for key, value := range r.props {
			propsList = append(propsList, fmt.Sprintf("%v: \"%v\"", key, value))
		}
	}

	propsString := fmt.Sprintf("{%v}", strings.Join(propsList, ","))

	return html.Html().Children(
		html.Head().Children(headNodes...),
		html.Body().Children(
			html.Div().Id(css.Id("app-container")),
			html.Script().TextUnsafe(fmt.Sprintf("React.render(React.createElement(%v, %v), document.getElementById('app-container'));", r.name, propsString)),
		),
	)
}
