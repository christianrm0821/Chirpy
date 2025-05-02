# API for Chirpy

## Users 

```json
{
    "id": "user_id in uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp", 
    "email": "example@email.com",
    "is_chirpy_red": "Bool value(true if premium subscription, false by default)",
    "hashed_password": "password"
}
```

### GET /admin/********

Shows how many times chirp has been visited

### POST /admin/******

Removes all the users that are registered

### "POST /api/users"

Creates a new user with the following email and password

Request Body: 

```json
{
    "email": "example@email.com",
    "hashed_password": "password"
}
```

### "PUT /api/users"

Updates the user email and password

Request Body: 

```json
{
    "email": "example@email.com",
    "hashed_password": "password"
}
```

### "POST /api/login"

Logs into the user account with the given email and password
Also makes 2 tokens one that lasts for 1 hour and the other that lasts for 60 days

Request Body: 

```json
{
    "email": "example@email.com",
    "hashed_password": "password"
}
```

### "POST /api/refresh"

makes a new token that expires in an hour
No request body required

### "POST /api/revoke"

Revokes the current token(sets the time to current time and true)
No request body required

### "POST /api/chirps"

Allows user to post a chirp
If the chirp is invalid(length too long) it flags it
If it uses key words that cannot be used it replaces them with "****"

Request Body: 

```json
{
    "body": "text that you wish to be posted"
}
```

### "GET /api/chirps"

Has query parameters author_id and sort

GET /api/chirps: returns every chirp ordered in ascending order from the time they were created at
GET /api/chirps?author_id=(id for given user): returns every chirp from that specific user ordered in ascending order from the time they were created
GET /api/chirps?sorted=desc: returns every chirp ordered in descending order from the time that they were created
GET /api/chirps?author_id=(author id)&sort=desc: returns all of the chirps from the specific user ordered in descending order from the time they were created

No Request Body required

### "GET /api/chirps/{chirpID}"

Gets the chirp from the chirpID provided

No request body required

### "DELETE /api/chirps/{chirpID}"

Deletes the chirp with the ChirpID provided if user is authorized

No Request Body required

###  "POST /api/polka/webhooks"

Request Body:
```json
{
	"event":"user.upgraded",
	"data": {
		"user_id": "user to be upgraded"
	} "data stuct"
}

```



