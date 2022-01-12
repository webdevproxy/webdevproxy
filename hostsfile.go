package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type HostsFileEntry struct {
	LineNumber  int    `json:"lineNumber"`
	LineContent string `json:"lineContent"`
	Ip          string `json:"ip"`
	Host        string `json:"host"`
	Proxied     bool   `json:"proxied"`
	ProxyIp     string `json:"proxyIp"`
	ProxyPort   int    `json:"proxyPort"`
	ProxyHost   string `json:"proxyHost"`
}

type HostsFileSyntaxError struct {
	LineNumber  int    `json:"lineNumber"`
	LineContent string `json:"lineContent"`
	SyntaxError string `json:"syntaxError"`
}

type Hostsfile struct {
	Path         string                 `json:"path"`
	Contents     string                 `json:"contents"`
	Entries      []HostsFileEntry       `json:"entries"`
	SyntaxErrors []HostsFileSyntaxError `json:"syntaxErrors"`
}

func (h *HostsFileEntry) proxyUrl() *url.URL {
	url, _ := url.Parse(fmt.Sprintf("http://%s:%d/", h.ProxyIp, h.ProxyPort))
	return url
}

func NewHostsFile(path string) *Hostsfile {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		contents = make([]byte, 0)
	}

	entries, syntaxErrors := ParseHostsFile(contents)

	hostsfile := Hostsfile{Path: path, Contents: string(contents), Entries: entries, SyntaxErrors: syntaxErrors}
	return &hostsfile
}

func ParseHostsFile(contents []byte) ([]HostsFileEntry, []HostsFileSyntaxError) {
	existingHosts := make(map[string]bool)
	entries := make([]HostsFileEntry, 0)
	syntaxErrors := make([]HostsFileSyntaxError, 0)
	proxyComment := regexp.MustCompile("^\\s*webdevproxy(.*)$")

	lines := strings.Split(strings.Trim(string(contents), " \t\r\n"), "\n")
	for lineIndex, lineContent := range lines {
		lineNumber := lineIndex + 1

		// trim and skip empty lines or lines just with comments
		trimmedLine := strings.Replace(strings.Trim(lineContent, " \t"), "\t", " ", -1)
		if len(trimmedLine) == 0 || trimmedLine[0] == '#' {
			continue
		}

		lineParts := strings.SplitN(trimmedLine, " ", 2)
		if len(lineParts) > 1 && len(lineParts[0]) > 0 {
			lineEntries := make([]HostsFileEntry, 0)

			// parse the text before any comment in form (ip host*)
			ip := lineParts[0]
			valueParts := strings.SplitN(lineParts[1], "#", 2)
			if hosts := strings.Fields(valueParts[0]); len(hosts) > 0 {
				for _, host := range hosts {
					if host != "localhost" {
						if existingHosts[host] {
							syntaxErrors = append(syntaxErrors, HostsFileSyntaxError{
								LineNumber:  lineNumber,
								LineContent: lineContent,
								SyntaxError: "Duplicate hostname",
							})
						} else {
							lineEntries = append(lineEntries, HostsFileEntry{LineNumber: lineNumber, LineContent: lineContent, Ip: ip, Host: host})
							existingHosts[host] = true
						}
					}
				}
			}

			// parse the comment(s)
			if len(lineEntries) > 0 && len(valueParts) > 1 {
				var proxied bool
				var proxyPort int
				var proxyHost string
				var proxyIp string

				matches := proxyComment.FindAllStringSubmatch(valueParts[1], -1)
				if len(matches) > 0 && len(matches[0]) > 1 {
					if args := strings.Fields(matches[0][1]); len(args) > 0 {
						for _, arg := range args {
							var err error
							argParts := strings.SplitN(arg, ":", 2)
							if len(argParts) > 1 {
								key := argParts[0]
								value := argParts[1]
								switch {
								case key == "to":
									proxyPort, proxyIp, err = ParseProxyPort(value)
									if err != nil {
										syntaxErrors = append(syntaxErrors, HostsFileSyntaxError{
											LineNumber:  lineNumber,
											LineContent: lineContent,
											SyntaxError: err.Error(),
										})
									} else {
										proxied = true
									}

								case key == "host":
									proxyHost = value

								default:
									syntaxErrors = append(syntaxErrors, HostsFileSyntaxError{
										LineNumber:  lineNumber,
										LineContent: lineContent,
										SyntaxError: fmt.Sprintf("Unknown key:value in the webdevproxy comment: %s", arg),
									})
								}
							}
						}
					}

					if proxied {
						if proxyIp == "" {
							proxyIp = ip
						}
						for i := 0; i < len(lineEntries); i++ {
							lineEntry := &lineEntries[i]
							lineEntry.Proxied = true
							lineEntry.ProxyPort = proxyPort
							lineEntry.ProxyIp = proxyIp
							if proxyHost != "" {
								lineEntry.ProxyHost = proxyHost
							} else {
								lineEntry.ProxyHost = lineEntry.Host
							}
						}
					} else {
						syntaxErrors = append(syntaxErrors, HostsFileSyntaxError{
							LineNumber:  lineNumber,
							LineContent: lineContent,
							SyntaxError: "No port value given, use to:PORT in webdevproxy comment",
						})

					}
				}
			}

			// add the line entries
			entries = append(entries, lineEntries...)
		}
	}

	return entries, syntaxErrors
}

func ParseProxyPort(value string) (int, string, error) {
	portParts := strings.SplitN(value, ":", 2)
	if len(portParts) > 2 {
		return 0, "", errors.New("The to: value must either be PORT or IP:PORT")
	}

	var ip string
	var port string

	if len(portParts) == 1 {
		port = portParts[0]
	}
	if len(portParts) == 2 {
		ip = portParts[0]
		port = portParts[1]
	}

	portValue, err := strconv.Atoi(port)
	if err != nil {
		return 0, "", errors.New(fmt.Sprintf("Invalid to: port value, must be an integer: %s", port))
	}
	if portValue < 1 {
		return 0, "", errors.New(fmt.Sprintf("Invalid to: port value, must be >= 1: %s", port))
	}

	return portValue, ip, nil
}
