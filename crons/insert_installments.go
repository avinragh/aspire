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
	loans, err := database.FindLoans(0, constants.LoanStatusApproved, "", 0, 0)
	if err != nil {
		logger.Println(err)
		return
	}
	for _, loan := range loans {
		installments, err := database.FindInstallments(*loan.ID, "", "", 0, 0)
		if err != nil {
			logger.Println(err)
			return
		}
		if installments != nil && len(installments) == 0 {
			currentTime := time.Now().UTC()
			sqlInsert := `
			INSERT INTO installments(installment_amount,due_date,state,loan_id,created_on,modified_on)
			VALUES ($1, $2, $3, $4, $5)`

			installments := util.GetInstallments(*loan.ID, *loan.Amount, *loan.Term)
			tctx := packcontext.Background()
			tx, err := database.BeginTx(tctx, nil)
			if err != nil {
				logger.Println(err)
			}
			for _, installment := range installments {
				installment.CreatedOn = strfmt.DateTime(currentTime)
				installment.ModifiedOn = strfmt.DateTime(currentTime)
				_, err = tx.ExecContext(tctx, sqlInsert, installment.InstallmentAmount, installment.DueDate, installment.State, *loan.ID, installment.CreatedOn, installment.ModifiedOn)
				if err != nil {
					tx.Rollback()
					logger.Println(err)
				}
			}
			err = tx.Commit()
			if err != nil {
				logger.Println(err)
			}
		}
	}
}
