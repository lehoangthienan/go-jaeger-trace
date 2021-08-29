package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lehoangthienan/go-jaeger-trace/utils/tracing"
	"github.com/opentracing/opentracing-go"
)

const thisServiceName = "service-b"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Fail to load .env %v \n", err)
	}

	tracer, closer := tracing.Init(thisServiceName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	http.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
		span := tracing.StartSpanFromRequest(tracer, r)
		defer span.Finish()

		w.Write([]byte(fmt.Sprintf("%s", thisServiceName)))
	})
	svc2Port, okay := os.LookupEnv("PORT")
	if !okay {
		svc2Port = "3001"
	}

	log.Printf("Listening on localhost:%s", svc2Port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", svc2Port), nil))
}
