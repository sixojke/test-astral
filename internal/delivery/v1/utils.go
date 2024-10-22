package v1

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func (h *Handler) filePathGenerator(userId, fileName string) string {
	return fmt.Sprintf("%v/%v/%v", h.config.Documents.UploadsDir, userId, fileName)
}

func getUserIdByContext(c *gin.Context) string {
	return c.MustGet("userId").(string)
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
