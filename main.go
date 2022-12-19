package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/opentelemetry-go-contrib/launcher"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("godotenv.Load", err)
	}

	shutdown, err := launcher.ConfigureOpenTelemetry(
		honeycomb.WithApiKey(os.Getenv("HONEYKOMB_API_KEY")),
		launcher.WithServiceName(os.Getenv("OTEL_SERVICE_NAME")),
	)
	if err != nil {
		log.Fatal("launcher.ConfigureOpenTelemetry", err)
	}
	defer shutdown()

	server := http.NewServeMux()

	server.Handle(
		"/my-endpoint",
		otelhttp.NewHandler(
			otelhttp.WithRouteTag("/my-endpoint", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("Received request")

				if _, err := w.Write([]byte("Hello World!\n")); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusAccepted)
			})),
			"my-operation",
			otelhttp.WithPublicEndpoint(),
		),
	)

	if err := http.ListenAndServe(":3030", server); err != nil {
		log.Fatal("http.ListenAndServe", err)
	}
}
