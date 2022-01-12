import * as React from "react"
import { useConfig } from "../hooks/use_config"

const AppComponent = () => {
  const { config, configError } = useConfig()

  if (!config) {
    return <div className="centered">Loading configuration</div>
  }

  return (
    <>
      {configError && <div className="error">{configError}</div>}
      <table>
        <thead>
          <tr>
            <th>Host</th>
            <th>Proxy To</th>
          </tr>
        </thead>
        <tbody>
          {config.hosts.entries.map((e, index) => (
            <tr key={`${e.host}-${e.lineNumber}-${index}`}>
              <td>{e.proxied ? <a href={`http://${e.host}`}>{e.host}</a> : e.host}</td>
              <td>{e.proxied ? <a href={`http://${e.proxyIp}:${e.proxyPort}`}>{`${e.proxyIp}:${e.proxyPort}`}</a> : "✗"}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  )
}

export default AppComponent