package mvc

import (
	"math"
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
			}
		}
	}

	vm := PayViewModel{purchaseName, errorString}
	cntl.templateMap["pay.html"].Execute(w, vm)
}

func (cntl *Controller) tryMakePayment(r *http.Request) error {
	purchaseName := r.URL.Query().Get("name")

	// Form data:
	paymentDate := r.Form.Get("PaymentDate")
	paymentAmount := r.Form.Get("Amount")

	parsedTime, err := time.Parse("2006-01-02", paymentDate)
	if err != nil {
		return err
	}

	parsedAmount, err := strconv.ParseFloat(paymentAmount, 64)
	if err != nil {
		return err
	}

	amountCurrenccy := int(math.Round(parsedAmount * 100))

	pmt := za.Payment{PaymentDate: parsedTime, PaymentAmount: amountCurrenccy}
	err = cntl.da.AddPayment(purchaseName, pmt)
	if err != nil {
		return err
	}

	err = cntl.da.Save()
	if err != nil {
		return err
	}

	return nil
}
