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

	"github.com/christianrm0821/Chirpy/internal/auth"
	"github.com/christianrm0821/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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

func mapChirpToValidChirp(myChirp database.Chirp) validChirp {
	valChirp := validChirp{
		ID:        myChirp.ID,
		CreatedAT: myChirp.CreatedAt,
		UpdatedAt: myChirp.UpdatedAt,
		Body:      myChirp.Body,
		UserID:    myChirp.UserID,
	}
	return valChirp
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
		PLATFORM:       os.Getenv("PLATFORM"),
		Secret:         os.Getenv("SECRET"),
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
		//make a decoder to get the request information
		decoder := json.NewDecoder(r.Body)
		request := email{}
		err := decoder.Decode(&request)
		if err != nil {
			errMsg := fmt.Sprintf("error decoding: %v", err)
			respondWithError(w, 500, errMsg)
			return
		}

		//hash password
		hashed_password, err := auth.HashPassword(request.Password)
		if err != nil {
			log.Fatal("error hashing the password: ", err)
			errMsg := fmt.Sprintf("error hashing password %v", err)
			respondWithError(w, 500, errMsg)
		}

		myEmailStruct := database.CreateUserParams{
			Email:          request.Email,
			HashedPassword: hashed_password,
		}

		user, err := counter.dbQueries.CreateUser(r.Context(), myEmailStruct)
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

	serveMux.HandleFunc("PUT /api/users", func(w http.ResponseWriter, r *http.Request) {
		//get token from header
		actualToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			errmsg := fmt.Sprintf("could not get token from header Error: %v", err)
			respondWithError(w, 401, errmsg)
			return
		}

		//get userID from token
		userID, err := auth.ValidateJWT(actualToken, counter.Secret)
		if err != nil {
			errmsg := fmt.Sprintf("error validating token Error: %v", err)
			respondWithError(w, 401, errmsg)
			return
		}

		//decode request
		decoder := json.NewDecoder(r.Body)
		request := email{}
		err = decoder.Decode(&request)
		if err != nil {
			errmsg := fmt.Sprintf("error decoding request Error: %v", err)
			respondWithError(w, 500, errmsg)
			return
		}

		myHashedPassword, err := auth.HashPassword(request.Password)
		if err != nil {
			errmsg := fmt.Sprintf("error hashing password Error: %v", err)
			respondWithError(w, 500, errmsg)
			return
		}
		err = auth.CheckPasswordHash(myHashedPassword, request.Password)
		if err != nil {
			errmsg := fmt.Sprintf("hashed password and password do not match Error: %v", err)
			respondWithError(w, 500, errmsg)
			return
		}

		newUserPasswordEmail := database.UpdatePasswordEmailFromUserIDParams{
			HashedPassword: myHashedPassword,
			Email:          request.Email,
			ID:             userID,
		}

		err = counter.dbQueries.UpdatePasswordEmailFromUserID(r.Context(), newUserPasswordEmail)
		if err != nil {
			errmsg := fmt.Sprintf("error updating password Error: %v", err)
			respondWithError(w, 500, errmsg)
			return
		}

		userInfo, err := counter.dbQueries.GetUserFromID(r.Context(), userID)
		if err != nil {
			errmsg := fmt.Sprintf("error getting user information Error: %v", err)
			respondWithError(w, 401, errmsg)
			return
		}

		respondWithJson(w, 200, userReturnEmail{
			ID:        userID,
			CreatedAT: userInfo.CreatedAt,
			UpdatedAt: userInfo.UpdatedAt,
			Email:     userInfo.Email,
		})

	})

	//register the login handler
	serveMux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		request := email{}
		err := decoder.Decode(&request)
		if err != nil {
			errMsg := fmt.Sprintf("error decoding: %v", err)
			respondWithError(w, 500, errMsg)
			return
		}

		user, err := counter.dbQueries.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}

		//checks if the password is correct
		err = auth.CheckPasswordHash(user.HashedPassword, request.Password)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}

		//This is getting a time of 1 hour which is the token life length
		expiredTimeDuration, err := time.ParseDuration("1h")
		if err != nil {
			respondWithError(w, 500, "could not convert time to duration")
			return
		}

		//makes a new token with current user ID, secret and expiration time
		token, err := auth.MakeJWT(user.ID, counter.Secret, expiredTimeDuration)
		if err != nil {
			respondWithError(w, 500, "could not make token")
			return
		}

		//make a new fresh token
		freshToken, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, 500, "could not produce a fresh token")
			return
		}

		//gets 60 days since that is the expire time of the refresh token
		expireTimeRefresh := time.Hour * 24 * 60

		//create a struct to input the refresh token into the database
		refreshTokenDataBase := database.CreateRefreshTokenParams{
			Token:     freshToken,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(expireTimeRefresh),
			RevokedAt: sql.NullTime{Valid: false},
		}

		_, err = counter.dbQueries.CreateRefreshToken(r.Context(), refreshTokenDataBase)
		if err != nil {
			errmsg := fmt.Sprintf("could not add refresh token to database: %v", err)
			respondWithError(w, 500, errmsg)
			return
		}

		respondWithJson(w, 200, userReturnEmail{
			ID:           user.ID,
			CreatedAT:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: freshToken,
		})

	})

	//gets a new token for the given user and sets the lifespan to 1 hour
	serveMux.HandleFunc("POST /api/refresh", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
		}
		user, err := counter.dbQueries.GetUserFromRefreshToken(r.Context(), token)
		if err != nil {
			respondWithError(w, 401, "token does not exist")
			return
		}
		if time.Now().After(user.ExpiresAt) {
			respondWithError(w, 401, "token has expired")
			return
		}

		if user.RevokedAt.Valid {
			respondWithError(w, 401, "token revoked")
			return
		}

		newToken, err := auth.MakeJWT(user.UserID, counter.Secret, time.Hour)
		if err != nil {
			respondWithError(w, 500, "could not make new token")
			return
		}
		respondWithJson(w, 200, tokenResponse{
			Token: newToken,
		})
	})

	//sets the revoke time to current time
	serveMux.HandleFunc("POST /api/revoke", func(w http.ResponseWriter, r *http.Request) {
		//gets token from header
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, 401, "could not get token from header")
			return
		}

		//get current user
		user, err := counter.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
		if err != nil {
			w.WriteHeader(204)
			return
		}
		if time.Now().After(user.ExpiresAt) {
			w.WriteHeader(204)
			return
		}

		if user.RevokedAt.Valid {
			w.WriteHeader(204)
			return
		}

		//gets current time to set in the database revoked time
		newTime := sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		updatedToken := database.RevokeRefreshTokenParams{
			RevokedAt: newTime,
			UpdatedAt: time.Now(),
			Token:     refreshToken,
		}

		//changes the revoke time, updated_at time for the given token
		err = counter.dbQueries.RevokeRefreshToken(r.Context(), updatedToken)
		if err != nil {
			respondWithError(w, 500, "error updating the database")
			return
		}

		//Sets header code to 204
		w.WriteHeader(204)
	})

	//register the validate_chirp handler
	//Makes sure chirp is valid
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

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			errmsg := fmt.Sprintf("%v", err)
			respondWithError(w, 401, errmsg)
			return
		}

		userID, err := auth.ValidateJWT(token, counter.Secret)
		if err != nil {
			//fmt.Printf("userID: %v      request.UserID: %v", userID.String(), request.UserID.String())
			respondWithError(w, 401, "Unauthorized")
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
			UserID: userID,
		}
		myChirp, err := counter.dbQueries.CreateChirp(r.Context(), input)
		if err != nil {
			errMsg := fmt.Sprintf("error creating chirp: %v", err)
			respondWithError(w, 500, errMsg)
			return
		}
		valChirp := mapChirpToValidChirp(myChirp)
		respondWithJson(w, 201, valChirp)
	})

	//register getting all chirps in database
	serveMux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		chirps, err := counter.dbQueries.GetAllChirps(r.Context())
		if err != nil {
			errmsg := fmt.Sprintf("error getting the chirps: %v", err)
			respondWithError(w, 500, errmsg)
			return
		}
		var valChirps []validChirp
		for _, val := range chirps {
			tmpChirp := mapChirpToValidChirp(val)
			valChirps = append(valChirps, tmpChirp)
		}
		respondWithJson(w, 200, valChirps)
	})

	//Gets a specific chirp given with the ID
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		chirpID := r.PathValue("chirpID")
		myChirp, err := counter.dbQueries.GetChirpWithID(r.Context(), uuid.MustParse(chirpID))
		if err != nil {
			errmsg := fmt.Sprintf("error getting this chirp: %v", err)
			respondWithError(w, 404, errmsg)
			return
		}
		respondWithJson(w, 200, mapChirpToValidChirp(myChirp))
	})

	//delete a specific chirp
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		userToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			errmsg := fmt.Sprintf("error getting token from header Error: %v", err)
			respondWithError(w, 401, errmsg)
			return
		}

		userIDToken, err := auth.ValidateJWT(userToken, counter.Secret)
		if err != nil {
			errmsg := fmt.Sprintf("could not validate user from token Error: %v", err)
			respondWithError(w, 500, errmsg)
			return
		}

		chirpID := r.PathValue("chirpID")
		myChirp, err := counter.dbQueries.GetChirpWithID(r.Context(), uuid.MustParse(chirpID))
		if err != nil {
			errmsg := fmt.Sprintf("error getting chirp with given ID Error: %v", err)
			respondWithError(w, 404, errmsg)
			return
		}

		if myChirp.UserID != userIDToken {
			respondWithError(w, 403, "Unauthorized")
			return
		}

		err = counter.dbQueries.DeleteChirpWithID(r.Context(), uuid.MustParse(chirpID))
		if err != nil {
			errmsg := fmt.Sprintf("could not delete chirp Error: %v", err)
			respondWithError(w, 500, errmsg)
			return
		}
		respondWithJson(w, 204, email{})
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
