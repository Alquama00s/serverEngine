package serverEngine

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	routing "github.com/Alquama00s/serverEngine/lib/routing"

	"github.com/Alquama00s/serverEngine/lib/logging/loggerFactory"
	"github.com/Alquama00s/serverEngine/lib/routing/registrar"
)

var (
	routeRegistrar *registrar.DefaultRegistrar
	once           sync.Once
	logger         = loggerFactory.GetLogger()
	serverPort     = "8080"
)

func Registrar() *registrar.DefaultRegistrar {
	once.Do(func() {
		logger.Info().Msg("Initializing route registrar")
		if routeRegistrar == nil {
			routeRegistrar = &registrar.DefaultRegistrar{
				Routers: make(map[string]*routing.Router),
				Handler: &routing.SimpleHandler{},
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

func Port(port int) {
	serverPort = strconv.Itoa(port)
}

func Sereve() {
	mx := http.NewServeMux()
	logger.Info().Msg("Starting server...")
	routeRegistrar.Initialize()
	routeRegistrar.FinalizeRoutes(mx)
	logger.Info().Msg("Server is running on port " + serverPort)
	lmx := LoggingMiddleware(mx)
	log.Fatal(http.ListenAndServe(":"+serverPort, lmx))
}
