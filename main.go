package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/christianrm0821/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// to count number of times site is visited(hits)
type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	PLATFORM       string
}

// response types
/*
type req struct {
	Body string `json:"body"`
}

type resClean struct {
	CleanedBody string `json:"cleaned_body"`
}
*/

type resErr struct {
	Error string `json:"error"`
}

type email struct {
	Email string `json:"email"`
}

type userReturnEmail struct {
	ID        uuid.UUID `json:"id"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type chirpPostReq struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type validChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	generalErr := resErr{
		Error: msg,
	}
	res, _ := json.Marshal(generalErr)
	w.Write(res)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	res, _ := json.Marshal(payload)
	w.Write(res)
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error: ", err)
		return
	}

	//making a newserveMux
	const port = ":8080"

	//keeps count of how many requests are being made
	counter := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
		PLATFORM:       "dev",
	}
	counter.fileserverHits.Store(0)

	//mux or multiplexer
	//it is a request router
	// it gets incoming http requests and decides which handler function should process the request
	//maps url patterns to handler functions
	serveMux := http.NewServeMux()

	//handlefunc register handlers with serveMux
	//takes in the "/healthz" endpoint
	//takes in a function with the signature "func(http.ResponseWriter, *http.Request)"
	//It automatically converts your function to a handler interface
	serveMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	//registers handler with serveMux
	//takes url path with an object with the method "ServeHTTP(http.ResponseWriter, *http.Request)"
	//used for pre-built handlers or custom handler type
	//want to use this over a handle in more complex situations such as with the fileserver handler or using miiddleware like stripPrefix
	//Strip prefix takes away the prefix "/app" from the handler
	//FileServer is a built in handler, automatically handles file serving, content types, and directory listings
	//FileServer serves static content
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", counter.MiddlewareMetricsInc(appHandler))

	//register the metrics handler
	serveMux.HandleFunc("GET /admin/metrics", counter.RequestNum)

	//register the reset handler
	serveMux.HandleFunc("POST /admin/reset", counter.resetComplete)

	//register the users handler
	serveMux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		request := email{}
		err := decoder.Decode(&request)
		if err != nil {
			errMsg := fmt.Sprintf("error decoding: %v", err)
			respondWithError(w, 500, errMsg)
			return
		}
		user, err := counter.dbQueries.CreateUser(r.Context(), request.Email)
		if err != nil {
			errMsg := fmt.Sprintf("error creating user: %v", err)
			respondWithError(w, 500, errMsg)
			return
		}
		respondWithJson(w, 201, userReturnEmail{
			ID:        user.ID,
			CreatedAT: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		})
	})

	//register the validate_chirp handler
	serveMux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		//get the information and putting it into request
		decoder := json.NewDecoder(r.Body)
		request := chirpPostReq{}
		err := decoder.Decode(&request)
		//handling general error with decoding
		if err != nil {
			errMsg := fmt.Sprintf("error decoding: %v", err)
			respondWithError(w, 500, errMsg)
			return
		}
		//handling if the length of the request body(the message) is too long
		if len(request.Body) > 140 {
			respondWithError(w, 400, "Chirp is too long")
			return
		}

		//handling if the request was successful
		fmt.Printf("request body: %v", request.Body)
		cleanText := ValidString(request.Body)
		fmt.Println(cleanText)

		input := database.CreateChirpParams{
			Body:   cleanText,
			UserID: request.UserID,
		}
		myChirp, err := counter.dbQueries.CreateChirp(r.Context(), input)
		if err != nil {
			errMsg := fmt.Sprintf("error creating chirp: %v", err)
			respondWithError(w, 500, errMsg)
			return
		}
		respondWithJson(w, 201, validChirp{
			ID:        myChirp.ID,
			CreatedAT: myChirp.CreatedAt,
			UpdatedAt: myChirp.UpdatedAt,
			Body:      myChirp.Body,
			UserID:    myChirp.UserID,
		})
	})

	//making the server struct
	myServer := &http.Server{
		Addr:    port,
		Handler: serveMux,
	}

	//start an http server with the port and handler we created above/ handles any errors
	log.Println("Starting server on port #", port)
	err = myServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal("server error: ", err)
	}
}
