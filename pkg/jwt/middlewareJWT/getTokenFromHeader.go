package middlewareJWT

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func ExtractHeaderToken(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return fmt.Sprint("no token found"), errors.New("no token provided")
	}
	tokenString := strings.TrimPrefix(header, "Bearer ")
	return tokenString, nil
}
