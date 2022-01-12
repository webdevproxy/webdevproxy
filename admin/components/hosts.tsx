import * as React from "react"
import { Link, Outlet } from "react-router-dom"
import { useConfigContext } from "../hooks/use_config"

const HostsComponent = () => {
  const config = useConfigContext()

  return (
    <>
      <table>
        <thead>
          <tr>
            <th>Host</th>
            <th>Proxies To</th>
          </tr>
        </thead>
        <tbody>
          {config?.hosts.entries.map((e, index) => (
            <tr key={`${e.host}-${e.lineNumber}-${index}`}>
              <td>
                <a href={`http://${e.host}`}>↗</a>
                <Link to={`/hosts/${e.host}`} key={e.host}>{e.host}</Link>
              </td>
              <td>
                {e.proxied ? `${e.proxyIp}:${e.proxyPort}` : "✗"}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <Outlet />
    </>
  )
}

export default HostsComponent