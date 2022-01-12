import * as React from "react"
import { Config } from "../hooks/use_config"

const DashboardComponent = ({ config }: {config: Config}) => {
  return (
    <>
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

export default DashboardComponent