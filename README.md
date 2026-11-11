# Grafana Ping Data Source Plugin

A Grafana backend data source plugin that pings any server and lets you plot
the round-trip time over time: **x = when the check was performed, y = how long
the server took to answer (ms)**.

## How it works

- The plugin runs inside Grafana's backend (Go binary).
- For every host used in a query, it starts pinging that host on a fixed
  interval (default 10s) and keeps the last 5000 results in memory.
- Each query returns the stored points within the dashboard time range as a
  time series frame (`time` + `rtt_ms`), so the standard **Time series** panel
  gives you exactly the diagram you want with no extra configuration.
- Failed/timeout pings are stored as `null` so gaps show up in the graph.

> The pings are sent from the machine running Grafana, and the Grafana host
> must have the `ping` (iputils/Windows) binary installed.

## Project layout

```
main.go                       plugin entrypoint (serves the backend)
pkg/plugin/datasource.go      ping scheduler, ring buffer, QueryData/CheckHealth
pkg/plugin/ping.go            system ping wrapper + output parsing
pkg/plugin/ping_test.go       parser unit tests
src/                          frontend (query editor, config editor, module.ts)
dist/                         build output (module.js, gpx_ping, plugin.json)
```

## Build

Backend (requires Go 1.22+):

```bash
go build -o dist/gpx_ping .
```

Frontend (requires Node 18+):

```bash
npm install
npm run build
```

`dist/` is the plugin folder that Grafana loads.

## Install

1. Copy this folder to the Grafana plugins directory, e.g.
   `/var/lib/grafana/plugins/grafana-ping-datasource` (Linux), containing `dist/`.
2. Allow the unsigned plugin in `grafana.ini` (or via env
   `GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=grafana-ping-datasource`):

   ```ini
   [plugins]
   allow_loading_unsigned_plugins = grafana-ping-datasource
   ```

3. Restart Grafana.

## Use

1. **Connections → Data sources → Add data source → Ping**.
   Optionally set a *Default interval* (e.g. `10s`) and *Ping timeout*
   (e.g. `3s`), then click **Save & test** (it pings 127.0.0.1 to verify).
2. Create a dashboard, add a **Time series** panel, select the Ping data
   source, and enter a host (e.g. `8.8.8.8`) and optional interval
   (e.g. `30s`).
3. Set the panel's time range and refresh. The first data points appear after
   the plugin has pinged the host a few times.

## Notes

- One query = one host. Add multiple queries to compare several servers on the
  same panel.
- History is in-memory only: it resets when Grafana restarts and points older
  than `interval × 5000` are evicted.
