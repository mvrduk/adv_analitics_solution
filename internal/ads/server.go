package ads

import (
	"fmt"
	realip "github.com/ferluci/fast-realip"
	"github.com/mssola/user_agent"
	"github.com/oschwald/geoip2-golang"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net"
)

type Server struct {
	geoip  *geoip2.Reader
	logger *zap.Logger
}

func NewServer(geoip *geoip2.Reader) *Server {

	return &Server{geoip: geoip}
}

func (s *Server) Listen() error {
	return fasthttp.ListenAndServe(":8080", s.handleHttp)

}

func (s *Server) handleHttp(ctx *fasthttp.RequestCtx) {
	ip := realip.FromRequest(ctx)
	ua := string(ctx.Request.Header.UserAgent())
	parsed := user_agent.New(ua)

	browserName, browserVersion := parsed.Browser()

	country, err := s.geoip.Country(net.ParseIP(ip))
	if err != nil {
		s.logger.Fatal("GeoIP error", zap.Error(err))
		return
	}

	ctx.WriteString(fmt.Sprintf("User-Agent: %s\n", ua))
	ctx.WriteString(fmt.Sprintf("Browser: %s %s\n", browserName, browserVersion))
	ctx.WriteString(fmt.Sprintf("IP: %s\n", ip))
	ctx.WriteString("Analytics service started on 8080")
	ctx.WriteString(fmt.Sprintf("Country: %s", country.Country.IsoCode))
}
