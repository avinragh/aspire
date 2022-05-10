# [Aspire Mini API](https://github.com/avinragh/aspire)

API in Go to add and view loans and payments

## Author
Avinash Raghunathan

## Features

* Implements Create, Fetch and Delete RESTful APIs
* Suitable for use in other software projects
* Tests that run from `docker-compose up`
* Simple and Concise
* JWT Authentication for users
* Cron jobs for calculating 

## Usage

### Run Tests and Start Server

```bash
docker-compose up
```
or Run in the Background

```bash
docker-compose up -d
```
```bash
docker-compose logs aspire-server
```

### Error Response

All API Errors are returned in a single format and with the appropriate HTTP Status.
The format for the ErrorResponse is:

```json
{
    "code": "<string| Standard error code from the server>",
    "message": "<string| Standard error message for the error code>",
    "detail": "<string| Specific detail on the error if any else empty string" 
}
```

### User Signup and Login

#### API - /v1/Signup

Request:

```json
{
    "username":"<username>", 
    "password":"<password", //required
    "role": "<role: admin/user>", //required
    "email" : "<email>" //required
}
```

Response:

```json
{
    "createdOn": "<timestamp |2022-05-10T13:03:45.354Z>",
    "email": "<string| user@gmail.com>",
    "id": 1, //integer
    "modifiedOn": "<timestamp|2022-05-10T13:03:45.354Z>",
    "password": "<timestamp| $2a$14$3HYiKBacLjvnxg3ZuKRA8eaQGC.SS7ZnuIw7fkspQkJeiEYWlVZy6 (encoded)",
    "role": "<string|(user/admin)>",
    "username": "<string|guestuser"
}
```

Error Responses:

##### 400 Bad Request:

If the request from user is not well formed. 

##### 500 Internal Server Error:

Error when theres an error in the system

##### 409 Conflict:

If a user with given email already exists


#### API - /v1/Login

Request:
```json
{
    "email" : "<email>", //required
    "password":"<password" //required
}
```

Response
```json
{
    "email": "<string |user@gmail.com>",
    "role": "<string|(user/admin)>",
    "tokenString": "string | token"
}
```
ErrorResponses:

##### 400 Bad Request:

If the email or password are not given or invalid

##### 500 Internal Server Error

Error when theres an error in the system

All consequent requests need to be sent with the Header Token and Value = tokenString







