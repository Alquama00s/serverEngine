package serverEngine

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Alquama00s/serverEngine/lib"
	"github.com/Alquama00s/serverEngine/loggerFactory"
	"github.com/Alquama00s/serverEngine/registrar"
)

var (
	routeRegistrar *registrar.DefaultRegistrar
	once           sync.Once
	logger         = loggerFactory.GetLogger()
)

func Registrar() *registrar.DefaultRegistrar {
	once.Do(func() {
		logger.Info().Msg("Initializing route registrar")
		if routeRegistrar == nil {
			routeRegistrar = &registrar.DefaultRegistrar{
				Routers: make(map[string]*lib.Router),
				Handler: &lib.SimpleHandler{},
			}
		}
	})
	return routeRegistrar
}

// LoggingMiddleware logs method, path, headers, and duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log method and URL
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Log headers
		for name, values := range r.Header {
			for _, value := range values {
				log.Printf("Header: %s: %s", name, value)
			}
		}

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log duration
		duration := time.Since(start)
		log.Printf("Completed in %v", duration)
	})
}

func Sereve() {
	logger.Info().Msg("Starting server...")
	routeRegistrar.Initialize()
	routeRegistrar.FinalizeRoutes()
	logger.Info().Msg("Server is running on port 8080")
	mx := http.NewServeMux()
	lmx := LoggingMiddleware(mx)
	log.Fatal(http.ListenAndServe(":8080", lmx))
}
