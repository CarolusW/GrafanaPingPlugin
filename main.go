package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	"github.com/example/grafana-ping-datasource/pkg/plugin"
)

func main() {
	if err := datasource.Serve(serveOpts()); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

func serveOpts() datasource.ServeOpts {
	ds := plugin.NewDatasource()

	return datasource.ServeOpts{
		CheckHealthHandler: ds,
		QueryDataHandler:   ds,
	}
}
