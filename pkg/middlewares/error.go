package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // menjalankan handler dulu

		// Ambil semua error yang dikumpulkan Gin
		errs := c.Errors
		if len(errs) == 0 {
			return
		}

		// Ambil error terakhir
		err := errs.Last().Err
		if err == nil {
			return
		}

		// Log error
		log.Printf("❌ Gin Error: %v", err)

		// Jika response sudah terkirim, hentikan
		if c.Writer.Written() {
			return
		}

		// Format respons error
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
	}
}
