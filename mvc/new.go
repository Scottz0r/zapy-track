package mvc

import (
	zeroaprlib "go-zero-apr-mgr/zero-apr-lib"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (cntl *Controller) newHandler(w http.ResponseWriter, r *http.Request) {
	// Get method is simple. Just return the template.
	if r.Method == "GET" {
		cntl.newGet(w, r)
	} else if r.Method == "POST" {
		cntl.newPost(w, r)
	} else {
		// Not get/post, so nothing can be done.
		r.Response.StatusCode = http.StatusBadRequest
	}
}

func (cntl *Controller) newGet(w http.ResponseWriter, r *http.Request) {
	cntl.templateMap["new.html"].Execute(w, nil)
}

func (cntl *Controller) newPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// TODO: Better errors!
		cntl.ErrorPage(w, r, err.Error())
		return
	}

	// Need to make sure a purchase with this name does not already exist.
	purName := strings.TrimSpace(r.Form.Get("name"))
	if purName == "" {
		// TODO: better errors
		cntl.ErrorPage(w, r, "Invalid name")
		return
	}

	if cntl.da.PurchaseExists(purName) {
		// TODO: Better errors!
		cntl.ErrorPage(w, r, "A purchase with this name already exists")
		return
	}

	parsedDate, err := time.Parse("2006-01-02", r.Form.Get("purchaseDate"))
	if err != nil {
		// TODO: Better errors!
		cntl.ErrorPage(w, r, err.Error())
		return
	}

	parsedAmount, err := strconv.ParseFloat(r.Form.Get("amount"), 64)
	if err != nil {
		// TODO: Better errors!
		cntl.ErrorPage(w, r, err.Error())
		return
	}
	amountCurrency := floatToCurrency(parsedAmount)

	parsedNumPayments, err := strconv.ParseInt(r.Form.Get("numPayments"), 10, 32)
	if err != nil {
		// TODO: Better errors!
		cntl.ErrorPage(w, r, err.Error())
		return
	}

	newPurchase := zeroaprlib.Purchase{
		Name:              r.Form.Get("name"),
		Account:           r.Form.Get("account"),
		Memo:              r.Form.Get("memo"),
		Date:              parsedDate,
		InitialAmount:     amountCurrency,
		WantedNumPayments: int(parsedNumPayments),
	}

	err = cntl.da.AddPurchase(newPurchase)
	if err != nil {
		cntl.ErrorPage(w, r, err.Error())
		return
	}

	err = cntl.da.Save()
	if err != nil {
		cntl.ErrorPage(w, r, err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
