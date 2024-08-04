package middleware

import (
	"net"
	"net/http"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/gin-gonic/gin"
)

func ValidateIP(trustedSubnetText string) (gin.HandlerFunc, error) {
	_, trustedSubnet, err := net.ParseCIDR(trustedSubnetText)

	if err != nil {
		return nil, err
	}

	return func(c *gin.Context) {
		// Получаем IP из заголовка
		IPText := c.Request.Header.Get(common.IPHeader)
		
		IP := net.ParseIP(IPText)

		// Проверяем входит ли IP в доверенные
		if !trustedSubnet.Contains(IP) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Пропускаем дальше
		c.Next()
	}, nil
}
