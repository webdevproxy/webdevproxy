import * as React from "react"
import { Link } from "react-router-dom"
import { useConfigContext } from "../hooks/use_config"
import HostsComponent from "./hosts"

const DashboardComponent = () => {
  return (
    <>
      <HostsComponent />
    </>
  )
}

export default DashboardComponent