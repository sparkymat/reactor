package reactor

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/sparkymat/webdsl/html"
)

type fileMapping struct {
	filePath string
	webPath  string
}

type reactor struct {
	name                  string
	javascriptFolders     []fileMapping
	cssFolders            []fileMapping
	customJavascriptLinks []string
	customCssLinks        []string
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

func (r *reactor) AddCustomJavascriptLink(link string) {
	r.customJavascriptLinks = append(r.customJavascriptLinks, link)
}

func (r *reactor) AddCustomCssLink(link string) {
	r.customCssLinks = append(r.customCssLinks, link)
}

func (r reactor) Html() html.HtmlDocument {
	javascriptNodes := []*html.Node{}
	headNodes := []*html.Node{}

	headNodes = append(headNodes, html.Title(r.name))

	// List javascript files
	for _, jsMap := range r.javascriptFolders {
		files, err := ioutil.ReadDir(jsMap.filePath)
		if err == nil {
			for _, jsFile := range files {
				if filepath.Ext(jsFile.Name()) == ".js" {
					finalPath := fmt.Sprintf("/%v/%v", jsMap.webPath, jsFile.Name())
					javascriptNodes = append(javascriptNodes, html.Script().Attr("type", "text/javascript").Attr("src", finalPath))
				}
			}
		}
	}

	for _, jsLink := range r.customJavascriptLinks {
		headNodes = append(headNodes, html.Script().Attr("type", "text/javascript").Attr("src", jsLink))
	}

	// List css files
	for _, cssMap := range r.cssFolders {
		files, err := ioutil.ReadDir(cssMap.filePath)

		if err == nil {
			for _, cssFile := range files {
				if filepath.Ext(cssFile.Name()) == ".css" {
					finalPath := fmt.Sprintf("/%v/%v", cssMap.webPath, cssFile.Name())
					headNodes = append(headNodes, html.Link().Attr("rel", "stylesheet").Attr("href", finalPath))
				}
			}
		}
	}

	for _, cssLink := range r.customCssLinks {
		headNodes = append(headNodes, html.Link().Attr("rel", "stylesheet").Attr("href", cssLink))
	}

	childNodes := []*html.Node{html.Div().Attr("id", "js-reactor-app")}

	childNodes = append(childNodes, javascriptNodes...)

	return html.Html(
		html.Head(headNodes...),
		html.Body(childNodes...),
	)
}
