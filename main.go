package main

import (
	"context"
	"flag"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	admissioncontrol "github.com/elithrar/admission-control"
	// Required to stay with v1beta1 due to https://github.com/elithrar/admission-control/issues/20
	admission "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	log "github.com/go-kit/kit/log"
)

type conf struct {
	TLSCertPath string
	TLSKeyPath  string
	HTTPOnly    bool
	Port        string
	Host        string
}

func main() {
	ctx := context.Background()

	// Get config
	conf := &conf{}
	flag.StringVar(&conf.Port, "port", "80", "The port to listen on (HTTP).")

	flag.Parse()

	// Set up logging
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "loc", log.DefaultCaller)

	// Set up the routes & logging middleware.
	r := mux.NewRouter().StrictSlash(true)
	// Show all available routes
	msg := "Funny Failing Webhook always rejects pod CREATE requests"
	r.Handle("/", printAvailableRoutes(r, logger, msg)).Methods(http.MethodGet)
	// Default health-check endpoint
	r.HandleFunc("/healthz", healthCheckHandler).Methods(http.MethodGet)

	r.Handle("/funny-failing-webhook", &admissioncontrol.AdmissionHandler{
		AdmitFunc: funnyFailingWebhook(),
		Logger:    logger,
	}).Methods(http.MethodPost)

	// HTTP server
	timeout := time.Second * 15
	srv := &http.Server{
		Handler:           admissioncontrol.LoggingMiddleware(logger)(r),
		Addr:              ":" + conf.Port,
		IdleTimeout:       timeout,
		ReadTimeout:       timeout,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
	}

	admissionServer, err := admissioncontrol.NewServer(
		srv,
		log.With(logger, "component", "server"),
	)
	if err != nil {
		fatal(logger, err)
		return
	}

	if err := admissionServer.Run(ctx); err != nil {
		fatal(logger, err)
		return
	}
}

func fatal(logger log.Logger, err error) {
	logger.Log(
		"status", "fatal",
		"err", err,
	)

	os.Exit(1)
	return
}

// healthCheckHandler returns a HTTP 200, everytime.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// printAvailableRoutes prints all routes attached to the provided Router, and
// prepends a message to the response.
func printAvailableRoutes(router *mux.Router, logger log.Logger, msg string) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		var routes []string
		err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, err := route.GetPathTemplate()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Log("msg", "walkFunc failed", err, err.Error())
				return err
			}

			routes = append(routes, path)
			return nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log("msg", "walkFunc failed", err, err.Error())
			return
		}

		fmt.Fprintln(w, msg)
		fmt.Fprintln(w, "Available routes:")
		for _, path := range routes {
			fmt.Fprintln(w, path)
		}
	}

	return http.HandlerFunc(fn)
}

func funnyFailingWebhook() admissioncontrol.AdmitFunc {
	return func(admissionReview *admission.AdmissionReview) (*admission.AdmissionResponse, error) {
		resp := &admission.AdmissionResponse{
			Allowed: false,
			Result:  &metav1.Status{},
		}
		err := fmt.Errorf("Yeah, no, not happening")
		return resp, err
	}
}
