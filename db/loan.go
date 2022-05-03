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
		SELECT id,amount,created_on,term,currency,state,modified_on,start_date,user_Id FROM loans WHERE id=$1`

	err := db.QueryRow(sqlFindById, id).Scan(&loan.ID, &loan.Amount, &loan.CreatedOn, &loan.Term, &loan.Currency, &loan.State, &loan.ModifiedOn, &loan.StartDate, &loan.UserID)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

func (db *DB) FindLoans(userId int64, state string, sort string, limit int64, page int64) ([]*models.Loan, error) {
	loans := []*models.Loan{}
	if sort == "" {
		sort = "createdOn.desc"
	}
	var columnKey, operatorKey, columnName, operator string

	sortKeys := strings.Split(sort, ".")
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
		SELECT id,amount,created_on,term,currency,state,modified_on,start_date,user_Id FROM loans`
	sqlFindByUser := `
		SELECT id,amount,created_on,term,currency,state,modified_on,start_date,user_Id FROM loans where user_id=$1`
	sqlFindByStatus := `
		SELECT id,amount,created_on,term,currency,state,modified_on,start_date,user_Id FROM loans where state=$1`
	sqlFindByUserAndStatus := `
		SELECT id,amount,created_on,term,currency,state,modified_on,start_date,user_Id FROM loans where user_id=$1 and state=$2`
	var rows *sql.Rows
	var err error
	if userId != 0 {
		if state != "" {
			sqlFindByUserAndStatus = sqlFindByUserAndStatus + " " + sortClause
			if limit != 0 && page != 0 {
				sqlFindByUserAndStatus = sqlFindByUserAndStatus + " " + "LIMIT $3 " + "OFFSET $4"
				rows, err = db.Query(sqlFindByUserAndStatus, userId, state, limit, page*limit)
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
			if limit != 0 && page != 0 {
				sqlFindByUser = sqlFindByUser + " " + "LIMIT $2 " + "OFFSET $3"
				rows, err = db.Query(sqlFindByUser, userId, limit, page*limit)
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
		if state != "" {
			sqlFindByStatus = sqlFindByStatus + " " + sortClause
			if limit != 0 && page != 0 {
				sqlFindByStatus = sqlFindByStatus + " " + "LIMIT $2 " + "OFFSET $3"
				rows, err = db.Query(sqlFindByStatus, state, limit, page*limit)
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
			if limit != 0 && page != 0 {
				sqlFind = sqlFind + " " + "LIMIT $1 " + "OFFFSET $2"
				rows, err = db.Query(sqlFind, limit, page*limit)
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
		if err := rows.Scan(loan.ID, loan.Amount, loan.CreatedOn, loan.Term, loan.Currency, loan.State, loan.ModifiedOn, loan.StartDate, loan.UserID); err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, nil
}

func (db *DB) AddLoan(loan *models.Loan, userId int64) (*models.Loan, error) {
	defaultTime := strfmt.DateTime(time.Time{}.UTC())
	currentDate := strfmt.DateTime(time.Now().UTC())
	loan.CreatedOn = currentDate
	loan.ModifiedOn = currentDate
	loan.StartDate = defaultTime
	loan.UserID = userId
	sqlInsert := `
		INSERT INTO loans(amount,created_on,term,currency,state,modified_on,start_date,user_id)
		VALUES ($1, $2, $3, $4,$5, $6, $7, $8)
		RETURNING id`
	var id int64
	err := db.QueryRow(sqlInsert, loan.Amount, loan.CreatedOn, loan.Term, loan.Currency, loan.State, loan.ModifiedOn, loan.StartDate, loan.UserID).Scan(&id)
	if err != nil {
		return nil, err
	}
	loan.ID = &id
	return loan, nil
}

func (db *DB) ApproveLoan(loanId int64) error {
	sqlUpdate := `
		UPDATE loans
		SET state = $1, modifiedOn = $2, startDate = $3
		WHERE id= $4;`

	currentTime := time.Now().UTC()
	_, err := db.Exec(sqlUpdate, constants.LoanStatusApproved, strfmt.DateTime(currentTime), strfmt.DateTime(util.GetDate(currentTime)), loanId)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteLoan(id int64) error {
	sqlDelete := `
	DELETE FROM loans WHERE id=$1;`
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
	SET state= $1,modifiedOn =$2`
	_, err := db.Exec(sqlUpdate, constants.LoanStatusPaid, strfmt.DateTime(currentDate))
	if err != nil {
		return err
	}
	return nil
}
