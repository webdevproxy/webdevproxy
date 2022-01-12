import { useEffect, useState } from "react"

interface ProxyConfig {
}

interface HostsFileEntry {
  lineNumber:    number
  lineContent:   string
	ip:            string
	host:          string
	proxied:       boolean
	proxyPort:     number
	proxyIp:       string
	proxyHost:     string
}

interface HostsFileSyntaxError {
	lineNumber:  number
	lineContent: string
	syntaxError: string
}

interface Hostsfile {
  path: string
  contents: string
  entries: HostsFileEntry[]
  syntaxErrors: HostsFileSyntaxError[]
}

interface Config {
  proxy: ProxyConfig,
  hosts: Hostsfile
}

export const useConfig = () => {
  const [config, setConfig] = useState<Config | undefined>(undefined)
  const [configError, setConfigError] = useState<string|undefined>(undefined)

  useEffect(() => {
    const source = new EventSource("/api/watch-config")
    source.addEventListener("message", (e) => {
      try {
        setConfig(JSON.parse(e.data))
        setConfigError(undefined)
      } catch (err: any) {
        setConfigError(err.toString())
      }
    })
    source.addEventListener("error", () => {
      setConfigError("Error reading config from local proxy")
    })

    return () => source.close()
  }, [])
  return { config, configError }
}
