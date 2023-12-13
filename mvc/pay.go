package mvc

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	za "go-zero-apr-mgr/zero-apr-lib"
)

type PayViewModel struct {
	Name         string
	ErrorMessage string
}

func (cntl *Controller) payHandler(w http.ResponseWriter, r *http.Request) {
	purchaseName := r.URL.Query().Get("name")
	// TODO: make sure purchase exists.

	var errorString string

	// TODO - this needs to be refactored
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			errorString = err.Error()
		} else {
			err = cntl.tryMakePayment(r)
			if err != nil {
				errorString = err.Error()
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		cntl.ErrorPage(w, err.Error())
		return
	} else if r.Method == "GET" {
		vm := PayViewModel{purchaseName, errorString}
		cntl.ExecuteView(w, "pay", vm)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (cntl *Controller) tryMakePayment(r *http.Request) error {
	purchaseName := r.URL.Query().Get("name")

	paymentDate := r.Form.Get("paymentDate")
	paymentAmount := r.Form.Get("amount")

	parsedTime, err := time.Parse("2006-01-02", paymentDate)
	if err != nil {
		return err
	}

	parsedAmount, err := strconv.ParseFloat(paymentAmount, 64)
	if err != nil {
		return err
	}

	amountCurrenccy := floatToCurrency(parsedAmount)
	if amountCurrenccy <= 0 {
		return fmt.Errorf("Payments cannot be zero or negative")
	}

	pmt := za.Payment{PaymentDate: parsedTime, PaymentAmount: amountCurrenccy}
	err = cntl.da.AddPayment(purchaseName, pmt)
	if err != nil {
		return err
	}

	return nil
}
