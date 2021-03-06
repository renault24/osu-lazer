package customlogger

import (
	"github.com/deissh/osu-lazer/api/pkg/services/oauth"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			req := c.Request()
			res := c.Response()
			start := time.Now()

			id := req.Header.Get(echo.HeaderXRequestID)

			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			logger := log.Logger

			if err = next(c); err != nil {
				c.Error(err)
			}

			stop := time.Now()

			evt := logger.Info()

			// getting token information if present
			data := c.Get("current_user_token")
			if userToken, ok := data.(*oauth.Token); ok {
				evt.Uint("user_id", userToken.UserID)
				evt.Str("user_scopes", userToken.Scopes)
			}
			if id != "" {
				evt.Str("request_id", id)
			}
			evt.Str("remote_ip", c.RealIP())
			evt.Str("host", req.Host)
			evt.Str("method", req.Method)
			evt.Str("uri", req.RequestURI)
			evt.Str("user_agent", req.UserAgent())
			evt.Int("status", res.Status)
			evt.Str("referer", req.Referer())

			if err != nil {
				evt.Err(err)
			}

			evt.Dur("latency", stop.Sub(start))
			evt.Str("latency_human", stop.Sub(start).String())

			cl := req.Header.Get(echo.HeaderContentLength)
			if cl == "" {
				cl = "0"
			}

			evt.Str("bytes_in", cl)
			evt.Str("bytes_out", strconv.FormatInt(res.Size, 10))
			evt.Msg("")

			return err
		}
	}
}
