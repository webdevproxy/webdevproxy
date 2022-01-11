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
	LineNumber    int
	Ip            string
	Hostname      string
	Proxied       bool
	ProxyPort     int
	ProxyHost     string
	ProxyHostname string
}

type HostsFileSyntaxError struct {
	LineNumber  int
	Line        string
	SyntaxError string
}

type Hostsfile struct {
	Path         string
	Contents     string
	Entries      []HostsFileEntry
	SyntaxErrors []HostsFileSyntaxError
}

func (h *HostsFileEntry) proxyUrl() *url.URL {
	url, _ := url.Parse(fmt.Sprintf("http://%s:%d/", h.ProxyHost, h.ProxyPort))
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
	entries := make([]HostsFileEntry, 0)
	syntaxErrors := make([]HostsFileSyntaxError, 0)
	proxyComment := regexp.MustCompile("^\\s*webdevproxy(.*)$")

	lines := strings.Split(strings.Trim(string(contents), " \t\r\n"), "\n")
	for lineIndex, line := range lines {
		lineNumber := lineIndex + 1

		// trim and skip empty lines or lines just with comments
		line = strings.Replace(strings.Trim(line, " \t"), "\t", " ", -1)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		lineParts := strings.SplitN(line, " ", 2)
		if len(lineParts) > 1 && len(lineParts[0]) > 0 {
			lineEntries := make([]HostsFileEntry, 0)

			// parse the text before any comment in form (ip hostname*)
			ip := lineParts[0]
			valueParts := strings.SplitN(lineParts[1], "#", 2)
			if hostnames := strings.Fields(valueParts[0]); len(hostnames) > 0 {
				for _, hostname := range hostnames {
					lineEntries = append(lineEntries, HostsFileEntry{LineNumber: lineNumber, Ip: ip, Hostname: hostname})
				}
			}

			// parse the comment(s)
			if len(valueParts) > 1 {
				var proxied bool
				var proxyPort int
				var proxyHostname string
				var proxyHost string

				matches := proxyComment.FindAllStringSubmatch(valueParts[1], -1)
				if len(matches) > 0 && len(matches[0]) > 1 {
					if args := strings.Fields(matches[0][1]); len(args) > 0 {
						for _, arg := range args {
							var err error
							argParts := strings.SplitN(arg, ":", 2)
							if len(argParts) > 1 {
								name := argParts[0]
								switch {
								case name == "to":
									proxyPort, proxyHost, err = ParseProxyPort(argParts[1])
									if err != nil {
										syntaxErrors = append(syntaxErrors, HostsFileSyntaxError{
											LineNumber:  lineNumber,
											Line:        line,
											SyntaxError: err.Error(),
										})
									} else {
										proxied = true
									}

								case name == "host":
									proxyHostname = argParts[1]

								default:
									syntaxErrors = append(syntaxErrors, HostsFileSyntaxError{
										LineNumber:  lineNumber,
										Line:        line,
										SyntaxError: fmt.Sprintf("Unknown key:value in the webdevproxy comment: %s", arg),
									})
								}
							}
						}
					}

					if proxied {
						if proxyHost == "" {
							proxyHost = ip
						}
						for i := 0; i < len(lineEntries); i++ {
							lineEntry := &lineEntries[i]
							if lineEntry.Hostname == "localhost" {
								syntaxErrors = append(syntaxErrors, HostsFileSyntaxError{
									LineNumber:  lineNumber,
									Line:        line,
									SyntaxError: "A webdevproxy comment is not allowed for localhost",
								})
							} else {
								lineEntry.Proxied = true
								lineEntry.ProxyPort = proxyPort
								lineEntry.ProxyHost = proxyHost
								if proxyHostname != "" {
									lineEntry.ProxyHostname = proxyHostname
								} else {
									lineEntry.ProxyHostname = lineEntry.Hostname
								}
							}
						}
					} else {
						syntaxErrors = append(syntaxErrors, HostsFileSyntaxError{
							LineNumber:  lineNumber,
							Line:        line,
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

// parse to:PORT or to:HOST:PORT
func ParseProxyPort(value string) (int, string, error) {
	portParts := strings.SplitN(value, ":", 2)
	if len(portParts) > 2 {
		return 0, "", errors.New("The to: value must either be PORT or HOST:PORT")
	}

	var port string
	var host string

	if len(portParts) == 1 {
		port = portParts[0]
	}
	if len(portParts) == 2 {
		host = portParts[0]
		port = portParts[1]
	}

	portValue, err := strconv.Atoi(port)
	if err != nil {
		return 0, "", errors.New(fmt.Sprintf("Invalid to: port value, must be an integer: %s", port))
	}
	if portValue < 1 {
		return 0, "", errors.New(fmt.Sprintf("Invalid to: port value, must be >= 1: %s", port))
	}

	return portValue, host, nil
}
