package ads

import (
	"adv_analitics_solution/internal/stats"
	"fmt"
	realip "github.com/ferluci/fast-realip"
	"github.com/mssola/user_agent"
	"github.com/oschwald/geoip2-golang"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net"
	"net/http"
)

type Server struct {
	geoip  *geoip2.Reader
	logger *zap.Logger
	stats  *stats.Manager
}

func NewServer(geoip *geoip2.Reader, stats *stats.Manager) *Server {

	return &Server{geoip: geoip, stats: stats}
}

func (s *Server) Listen() error {
	return fasthttp.ListenAndServe(":8080", s.handler)

}

func (s *Server) handler(ctx *fasthttp.RequestCtx) {

	remoteIp := realip.FromRequest(ctx)
	ua := string(ctx.Request.Header.UserAgent())
	parsed := user_agent.New(ua)
	browserName, browserVersion := parsed.Browser()

	statsKey := stats.NewKey(stats.Key{
		Os:      parsed.OS(),
		Browser: browserName,
	})

	statsValue := stats.Value{Requests: 1}

	defer func() {
		s.stats.Append(
			statsKey, statsValue)
	}()

	country, err := s.geoip.Country(net.ParseIP(remoteIp))
	if err != nil {
		s.logger.Fatal("GeoIP error", zap.Error(err))
		return
	}

	statsKey.Country = country.Country.IsoCode

	u := &User{Country: country.Country.IsoCode, Browser: browserName}

	campaigns := GetCampaigns()

	winner := MakeAuction(campaigns, u)

	if winner == nil {
		ctx.Redirect("https//example.com", http.StatusSeeOther)
		return
	}

	statsKey.CampaignId = winner.Id
	statsValue.Impressions++

	ctx.Redirect(winner.ClickUrl, http.StatusSeeOther)
	ctx.WriteString(fmt.Sprintf("User-Agent: %s\n", ua))
	ctx.WriteString(fmt.Sprintf("Browser: %s %s\n", browserName, browserVersion))
	ctx.WriteString(fmt.Sprintf("IP: %s\n", remoteIp))
	ctx.WriteString("Analytics service started on 8080")
	ctx.WriteString(fmt.Sprintf("Country: %s", country.Country.IsoCode))
}
