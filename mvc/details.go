package mvc

import (
	za "go-zero-apr-mgr/zero-apr-lib"
	"net/http"
)

type PaymentDetailViewModel struct {
	Date   string
	Amount string
}

type PurchaseDetailViewModel struct {
	Name                string
	Account             string
	Memo                string
	PurchaseDate        string
	InitialAmount       string
	WantedPayments      int
	RemainingAmount     string
	TotalPaid           string
	TotalPayments       int
	PeriodPaymentAmount string

	Payments []PaymentDetailViewModel
}

func (cntl *Controller) detailsHandler(w http.ResponseWriter, r *http.Request) {
	purchaseName := r.URL.Query().Get("name")

	// Send 404 if the purchase does not exist.
	pur := cntl.da.GetPurchase(purchaseName)
	if pur == nil {
		r.Response.StatusCode = http.StatusNotFound
		return
	}

	summary := za.CalculatePurchaseSummary(pur)

	vm := PurchaseDetailViewModel{
		Name:                pur.Name,
		Account:             pur.Account,
		Memo:                pur.Memo,
		PurchaseDate:        pur.Date.Format("2006-01-02"),
		InitialAmount:       currencyToString(pur.InitialAmount),
		WantedPayments:      pur.WantedNumPayments,
		RemainingAmount:     currencyToString(summary.RemainingAmount),
		TotalPaid:           currencyToString(summary.TotalPaid),
		TotalPayments:       summary.TotalPayments,
		PeriodPaymentAmount: currencyToString(summary.PeriodPaymentAmount),
	}

	for i := range pur.Payments {
		payVm := PaymentDetailViewModel{
			Date:   pur.Payments[i].PaymentDate.Format("2006-01-02"),
			Amount: currencyToString(pur.Payments[i].PaymentAmount),
		}
		vm.Payments = append(vm.Payments, payVm)
	}

	cntl.templateMap["details.html"].Execute(w, vm)
}
