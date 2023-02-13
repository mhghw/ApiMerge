package ports

import (
	"mergeApi/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

func (s HttpServer) HandleGetRates(c *gin.Context) {
	var body RequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": "bad body for request",
		})
		return
	}

	for _, cur := range body.Currencies {
		if !funk.ContainsString(config.AppConfig.Currencies, strings.ToUpper(cur)) &&
			strings.ToUpper(cur) != config.AppConfig.BaseCurrency {
			c.JSON(http.StatusBadRequest, map[string]any{
				"error": "currencies are not supported",
			})
			return
		}
	}
	if !funk.ContainsString(config.AppConfig.Currencies, strings.ToUpper(body.To)) {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": "to currency is not supported",
		})
	}

	rates, err := s.app.Queries.GetRates.Handle(c.Request.Context(), body.Currencies, body.To)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resp := make([]Response, 0)
	for _, rate := range rates {
		resp = append(resp, Response(rate))
	}

	c.JSON(http.StatusOK, resp)
}
