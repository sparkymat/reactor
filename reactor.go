package reactor

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sparkymat/webdsl/html"
)

type fileMapping struct {
	filePath string
	webPath  string
}

type reactor struct {
	appId                 string
	name                  string
	javascriptFolders     []fileMapping
	cssFolders            []fileMapping
	customJavascriptLinks []string
	customCssLinks        []string
}

func New(name string, appId string) reactor {
	r := reactor{name: name, appId: appId}
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
	javascriptPaths := []string{}
	for _, jsMap := range r.javascriptFolders {
		err := filepath.Walk(jsMap.filePath, func(path string, f os.FileInfo, err error) error {
			if filepath.Ext(f.Name()) == ".js" {
				javascriptPaths = append(javascriptPaths, strings.Replace(path, jsMap.filePath, jsMap.webPath, 1))
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	for _, javascriptPath := range javascriptPaths {
		javascriptNodes = append(javascriptNodes, html.Script().Attr("type", "text/javascript").Attr("src", javascriptPath))
	}

	for _, jsLink := range r.customJavascriptLinks {
		headNodes = append(headNodes, html.Script().Attr("type", "text/javascript").Attr("src", jsLink))
	}

	// List css files
	cssPaths := []string{}
	for _, cssMap := range r.cssFolders {
		err := filepath.Walk(cssMap.filePath, func(path string, f os.FileInfo, err error) error {
			if filepath.Ext(f.Name()) == ".css" {
				cssPaths = append(cssPaths, strings.Replace(path, cssMap.filePath, cssMap.webPath, 1))
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	for _, cssPath := range cssPaths {
		headNodes = append(headNodes, html.Link().Attr("rel", "stylesheet").Attr("href", cssPath))
	}

	for _, cssLink := range r.customCssLinks {
		headNodes = append(headNodes, html.Link().Attr("rel", "stylesheet").Attr("href", cssLink))
	}

	childNodes := []*html.Node{html.Div().Attr("id", r.appId)}

	childNodes = append(childNodes, javascriptNodes...)

	return html.Html(
		html.Head(headNodes...),
		html.Body(childNodes...),
	)
}
