package constants

//Loan Status Constants
const (
	LoanStatusApproved = "APPROVED"
	LoanStatusRejected = "REJECTED"
	LoanStatusPending  = "PENDING"
	LoanStatusPaid     = "PAID"
)

//Installment Status
const (
	InstallmentStatusPending = "PENDING"
	InstallmentStatusPaid    = "PAID"
)

//User Role Constants
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

//Database Credentials
const (
	DbUsername = "avinragh"
	DbPassword = "toor"
	DbDatabase = "aspire"
	DbHost     = "host.docker.internal" //"localhost" //
	DbPort     = "5432"
)

//Environment Variables
const (
	EnvDbUsername = "DB_USERNAME"
	EnvDbPassword = "DB_PASSWORD"
	EnvDbDatabase = "DB_DATABASE"
	EnvDbHost     = "DB_HOST"
	EnvDbPort     = "DB_PORT"
)
