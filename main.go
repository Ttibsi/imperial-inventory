package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Ttibsi/imperial-inventory/pkg"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
)

func main() {
	// Setting up the request router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Logging always helps too
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("| %-6s|", i)) }
	output.FormatMessage = func(i interface{}) string { return fmt.Sprintf("*%s*", i) }
	output.FormatFieldName = func(i interface{}) string { return fmt.Sprintf("%s:", i) }
	output.FormatFieldValue = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("%s", i)) }
	lgr := zerolog.New(output).With().Timestamp().Logger()

	// Open our persistant connection to the db
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/testdata")
	if err != nil {
		lgr.Fatal().Msg("Database not found")
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	conn := pkg.Conn{Db: db, Zlog: lgr}

	// A simple hello world to test the router with cURL
	// curl localhost:7333/
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			lgr.Warn().Msg("Error in root call: " + err.Error())
		}
	})

	r.Route("/ships", func(r chi.Router) {
		r.Get("/", conn.ListShips)             // GET /ships
		r.Get("/{shipID}", conn.SingleShip)    //GET /ships/123
		r.Delete("/{shipID}", conn.DeleteShip) //DELETE /ships/123

		// These two functions require a request body as well. See the comments
		// above the function definitions for details and examples
		r.Post("/new", conn.NewShip)        // POST /ships/new
		r.Put("/{shipID}", conn.UpdateShip) //PUT /ships/123
	})

	// 7333 is R3D3 on a T9 keypad
	fmt.Println("Serving on localhost:7333")
	lgr.Fatal().Msg(http.ListenAndServe(":7333", r).Error())
}
