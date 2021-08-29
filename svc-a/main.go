package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lehoangthienan/go-jaeger-trace/utils/tracing"
	"github.com/opentracing/opentracing-go"

	httpCustom "github.com/lehoangthienan/go-jaeger-trace/utils/http"
)

const thisServiceName = "service-a"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Fail to load .env %v \n", err)
	}

	tracer, closer := tracing.Init(thisServiceName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	svc2Port, ok := os.LookupEnv("SERVICE_B_PORT")
	if !ok {
		svc2Port = "3001"
	}
	svc2Host := fmt.Sprintf("localhost:%s", svc2Port)

	http.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
		span := tracing.StartSpanFromRequest(tracer, r)
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(context.Background(), span)
		response, err := DoRequestGet(ctx, svc2Host)
		if err != nil {
			log.Fatalf("Error occurred: %s", err)
		}
		w.Write([]byte(fmt.Sprintf("%s -> %s", thisServiceName, response)))
	})

	svc1Port, okay := os.LookupEnv("PORT")
	if !okay {
		svc1Port = "3000"
	}

	log.Printf("Listening on localhost:%s", svc1Port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", svc1Port), nil))
}

func DoRequestGet(ctx context.Context, hostPort string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ping-send")
	defer span.Finish()

	url := fmt.Sprintf("http://%s/go", hostPort)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if err := tracing.Inject(span, req); err != nil {
		return "", err
	}
	return httpCustom.Do(req)
}
