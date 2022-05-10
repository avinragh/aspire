package db

import (
	"aspire/constants"
	"aspire/models"
	"aspire/util"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
)

func (db *DB) FindLoanById(id int64) (*models.Loan, error) {
	loan := &models.Loan{}
	sqlFindById := `
		SELECT loans.id,loans.amount,loans.created_on,loans.term,loans.currency,loans.state,loans.modified_on,loans.start_date,loans.user_Id,loans.installments_created FROM loans WHERE loans.id=$1`

	err := db.QueryRow(sqlFindById, id).Scan(&loan.ID, &loan.Amount, &loan.CreatedOn, &loan.Term, &loan.Currency, &loan.State, &loan.ModifiedOn, &loan.StartDate, &loan.UserID, &loan.InstallmentsCreated)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

func (db *DB) FindLoans(userId int64, state *string, sort *string, limit *int64, page *int64) ([]*models.Loan, error) {
	loans := []*models.Loan{}
	if sort == nil {
		sortString := "createdOn.desc"
		sort = &sortString
	}
	var columnKey, operatorKey, columnName, operator string

	sortKeys := strings.Split(*sort, ".")
	if len(sortKeys) == 2 {
		columnKey = sortKeys[0]
		operatorKey = sortKeys[1]
	} else {
		return nil, errors.New("sort value not proper")
	}

	switch columnKey {
	case "createdOn":
		columnName = "created_on"
		break
	case "amount":
		columnName = "amount"
		break
	case "term":
		columnName = "term"
		break
	case "state":
		columnName = "state"
	case "startDate":
		columnName = "start_date"

	}

	switch operatorKey {
	case "asc":
		operator = "ASC"
	case "desc":
		operator = "DESC"
	}

	sortClause := fmt.Sprintf("ORDER BY %s %s", columnName, operator)

	sqlFind := `
		SELECT loans.id,loans.amount,loans.created_on,loans.term,loans.currency,loans.state,loans.modified_on,loans.start_date,loans.user_Id,loans.installments_created FROM loans`
	sqlFindByUser := `
		SELECT loans.id,loans.amount,loans.created_on,loans.term,loans.currency,loans.state,loans.modified_on,loans.start_date,loans.user_Id,loans.installments_created FROM loans where loans.user_id=$1`
	sqlFindByStatus := `
		SELECT loans.id,loans.amount,loans.created_on,loans.term,loans.currency,loans.state,loans.modified_on,loans.start_date,loans.user_Id,loans.installments_created FROM loans where loans.state=$1`
	sqlFindByUserAndStatus := `
		SELECT loans.id,loans.amount,loans.created_on,loans.term,loans.currency,loans.state,loans.modified_on,loans.start_date,loans.user_Id,loans.installments_created FROM loans where loans.user_id=$1 and loans.state=$2`
	var rows *sql.Rows
	var err error
	if userId != 0 {
		if state != nil {
			sqlFindByUserAndStatus = sqlFindByUserAndStatus + " " + sortClause
			if limit != nil && page != nil {
				sqlFindByUserAndStatus = sqlFindByUserAndStatus + " " + "LIMIT $3 " + "OFFSET $4"
				rows, err = db.Query(sqlFindByUserAndStatus, userId, state, limit, ((*page)-1)*(*limit))
				if err != nil {
					return nil, err
				}

			} else {
				rows, err = db.Query(sqlFindByUserAndStatus, userId, state)
				if err != nil {
					return nil, err
				}
			}

		} else {
			sqlFindByUser = sqlFindByUser + " " + sortClause
			if limit != nil && page != nil {
				sqlFindByUser = sqlFindByUser + " " + "LIMIT $2 " + "OFFSET $3"
				rows, err = db.Query(sqlFindByUser, userId, limit, ((*page)-1)*(*limit))
				if err != nil {
					return nil, err
				}

			} else {
				rows, err = db.Query(sqlFindByUser, userId)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		if state != nil {
			sqlFindByStatus = sqlFindByStatus + " " + sortClause
			if limit != nil && page != nil {
				sqlFindByStatus = sqlFindByStatus + " " + "LIMIT $2 " + "OFFSET $3"
				rows, err = db.Query(sqlFindByStatus, state, limit, ((*page)-1)*(*limit))
				if err != nil {
					return nil, err
				}

			} else {
				rows, err = db.Query(sqlFindByStatus, state)
				if err != nil {
					return nil, err
				}
			}

		} else {
			sqlFind = sqlFind + " " + sortClause
			if limit != nil && page != nil {
				sqlFind = sqlFind + " " + "LIMIT $1 " + "OFFFSET $2"
				rows, err = db.Query(sqlFind, limit, ((*page)-1)*(*limit))
				if err != nil {
					return nil, err
				}

			} else {
				rows, err = db.Query(sqlFind)
				if err != nil {
					return nil, err
				}
			}
		}

	}
	defer rows.Close()
	for rows.Next() {
		loan := &models.Loan{}
		if err := rows.Scan(&loan.ID, &loan.Amount, &loan.CreatedOn, &loan.Term, &loan.Currency, &loan.State, &loan.ModifiedOn, &loan.StartDate, &loan.UserID, &loan.InstallmentsCreated); err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, nil
}

func (db *DB) AddLoan(loan *models.Loan, userId int64) (*models.Loan, error) {
	if *loan.Amount == 0 || *loan.Term == 0 {
		return nil, errors.New("Loan and term Cannot be zero")
	}
	defaultTime := strfmt.DateTime(time.Time{}.UTC())
	currentDate := strfmt.DateTime(time.Now().UTC())
	loan.CreatedOn = currentDate
	loan.ModifiedOn = currentDate
	loan.StartDate = defaultTime
	loan.UserID = userId
	installmentCreated := false
	loan.InstallmentsCreated = &installmentCreated
	if loan.Currency == "" {
		loan.Currency = "USD"
	}
	loan.State = constants.LoanStatusPending

	sqlInsert := `
		INSERT INTO loans(amount,created_on,term,currency,state,modified_on,start_date,user_id,installments_created)
		VALUES ($1, $2, $3, $4,$5, $6, $7, $8, $9)
		RETURNING loans.id,loans.amount,loans.created_on,loans.term,loans.currency,loans.state,loans.modified_on,loans.start_date,loans.user_Id,loans.installments_created`
	err := db.QueryRow(sqlInsert, loan.Amount, loan.CreatedOn, loan.Term, loan.Currency, loan.State, loan.ModifiedOn, loan.StartDate, loan.UserID, loan.InstallmentsCreated).Scan(&loan.ID, &loan.Amount, &loan.CreatedOn, &loan.Term, &loan.Currency, &loan.State, &loan.ModifiedOn, &loan.StartDate, &loan.UserID, &loan.InstallmentsCreated)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

func (db *DB) ApproveLoan(loanId int64) (*models.Loan, error) {
	loan := &models.Loan{}
	sqlUpdate := `
		UPDATE loans
		SET state = $1, modified_on= $2, start_date=$3
		WHERE loans.id= $4
		RETURNING loans.id,loans.amount,loans.created_on,loans.term,loans.currency,loans.state,loans.modified_on,loans.start_date,loans.user_Id,loans.installments_created`

	currentTime := time.Now().UTC()
	err := db.QueryRow(sqlUpdate, constants.LoanStatusApproved, strfmt.DateTime(currentTime), strfmt.DateTime(util.GetDate(currentTime)), loanId).Scan(&loan.ID, &loan.Amount, &loan.CreatedOn, &loan.Term, &loan.Currency, &loan.State, &loan.ModifiedOn, &loan.StartDate, &loan.UserID, &loan.InstallmentsCreated)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

func (db *DB) DeleteLoan(id int64) error {
	sqlDelete := `
	DELETE FROM loans WHERE loans.id=$1`
	_, err := db.Exec(sqlDelete, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateToPaid(id int64) error {
	currentDate := time.Now().UTC()
	sqlUpdate := `
	UPDATE loans
	SET state= $1,modified_on =$2`
	_, err := db.Exec(sqlUpdate, constants.LoanStatusPaid, strfmt.DateTime(currentDate))
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) FindLoanIdsForUser(userId int64) ([]int64, error) {
	loanIds := []int64{}
	loans, err := db.FindLoans(userId, nil, nil, nil, nil)
	if err != nil {
		return loanIds, err
	}
	for _, loan := range loans {
		loanIds = append(loanIds, loan.ID)
	}
	return loanIds, nil
}
