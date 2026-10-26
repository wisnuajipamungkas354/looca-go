package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"message": msg})
}

func ParseID(c *gin.Context) (uint, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return uint(id), true
}
