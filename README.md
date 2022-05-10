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

## Content

**[1. Usage](#heading--1)**

  * [1.1. Run Tests and Start Server](#heading--1-1)
  * [1.2. Error Response](#heading--1-2)
  * [1.3. User Signup and Login](#heading--1-3)
    * [1.3.1. API -  POST /v1/Signup](#heading--1-3-1)
    * [1.3.2. API - POST /v1/Login](#heading--1-3-2)
  * [1.4. Loan APIs](#heading--1-4)
    * [1.4.1. API - GET /v1/Loans/{id}](#heading--1-4-1)
    * [1.4.2. API - GET /v1/Loans](#heading--1-4-2)
    * [1.4.3. API - POST /v1/Loans](#heading--1-4-3)
    * [1.4.4. API - PATCH /v1/Loans/{id}/Approve](#heading--1-4-4)
  * [1.5. Installment APIs](#heading--1-5)
    * [1.5.1. API - GET /v1/Installments/{id}](#heading--1-5-1)
    * [1.5.2. API - GET /v1/Installments](#heading--1-5-2)
    * [1.5.3. API - POST /v1/Installments](#heading--1-5-3)
    * [1.5.4. API - PATCH /v1/Installments/{id}/Repay](#heading--1-5-4)

**[2. Crons](#heading--2)**

  * [2.1. Insert Installments](#heading--2-1)
  * [2.2. Update Paid Loans](#heading--2-2) 


<div id="heading--1"/>

## Usage

<div id="heading--1-1"/>

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

<div id="heading--1-2"/>

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

<div id="heading--1-3"/>

### User Signup and Login

<div id="heading--1-3-1"/>

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

<div id="heading--1-3-2"/>

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

<div id="heading--1-4"/>

### Loan APIs

<div id="heading--1-4-1"/>

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

<div id="heading--1-4-2"/>

#### API - GET /v1/Loans

Request parameters:

###### URL Parameters:

**_userId_** : filter on User ID

**_state_**: filter on Loan state

**_sort_**: 

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;_format_: _field_name.order_

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;_example_: createdOn.desc

**_limit_**: limit number of results per page

**_page_**: page number (starts from 1)


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

<div id="heading--1-4-3"/>

#### API - POST /v1/Loans

<div id="heading--1-4-4"/>

#### API - PATCH /v1/Loans/{id}/Approve

<div id="heading--1-5"/>

### Installment APIs

<div id="heading--1-5-1"/>

#### API - GET /v1/Installments/{id}

<div id="heading--1-5-2"/>

#### API - GET /v1/Installments

<div id="heading--1-5-3"/>

#### POST /v1/Installments

<div id="heading--1-5-4"/>

#### PATCH /v1/Installments/{id}/Repay

<div id="heading--2"/>

## Crons

<div id="heading--2-1"/>

### Insert Installments

<div id="heading--2-2"/>

### Update Paid Loans














