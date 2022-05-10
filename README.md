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
* Filter, Sort and Pagination Support

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
<br/>

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
<br/>

### User Signup and Login

#### API -  POST /v1/Signup

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

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;If the request from user is not well formed. 

##### 500 Internal Server Error:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Error when there is an error in the system

##### 409 Conflict:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;If a user with given email already exists

<br/>

#### API - POST /v1/Login

Request:
```json
{
    "email" : "<email>", //required
    "password":"<password" //required
}
```

Response:

Status: 200 OK
```json
{
    "email": "<string |user@gmail.com>",
    "role": "<string|(user/admin)>",
    "tokenString": "string | token"
}
```
ErrorResponses:

##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;If the email or password are not given or invalid

##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Error when there is an error in the system

<br/>

All consequent requests need to be sent with the Header **Token** and Value = _tokenString_

<br/>

#### API - GET /v1/Loans/{id}

Response:

Status: 200 OK
```json
{
    "amount": 100.0, //Decimal| Loan Amount
    "createdOn": "<timestamp>|2022-05-10T13:07:42.730Z",
    "currency": "<string|USD>",
    "id": 5, // Integer
    "installmentsCreated": false, //boolean | if installments have been created for the loan
    "modifiedOn": "<timestamp|2022-05-10T13:07:42.730Z>",
    "startDate": "<timestamp|0001-01-01T00:00:00.000Z>",
    "state": "<string|PENDING>",
    "term": 4,//Integer | number of terms of the loan
    "userId": 8 //Integer
}
```

Error Responses:

##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;If the Id in the request is malformed

##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Error when there is an error in the system

<br/>



#### API - GET /v1/Loans

Request parameters:

###### URL Parameters:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;**_userId_** : filter on User ID

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;**_state_**: filter on Loan state

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;**_sort_**: 

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;_format_: field_name.order

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;_example_: createdOn.desc

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;**_limit_**: limit number of results per page

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;**_page_**: page number (starts from 1)

Response:

Status: 200 OK

```json
//List of Loan
[
    {
    "amount": 100.0, //Decimal| Loan Amount
    "createdOn": "<timestamp>|2022-05-10T13:07:42.730Z",
    "currency": "<string|USD>",
    "id": 5, // Integer
    "installmentsCreated": false, //boolean | if installments have been created for the loan
    "modifiedOn": "<timestamp|2022-05-10T13:07:42.730Z>",
    "startDate": "<timestamp|0001-01-01T00:00:00.000Z>",
    "state": "<string|PENDING>",
    "term": 4,//Integer | number of terms of the loan
    "userId": 8 //Integer        
    },
    {....}
]
```

Error Responses:

##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;If the params in the request are malformed

##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Error when there is an error in the system

<br/>







