package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type customResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w customResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		writer := &customResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		statusCode := c.Writer.Status()

		var data interface{}
		_ = json.Unmarshal(writer.body.Bytes(), &data)

		if _, ok := data.(map[string]interface{})["status"]; ok {
			return
		}

		response := ResponseBody{
			Status:  statusCode,
			Message: http.StatusText(statusCode),
			Data:    data,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeaderNow()
		_ = json.NewEncoder(c.Writer).Encode(response)
	}
}
