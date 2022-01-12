import { useEffect } from "react"

export const useLiveAdmin = () => {
  useEffect(() => {
    if (window.localStorage.getItem("__live_admin")) {
      console.info("Listening for live admin updates...")

      const source = new EventSource("/__live_admin")
      source.addEventListener("message", (e) => {
        if (/\.css$/.test(e.data)) {
          console.info("Reloading css...")

          const queryString = `?reload=${new Date().getTime()}`
          document.querySelectorAll('head > link').forEach(link => {
            const href = link.getAttribute("href")
            if (href) {
              const newHref = href?.replace(/\?.*|$/, queryString)
              link.setAttribute("href", newHref)
            }
          })
        } else {
          window.location.reload()
        }
      })

      return () => source.close()
    }
  }, [])
}
