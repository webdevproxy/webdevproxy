import * as React from "react"

import { useConfig } from "../hooks/use_config"
import { useLiveAdmin } from "../hooks/use_live_admin"
import DashboardComponent from "./dashboard"

const AppComponent = () => {
  useLiveAdmin()
  const { config, configError } = useConfig()

  if (!config) {
    return <div className="centered">Loading configuration</div>
  }

  return (
    <>
      {configError && <div className="error">{configError}</div>}
      <DashboardComponent config={config} />
    </>
  )
}

export default AppComponent