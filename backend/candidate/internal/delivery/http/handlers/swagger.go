package handlers

import (
	"net/http"
	"strings"

	"github.com/HghaVlad/trainee-match/backend/candidate/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SwaggerHandler
// @Summary Swagger UI
// @Description This endpoint serves the Swagger UI for the API documentation.
// @Tags swagger
// @Accept json
// @Produce html
// @Success 200 {string} string "Swagger UI served successfully"
// @Router /swagger/ [get]
func SwaggerHandler(w http.ResponseWriter, r *http.Request) {
	docs.SwaggerInfo.Host = r.Host

	if r.TLS != nil || strings.HasPrefix(r.Proto, "HTTPS") {
		docs.SwaggerInfo.Schemes = []string{"https"}
	}
	if prefix := r.Header.Get("X-Forwarded-Prefix"); prefix != "" {
		docs.SwaggerInfo.BasePath = prefix + "/api/v1"
	} else {
		docs.SwaggerInfo.BasePath = "/api/v1"
	}
	httpSwagger.WrapHandler(w, r)
}
