package apiCmd

import (
	"context"
	"crypto/tls"
	"dgraph-client/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgraph-io/dgo/v2"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start api server",
	Long: `starts the REST api server to allow remote connections to the specified
			address.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start API Server
		if err := startServer(apiCfg); err != nil {
			log.Fatalln("fatal error starting api server -", err)
		}
	},
}

func init() {
	Cmd.AddCommand(startCmd)
}

type API struct {
	DGraph *dgo.Dgraph
}

func (a *API) routes() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/query", a.query)

	return mux
}

func (a *API) home(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status string `json:"status"`
		URL    string `json:"urls"`
	}{
		Status: "active",
		URL:    "/query",
	}

	writeJson(w, data)
}

func (a *API) query(w http.ResponseWriter, r *http.Request) {
	var data interface{}
	writeJson(w, data)
}

func startServer(cfg *config.APIConfig) error {
	addr := fmt.Sprint(cfg.ApiAddr)
	log.Println("starting API server -", addr)
	defer log.Println("gracefully shutting down API server")

	var a API
	/*
		// Start dgraph client
		dgraph, cnclFunc := data.NewDGClient(cfg)
		defer cnclFunc()
	*/

	// TODO - add TLS support for production
	// tlscert := fmt.Sprintf("%s/app.crt", cfg.CertsDir)
	// tlskey := fmt.Sprintf("%s/app.key", cfg.CertsDir)
	aServe := http.Server{
		Addr:    addr,
		Handler: a.routes(),
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ReadTimeout:  cfg.ApiReadTimeout,
		WriteTimeout: cfg.ApiWriteTimeout,
		IdleTimeout:  cfg.ApiIdleTimeout,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		// TODO - add TLS support for production
		// serverErrors <- api.ListenAndServeTLS(tlscert, tlskey)
		serverErrors <- aServe.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("api server error - %w", err)

	case sig := <-shutdown:
		log.Printf("Graceful shutdown started - %s\n", sig)
		defer log.Printf("Graceful shutdown completed - %s\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()

		if err := aServe.Shutdown(ctx); err != nil {
			aServe.Close()
			return fmt.Errorf("fuck man...Server did not shutdown gracefully: %w", err)
		}
	}
	return nil
}

func writeJson(w http.ResponseWriter, data interface{}) {
	out, _ := json.MarshalIndent(data, "", " ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}
