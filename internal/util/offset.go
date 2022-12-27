package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetOffsetParam(c *gin.Context) (int, error) {
	query := c.Param("offset")
	offset := 0
	if value := query; value != "" {
		offset, err := strconv.Atoi(value)
		if err != nil || offset < 0 || offset > 100 {
			return 0, errors.New("invalid param")
		}
	}
	return offset, nil
}
