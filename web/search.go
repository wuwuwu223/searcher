package web

import (
	"github.com/gin-gonic/gin"
	"searcher/search/search"
)

// Search 搜索API
func Search(c *gin.Context) {
	str := c.Query("s")
	list := search.Search(str)
	c.JSON(200, list)
}
