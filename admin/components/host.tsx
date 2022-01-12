import * as React from "react"
import { useParams } from "react-router-dom"
import { useConfigContext, useHostsFileEntry } from "../hooks/use_config"

const HostComponent = () => {
  const { hostId } = useParams()
  const hostFileEntry = useHostsFileEntry(hostId)

  return (
    <div>
      {hostId}
      <p>
        <pre>{JSON.stringify(hostFileEntry, null, 2)}</pre>
      </p>
    </div>
  )
}

export default HostComponent