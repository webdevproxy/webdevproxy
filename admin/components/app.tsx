import * as React from "react"
import { Link, Outlet } from "react-router-dom"

import { ConfigContext, useConfig } from "../hooks/use_config"
import { useLiveAdmin } from "../hooks/use_live_admin"

const AppComponent = () => {
  useLiveAdmin()
  const { config, configError } = useConfig()

  if (!config) {
    return <div className="centered">Loading configuration</div>
  }

  return (
    <ConfigContext.Provider value={config}>
      {configError && <div className="error">{configError}</div>}
      <p>
        <Link to="/">Dashboard</Link> | <Link to="/hosts">All Hosts</Link>
      </p>
      <Outlet />
    </ConfigContext.Provider>
  )
}

export default AppComponent