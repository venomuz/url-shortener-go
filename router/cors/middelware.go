package cors

import (
	"github.com/rs/cors"
	"net/http"
)

func GinCorsMiddleware() Options {
	o := cors.Options{
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodOptions,
		},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders: []string{
			"Content-Type",
			"Content-Length",
			"accept",
			"Accept-Encoding",
			"origin",
			"Cache-Control",
			"X-Requested-With",
			"application/json",
		},
		OptionsPassthrough: false,
		ExposedHeaders: []string{
			"application/json",
			"Content-Type",
		},
		Debug:                true,
		OptionsSuccessStatus: 200,
	}

	return o
}
