package zeroaprlib

type PurchaseCalcs struct {
	TotalPaid            int
	TotalPayments        int
	RemainingAmount      int
	RemainingNumPayments int
	PeriodPaymentAmount  int
	WarningMessage       string
}

func CalculatePurchaseSummary(purchase *Purchase) PurchaseCalcs {

	result := PurchaseCalcs{}

	result.TotalPayments = len(purchase.Payments)
	if result.TotalPayments < purchase.WantedNumPayments {
		result.RemainingNumPayments = purchase.WantedNumPayments - result.TotalPayments
	}

	result.TotalPaid = 0
	for i := range purchase.Payments {
		result.TotalPaid += purchase.Payments[i].PaymentAmount
	}

	result.RemainingAmount = purchase.InitialAmount - result.TotalPaid

	if result.RemainingNumPayments > 0 {
		result.PeriodPaymentAmount = result.RemainingAmount / result.RemainingNumPayments
	} else {
		result.PeriodPaymentAmount = 0
	}

	if result.RemainingAmount < 0 {
		result.WarningMessage = "overpaid"
	}

	if result.RemainingNumPayments == 0 && result.RemainingAmount > 0 {
		result.WarningMessage = "underpaid"
	}

	return result
}
