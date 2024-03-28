// Sample application for webserver with opentelemetry-go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	db "github.com/nnaka2992/otel-database/backend/gen/sqlc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Create channel to listen for signals.
var (
	signalChan chan (os.Signal) = make(chan os.Signal, 1)
	query      *db.Queries
)

func main() {
	db_password := os.Getenv("DB_PASSWORD")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")
	port := os.Getenv("PORT")

	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background())

	connStr := "postgres://" + db_user + ":" + db_password + "@" + db_host + ":" + db_port + "/" + db_name + "?sslmode=disable"
	pool, err := db.NewDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	query = db.New(pool)

	ctx := context.Background()
	// SIGINT handles Ctrl+C locally.
	// SIGTERM handles Cloud Run termination signal.
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server.
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: nil,
	}
	go func() {
		http.Handle("/user/", otelhttp.NewHandler(http.HandlerFunc(getUserHandler), "getUserHandler"))
		http.Handle("/user/new", otelhttp.NewHandler(http.HandlerFunc(postUserAddHandler), "postUserAddHandler"))
		http.Handle("/user/new_many", otelhttp.NewHandler(http.HandlerFunc(postUserManyAddHandler), "postUserManyAddHandler"))
		http.Handle("/user/delete", otelhttp.NewHandler(http.HandlerFunc(deleteUserDeleteHandler), "deleteUserDeleteHandler"))
		http.Handle("/user/update", otelhttp.NewHandler(http.HandlerFunc(postUserUpdateHandler), "postUserUpdateHandler"))
		http.Handle("/health", otelhttp.NewHandler(http.HandlerFunc(healthHandler), "healthHandler"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Receive output from signalChan.
	sig := <-signalChan
	log.Printf("%s signal caught. Graceful Shutdown.", sig)

	// Gracefully shutdown the server by waiting on existing requests (except websockets).
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %+v", err)
	}
	log.Print("server exited")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
