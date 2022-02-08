package main

import (
	"code.cloudfoundry.org/lager"
	"lbmaster-advanced-groups-api/handler"
	"lbmaster-advanced-groups-api/internal"
	"lbmaster-advanced-groups-api/internal/adapter"
	"net/http"
	"os"
	"strconv"
)

func main() {
	logger := lager.NewLogger("lbmaster-advanced-groups-api")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))

	c, err := internal.NewConfig("./config.json", logger)
	if err != nil {
		logger.Fatal("config", err)
	}

	r := adapter.NewJsonPrefixGroupRepository(c.AdvancedGroupsConfigPath)
	h := handler.NewHandler(r, c.ApiKey, logger)

	logger.Info("start-listener", lager.Data{"port": c.Port})
	err = http.ListenAndServe(":"+strconv.Itoa(c.Port), h)
	if err != nil {
		logger.Fatal("start-listener", err)
	}
}
