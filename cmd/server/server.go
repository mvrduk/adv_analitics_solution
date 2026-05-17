package main

import (
	"adv_analitics_solution/internal/ads"
	"github.com/oschwald/geoip2-golang"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {

	geoip, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		logger.Fatal("Geo error", zap.Error(err))
	}

	s := ads.NewServer(geoip)
	if err := s.Listen(); err != nil {
		logger.Fatal("Something went wrong", zap.Error(err))
	}
}
