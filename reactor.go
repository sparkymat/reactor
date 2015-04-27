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
}

func New(name string) reactor {
	r := reactor{name: name}
	return r
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
					finalPath := fmt.Sprintf("%v/%v", cssMap.filePath, cssFile.Name())
					finalPath = strings.Replace(finalPath, cssMap.filePath, cssMap.webPath, 1)
					headNodes = append(headNodes, html.Link().Rel("stylesheet").Href(finalPath))
				}
			}
		}
	}

	return html.Html().Children(
		html.Head().Children(headNodes...),
		html.Body().Children(
			html.Div().Id(css.Id("app-container")),
			html.Script().TextUnsafe(fmt.Sprintf("React.render(React.createElement(%v, null), document.getElementById('app-container'));", r.name)),
		),
	)
}
