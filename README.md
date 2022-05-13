# [Aspire Mini API](https://github.com/avinragh/aspire)

API in Go to add and view loans and payments

## Author
Avinash Raghunathan

## Features

* Implements Create, Fetch, Patch and Delete RESTful APIs
* Tests that run from `docker-compose up`
* Simple and Concise
* JWT Authentication for users
* Cron jobs for calculating 
* Filter, Sort and Pagination Support
* API DOCS on root page.

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
  * [1.6. Swagger API documentation](#heading--1.6)
  

**[2. Crons](#heading--2)**

  * [2.1. Insert Installments](#heading--2-1)
  * [2.2. Update Paid Loans](#heading--2-2) 

**[3. Run Environment and Dependencies](#heading--3)**
 
  * [3.1. Database](#heading--3-1)
  * [3.2. Models](#heading--3-2)
  * [3.3. Crons](#heading--3-3)
  * [3.4. Unit Tests](#heading--3-4)

**[4. Design](#heading--4)**
  * [4.1 Database](#heading--4-1)
  * [4.2 Approval API and Repay API](#heading--4-2)
  * [4.3 Installment calculation and Loan Paid Updation as crons](#heading--4-3)
  * [4.4 User Roles](#heading--4-4)


  
<div id="heading--1"/>

## Usage

<div id="heading--1-1"/>

### Run Tests and Start Server
<p>
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
</p>
<br/>

<div id="heading--1-2"/>

### Error Response

<p>All API Errors are returned in a single format and with the appropriate HTTP Status.
The format for the ErrorResponse is:

```json
{
    "code": "<string| Standard error code from the server>",
    "message": "<string| Standard error message for the error code>",
    "detail": "<string| Specific detail on the error if any else empty string" 
}
```
</p>
<br/>

<div id="heading--1-3"/>

### User Signup and Login

<div id="heading--1-3-1"/>

#### API -  POST /v1/Signup

<p>Request:

```json
{
    "username":"<username>", 
    "password":"<password", //required
    "role": "<role: admin/user>", //required
    "email" : "<email>" //required
}
```

Response:

Status: 200 OK
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

- ##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If the request from user is not well formed. 

- ##### 500 Internal Server Error:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Error when there is an error in the system

- ##### 409 Conflict:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If a user with given email already exists

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

- ##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If the email or password are not given or invalid

- ##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Error when there is an error in the system

<br/>

All consequent requests need to be sent with the Header **Token** and Value = _tokenString_
</p>
<br/>

<div id="heading--1-4"/>

### Loan APIs

<div id="heading--1-4-1"/>

#### API - GET /v1/Loans/{id}

<p>Response:

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

- ##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If the Id in the request is malformed

- ##### 403 Forbidden:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If Unauthorized

- #### 404 Not Found:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If resource Not found

- ##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Error when there is an error in the system

</p><br/>

<div id="heading--1-4-2"/>

#### API - GET /v1/Loans

<p>Request parameters:

###### URL Parameters:

- **_userId_** : filter on User ID

- **_state_**: filter on Loan state

- **_sort_**: 

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- _format_: _field_name.order_

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- _example_: createdOn.desc

- **_limit_**: limit number of results per page

- **_page_**: page number (starts from 1)


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

- ##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If the params in the request are malformed

- ##### 403 Forbidden:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If Unauthorized

- ##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Error when there is an error in the system

</p><br/>

<div id="heading--1-4-3"/>

#### API - POST /v1/Loans

<p>Request:

###### Request Body:
```json
{
    "amount": 100, //Decimal - Loan Amount
    "term": 4 //Integer- No of terms for the Loan
}
```

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

- ##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If the request is malformed

- ##### 403 Forbidden:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If Unauthorized

- ##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Error when there is an error in the system

</p><br/>

<div id="heading--1-4-4"/>

#### API - PATCH /v1/Loans/{id}/Approve

<p>User must be Admin

Response:
```json
{
    "amount": 100.0, //Decimal| Loan Amount
    "createdOn": "<timestamp>|2022-05-10T13:07:42.730Z",
    "currency": "<string|USD>",
    "id": 5, // Integer
    "installmentsCreated": false, //boolean | if installments have been created for the loan
    "modifiedOn": "<timestamp|2022-05-10T13:07:42.730Z>",
    "startDate": "<timestamp|0001-01-01T00:00:00.000Z>",
    "state": "<string|APPROVED>",
    "term": 4,//Integer | number of terms of the loan
    "userId": 8 //Integer
}
```

Error Responses:

- ##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If the id is malformed

- ##### 403 Forbidden:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If Unauthorized and if not admin

- #### 404 Not Found:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If resource Not found

- ##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Error when there is an error in the system

</p><br/>

<div id="heading--1-5"/>

### Installment APIs

<div id="heading--1-5-1"/>

#### API - GET /v1/Installments/{id}

<p>Response:
```json
{
    "createdOn": "<timestamp|2022-05-10T12:36:00.031Z",
    "dueDate": "<timestamp|2022-05-10T00:00:00.000Z",
    "id": 5, //Integer
    "installmentAmount": 50, //Decimal-The calculated installment amount as per schedule
    "loanId": 2, //Integer-The loan Id the installment belong to
    "modifiedOn": "timestamp|2022-05-10T12:36:00.031Z",
    "repaymentAmount": 0, //Integer- the actual amount paid when repaying
    "repaymentTime": "<timestamp|0001-01-01T00:00:00.000Z",
    "state": "<string|PENDING>"
}
```

ErrorResponses:

- ##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If the id is malformed

- ##### 403 Forbidden:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If Unauthorized and if not admin

- #### 404 Not Found:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If resource Not found

- ##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Error when there is an error in the system

</p><br/>

<div id="heading--1-5-2"/>

#### API - GET /v1/Installments

<p>Request parameters:

###### URL Parameters:

- **_loanId_** : filter on loan ID

- **_state_**: filter on installment state

- **_sort_**: 

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- _format_: _field_name.order_

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- _example_: createdOn.asc

- **_limit_**: limit number of results per page

- **_page_**: page number (starts from 1)

Response:

Status: 200 OK

```json
//List of Installments
[
    {
    "createdOn": "<timestamp|2022-05-10T12:36:00.031Z",
    "dueDate": "<timestamp|2022-05-10T00:00:00.000Z",
    "id": 5, //Integer
    "installmentAmount": 50, //Decimal-The calculated installment amount as per schedule
    "loanId": 2, //Integer-The loan Id the installment belong to
    "modifiedOn": "timestamp|2022-05-10T12:36:00.031Z",
    "repaymentAmount": 0, //Integer- the actual amount paid when repaying
    "repaymentTime": "<timestamp|0001-01-01T00:00:00.000Z",
    "state": "<string|PENDING>"
    },
    {...}
]
```

Error Responses:

- ##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If the params are malformed

- ##### 403 Forbidden:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If Unauthorized

- ##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Error when there is an error in the system

</p><br/>

<div id="heading--1-5-3"/>

#### POST /v1/Installments

<p>Request:

Only if Admin

###### Request Body:
```json
{
    "installmentAmount": 50, //Decimal-The calculated installment amount as per schedule
    "loanId": 2, //Integer-The loan Id the installment belong to  
}
```

Response:
```json
{
    "createdOn": "<timestamp|2022-05-10T12:36:00.031Z",
    "dueDate": "<timestamp|2022-05-10T00:00:00.000Z",
    "id": 5, //Integer
    "installmentAmount": 50, //Decimal-The calculated installment amount as per schedule
    "loanId": 2, //Integer-The loan Id the installment belong to
    "modifiedOn": "timestamp|2022-05-10T12:36:00.031Z",
    "repaymentAmount": 0, //Integer- the actual amount paid when repaying
    "repaymentTime": "<timestamp|0001-01-01T00:00:00.000Z",
    "state": "<string|PENDING>"
}
```

Error Responses:

- ##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If the request is malformed

- ##### 403 Forbidden:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If Unauthorized

- ##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Error when there is an error in the system

</p><br/>

<div id="heading--1-5-4"/>

#### PATCH /v1/Installments/{id}/Repay

<p>Repayment can be done only once. Subsequent request to the same installment will not cause any change.

Request:
```json
{
    "repaymentAmount": 25 //Decimal - Actual amount paid as repayment
}
```

Response:
```json
{
    
    "createdOn": "<timestamp|2022-05-10T12:36:00.031Z",
    "dueDate": "<timestamp|2022-05-10T00:00:00.000Z",
    "id": 5, //Integer
    "installmentAmount": 50, //Decimal-The calculated installment amount as per schedule
    "loanId": 2, //Integer-The loan Id the installment belong to
    "modifiedOn": "timestamp|2022-05-10T12:36:00.031Z",
    "repaymentAmount": 0, //Integer- set with repayment amount from request
    "repaymentTime": "<timestamp|0001-01-01T00:00:00.000Z",
    "state": "<string|PAID>"
}
```

Error Responses:

- ##### 400 Bad Request:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If the id is malformed

- ##### 403 Forbidden:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If Unauthorized and if not admin

- #### 404 Not Found:

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- If resource Not found

- ##### 500 Internal Server Error

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Error when there is an error in the system

</p><br/>

<div id="heading--1.6"/>

### Swagger API documentation

<p>The home page serves the API documentation generated with help of the swagger spec. Users can reference the same for help with accesing the APIs.</p>

<div id="heading--2"/>

## Crons

There are 2 crons that run every minute:

<div id="heading--2-1"/>

### Insert Installments:

<p>Calculates the repayment schedule and inserts it into the DB once loan is approved</p>

<div id="heading--2-2"/>

### Update Paid Loans:

<p>Checks Installments and Updates loans that are paid</p>

<div id="heading--3"/>

## Run Environment

<div id="heading--3-1"/>

### Database

<p>The app uses PostgreSQL for its database needs.

**Credentials**:

```ini
DB_USERNAME = "avinragh"
DB_PASSWORD = "toor"
DB_DATABASE = "aspire"
DB_PORT = "5432"
```

</p>
<div id="heading--3-2"/>

### Models

<p>The app uses swagger to generate the models that are used within the app</p>

<div id="heading--3-3"/>

### Crons

<p>The app uses the Golang package [robfig/cron/v3](https://github.com/robfig/cron) for its cron operations</p>

<div id="heading--3-4"/>

### Unit Tests

<p>The app uses the Golang package [stretchr/testify](https://github.com/stretchr/testify) for unit test assertions</p>

<div id="heading--4"/>

## Design

This section highlights a few design decisions and explains why they were taken:

<div id="heading--4-1"/>

### Database 

<p>The following 3 main models are used within the app:

- User
- Loan
- Installment

These are the relationships that exist:

- A loan belongs to a User
- An Installment (Repayment) belongs to a Loan

Since there are relations existing within the models and entities like loans are very transactional in nature, I decided to go with a relational database.
Hence, PostgreSQL was selected as the DB of choice.</p>

<div id="heading--4-2"/>

### Approval API and Repay API

<p>For the Approval and Repay APIs I decided to go with Patch APIs. By convention, PUT is a method of modifying resource where the client sends data that updates the entire resource. Both Approval and Repay APIs were such that only focused partial updates were required. Hence, it made more sense and the API more meaningful to have PATCH requests for the operations.</p>

<div id="heading--4-3"/>

### Installment calculation and Loan Paid Updation as crons

<p>The repayment schedule calculation is a user agnostic task which is done by the system and does not require APIs to be exposed to the end user. Also, the calculations need to be done only once a Loan is approved. For this reason, I decided to go with a cron task to check for Approved loans and add Installments(Repayments) on to the DB. Since, the app is cloud native, I have added a field installmentsCreated to the model so that this operation is done only once.

The updation of loans to Paid status is also a user agnostic task and need be done only when all the installments are paid. Hence this is also written out as a cron task and checks for the status of the installment.

In a full production system, where theres huge traffic, the trigger for these crons would have been notifications from the db operations of update. But, to maintain the simplicity of the app I did not go as far as to enable this. I rely on simple scheduled crons that run by checking relevant flags/fields within the DB. </p>

<div id="heading--4-4"/>

### User Roles

<p>The user can signup to the app with the following 2 roles
- admin
- user

The user has access to only their own loans and installments while the admin has access to the entire system.
The Approve API is an admin only API and cannot be performed by a user with user role access.</p>


















