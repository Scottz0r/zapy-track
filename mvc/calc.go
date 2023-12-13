package mvc

import (
	"fmt"
	za "go-zero-apr-mgr/zero-apr-lib"
	"math"
	"net/http"
	"strconv"
)

type PayBreakdownViewModel struct {
	PurName   string
	PayAmount string
}

type CalcViewModel struct {
	Accounts       []string
	HasOutput      bool
	OutAccount     string
	OutAmount      string
	OtherPurchases string
	Breakdown      []PayBreakdownViewModel
}

func (cntl *Controller) calcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cntl.calcGet(w, r)
	} else if r.Method == "POST" {
		cntl.calcPost(w, r)
	} else {
		cntl.ErrorPage(w, r, "Invalid operation")
	}
}

func (cntl *Controller) calcGet(w http.ResponseWriter, r *http.Request) {
	allPurchases := cntl.da.GetAllPurchases()

	vm := CalcViewModel{}
	vm.Accounts = getUniqueAccounts(allPurchases)

	cntl.templateMap["calc.html"].Execute(w, vm)
}

func (cntl *Controller) calcPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		cntl.ErrorPage(w, r, err.Error())
		return
	}

	accountName := r.Form.Get("account")
	if accountName == "" {
		cntl.ErrorPage(w, r, "Account not selected")
		return
	}

	rawBalance := r.Form.Get("balance")
	parsedAmount, err := strconv.ParseFloat(rawBalance, 64)
	if err != nil {
		cntl.ErrorPage(w, r, err.Error())
		return
	}

	balanceCurrency := int(math.Round(parsedAmount * 100))
	pmtAmount, otherPurchases, breakdowns := calcPayment(cntl.da, accountName, balanceCurrency)

	allPurchases := cntl.da.GetAllPurchases()

	vm := CalcViewModel{
		Accounts:       getUniqueAccounts(allPurchases),
		HasOutput:      true,
		OutAccount:     accountName,
		OtherPurchases: currencyToString(otherPurchases),
		OutAmount:      currencyToString(pmtAmount),
		Breakdown:      breakdowns,
	}

	cntl.templateMap["calc.html"].Execute(w, vm)
}

func getUniqueAccounts(allPurchases []za.Purchase) []string {
	accounts := make([]string, 0)
	seen := make(map[string]bool)
	for i := range allPurchases {
		if !seen[allPurchases[i].Account] {
			accounts = append(accounts, allPurchases[i].Account)
			seen[allPurchases[i].Account] = true
		}
	}

	return accounts
}

// TODO - this should probably be an engine method?
// TODO - better results. Put in a structure.
func calcPayment(da *za.DataAccess, account string, balanceCurrency int) (int, int, []PayBreakdownViewModel) {
	// TODO: Don't want to take in DA in engine method.
	allPurchases := da.GetAllPurchases()

	totalRemain := 0
	sumPayments := 0
	breakdowns := make([]PayBreakdownViewModel, 0)

	for i := range allPurchases {
		if allPurchases[i].Account == account {
			summary := za.CalculatePurchaseSummary(&allPurchases[i])
			sumPayments += summary.PeriodPaymentAmount
			totalRemain += summary.RemainingAmount

			bdvm := PayBreakdownViewModel{
				PurName:   allPurchases[i].Name,
				PayAmount: currencyToString(summary.PeriodPaymentAmount),
			}
			breakdowns = append(breakdowns, bdvm)
		}
	}

	if sumPayments > balanceCurrency {
		// TODO: If this is the case, then something is wrong with the data.
		fmt.Println("Balance mismatch")
	}

	otherBalance := balanceCurrency - totalRemain
	paymentAmount := otherBalance + sumPayments

	return paymentAmount, otherBalance, breakdowns
}

// Current balance - (total remain) + sumPayments
// 	  |-------------------|
//		  gets how much added not to tracked purchase
//		  then add in payments to make.
