package util

import (
	"aspire/constants"
	"aspire/models"
	"time"

	"github.com/go-openapi/strfmt"
)

func GetInstallments(loanId int64, loanAmount float64, term int64) []*models.Installment {
	today := GetDate(time.Now())
	installments := []*models.Installment{}
	installmentAmount := loanAmount / float64(term)
	for i := 0; i < int(term); i++ {
		installment := &models.Installment{}
		installment.ID = GetInt64Pointer(loanId)
		installment.InstallmentAmount = &installmentAmount
		installment.DueDate = GetDateTimePointer(strfmt.DateTime(today.AddDate(0, 0, int(term)*7)))
		installment.State = constants.InstallmentStatusPending
		installments = append(installments, installment)
	}
	return installments
}
