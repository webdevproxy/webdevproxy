import * as React from "react"
import * as ReactDOM from "react-dom"
import { BrowserRouter, Routes, Route } from "react-router-dom";

import AppComponent from "./components/app"
import DashboardComponent from "./components/dashboard";
import HostComponent from "./components/host";
import HostsComponent from "./components/hosts";

ReactDOM.render(
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<AppComponent />}>
          <Route index element={<DashboardComponent />} />
          <Route path="hosts" element={<HostsComponent />} />
          <Route path="hosts/:hostId" element={<HostComponent />} />
          <Route
            path="*"
            element={
              <main style={{ padding: "1rem" }}>
                <p>There's nothing here!</p>
              </main>
            }
          />
        </Route>
      </Routes>
    </BrowserRouter>, document.getElementById("app"))
