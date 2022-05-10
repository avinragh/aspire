package crons

import (
	"aspire/constants"
	"aspire/context"
)

func UpdatePaidLoans(ctx *context.Context) {
	logger := ctx.GetLogger()
	database := ctx.GetDB()
	loans, err := database.FindLoans(0, nil, nil, nil, nil)
	if err != nil {
		logger.Println(err)
		return
	}
	for _, loan := range loans {
		if loan.State == constants.LoanStatusApproved && *loan.InstallmentsCreated {
			installments, err := database.FindInstallments(&loan.ID, nil, nil, nil, nil, []int64{})
			if err != nil {
				logger.Println(err)
				return
			}
			incompletePayments := false
			for _, installment := range installments {
				if installment.State == constants.InstallmentStatusPending {
					incompletePayments = true
					break
				}
			}
			if !incompletePayments {
				err := database.UpdateToPaid(loan.ID)
				if err != nil {
					logger.Println(err)
				}
			}
		}
	}
}
