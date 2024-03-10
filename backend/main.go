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

	"github.com/jackc/pgx/v4/pgxpool"

	db "github.com/nnaka2992/otel-database/backend/gen/sqlc"
)

// Create channel to listen for signals.
var (
	signalChan chan (os.Signal) = make(chan os.Signal, 1)
	pool       *pgxpool.Pool
	query      *db.Queries
)

func main() {
	db_password := os.Getenv("DB_PASSWORD")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")
	port := os.Getenv("PORT")

	connStr := "postgres://" + db_user + ":" + db_password + "@" + db_host + ":" + db_port + "/" + db_name + "?sslmode=disable"
	if err := initDB(connStr); err != nil {
		log.Fatal(err)
	}

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
		http.HandleFunc("/user/", getUserHandler)
		http.HandleFunc("/user/new", postUserAddHandler)
		http.HandleFunc("/user/delete", deleteUserDeleteHandler)
		http.HandleFunc("/user/update", postUserUpdateHandler)
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
