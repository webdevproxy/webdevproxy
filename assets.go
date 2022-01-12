package main

import (
	"embed"
	"io/fs"
	"os"
	"path"
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

func adminFS() fs.FS {
	liveAdminPath, useLiveAdminPath := checkForLiveAdminPath()
	if useLiveAdminPath {
		return os.DirFS(liveAdminPath)
	}
	content, _ := fs.Sub(assets, "assets/admin")
	return content
}

func checkForLiveAdminPath() (string, bool) {
	dir, err := os.Getwd()
	if err == nil {
		adminPath := path.Join(dir, "assets", "admin")
		pathInfo, err := os.Stat(adminPath)
		if err == nil && pathInfo.IsDir() {
			return adminPath, true
		}
	}
	return "", false
}
