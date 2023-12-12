package mvc

import (
	za "go-zero-apr-mgr/zero-apr-lib"
	"net/http"
)

type PurchaseSummaryViewModel struct {
	Name            string
	InitialAmount   string
	WantedPayments  int
	RemainingAmount string
	TotalPaid       string
	TotalPayments   int
	PaymentAmount   string
}

type HomeViewModel struct {
	Purchases []PurchaseSummaryViewModel
}

func (controller *Controller) homeHandler(w http.ResponseWriter, r *http.Request) {
	vm := HomeViewModel{}

	purchases := controller.da.GetAllPurchases()
	for i := range purchases {
		summary := za.CalculatePurchaseSummary(&purchases[i])

		purVm := PurchaseSummaryViewModel{
			Name:            purchases[i].Name,
			InitialAmount:   currencyToString(purchases[i].InitialAmount),
			WantedPayments:  purchases[i].WantedNumPayments,
			RemainingAmount: currencyToString(summary.RemainingAmount),
			TotalPaid:       currencyToString(summary.TotalPaid),
			TotalPayments:   summary.TotalPayments,
			PaymentAmount:   currencyToString(summary.PeriodPaymentAmount),
		}

		vm.Purchases = append(vm.Purchases, purVm)
	}

	controller.templateMap["home.html"].Execute(w, vm)
}
