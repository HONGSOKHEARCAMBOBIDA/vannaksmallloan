package share

import "github.com/gin-gonic/gin"

// Exported function
func RespondError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}
func ResponeSuccess(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"success": message})
}
func RespondDate(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{"data": data})
}
