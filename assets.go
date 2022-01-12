package main

import (
	"embed"
	"io/fs"
	"text/template"
)

//go:embed assets/*
var assets embed.FS

var parsedProxyErrorTemplate *template.Template

func proxyErrorTemplate() *template.Template {
	if parsedProxyErrorTemplate == nil {
		parsedProxyErrorTemplate, _ = template.ParseFS(assets, "assets/proxy/error.html")
	}
	return parsedProxyErrorTemplate
}

func adminContent() fs.FS {
	content, _ := fs.Sub(assets, "assets/admin")
	return content
}
