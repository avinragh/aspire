package util

import (
	"aspire/constants"
	"aspire/models"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
)

func GetInstallments(loanId int64, loanAmount float64, term int64) []*models.Installment {
	today := GetDate(time.Now())
	defaultTime := strfmt.DateTime(time.Time{})
	installments := []*models.Installment{}
	installmentAmount := loanAmount / float64(term)
	repaymentAmount := float64(0)
	for i := 0; i < int(term); i++ {
		installment := &models.Installment{}
		installment.ID = loanId
		installment.InstallmentAmount = &installmentAmount
		installment.DueDate = GetDateTimePointer(strfmt.DateTime(today.AddDate(0, 0, i*7)))
		installment.RepaymentAmount = &repaymentAmount
		installment.State = constants.InstallmentStatusPending
		installment.RepaymentTime = defaultTime
		installments = append(installments, installment)
	}
	return installments
}

func GetLoanIdsClause(loanIds []int64) string {
	loanIdString := ""
	if len(loanIds) > 0 {
		loanIdString = "("
		for i, loanId := range loanIds {
			loanIdString = loanIdString + strconv.FormatInt(loanId, 10)
			if i != len(loanIds)-1 {
				loanIdString = loanIdString + ","
			} else {
				loanIdString = loanIdString + ")"
			}
		}
	}
	return loanIdString

}
