package crons

import (
	"aspire/constants"
	"aspire/context"
	"aspire/util"
	packcontext "context"
	"time"

	"github.com/go-openapi/strfmt"
)

func InsertInstallments(ctx *context.Context) {
	logger := ctx.GetLogger()
	database := ctx.GetDB()
	state := constants.LoanStatusApproved
	loans, err := database.FindLoans(0, &state, nil, nil, nil)
	if err != nil {
		logger.Println(err)
		return
	}
	for _, loan := range loans {
		if !*loan.InstallmentsCreated {
			currentTime := time.Now().UTC()
			sqlInsert := `
				INSERT INTO installments(installment_amount,repayment_amount,due_date,state,loan_id,created_on,modified_on)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`

			sqlUpdate := `
				UPDATE loans
				SET installments_created = true, modified_on = $1
				WHERE id = $2`

			loan.ModifiedOn = strfmt.DateTime(currentTime)

			installments := util.GetInstallments(loan.ID, *loan.Amount, *loan.Term)
			tctx := packcontext.Background()
			tx, err := database.BeginTx(tctx, nil)
			if err != nil {
				logger.Println(err)
			}
			for _, installment := range installments {
				installment.CreatedOn = strfmt.DateTime(currentTime)
				installment.ModifiedOn = strfmt.DateTime(currentTime)

				_, err = tx.ExecContext(tctx, sqlInsert, installment.InstallmentAmount, installment.RepaymentAmount, installment.DueDate, installment.State, loan.ID, installment.CreatedOn, installment.ModifiedOn)
				if err != nil {
					tx.Rollback()
					logger.Println(err)
				}
			}
			_, err = tx.ExecContext(tctx, sqlUpdate, loan.ModifiedOn, loan.ID)
			if err != nil {
				tx.Rollback()
				logger.Println(err)
			}

			err = tx.Commit()
			if err != nil {
				logger.Println(err)
			}

		}
	}
}
