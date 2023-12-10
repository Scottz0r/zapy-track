package main

import (
	"fmt"
	zeroaprlib "go-zero-apr-mgr/zero-apr-lib"
	"html/template"
	"net/http"
	"path"
)

// Alias for the map used to look up loaded templates.
type TemplateMap map[string]*template.Template

type Controller struct {
	da          *zeroaprlib.DataAccess
	templateMap TemplateMap
}

type PurchaseViewModel struct {
	Name            string
	RemainingAmount string
	TotalPaid       string
	TotalPayments   int
	PaymentAmount   string
}

type HomeViewModel struct {
	Purchases []PurchaseViewModel
}

func (controller *Controller) homeHandler(w http.ResponseWriter, r *http.Request) {
	vm := HomeViewModel{}

	purchases := controller.da.GetAllPurchases()
	for i := range purchases {
		summary := zeroaprlib.CalculatePurchaseSummary(&purchases[i])

		purVm := PurchaseViewModel{
			Name:            purchases[i].Name,
			RemainingAmount: currencyToString(summary.RemainingAmount),
			TotalPaid:       currencyToString(summary.TotalPaid),
			TotalPayments:   summary.TotalPayments,
			PaymentAmount:   currencyToString(summary.PeriodPaymentAmount),
		}

		vm.Purchases = append(vm.Purchases, purVm)
	}

	controller.templateMap["home.html"].Execute(w, vm)
}

func currencyToString(value int) string {
	floatVal := convertFromCurrency(value)
	return fmt.Sprintf("%.2f", floatVal)
}

func serverMain(da *zeroaprlib.DataAccess) {
	templateMap, err := loadTemplates()
	if err != nil {
		fmt.Print(err)
		return
	}

	c := Controller{da, templateMap}

	http.HandleFunc("/", c.homeHandler)

	fmt.Println(http.ListenAndServe("localhost:8080", nil))
}

func loadTemplates() (TemplateMap, error) {
	templateMap := make(TemplateMap)

	templateCofnig := []string{
		"home.html",
	}

	for _, filename := range templateCofnig {
		fullPath := path.Join("templates", filename)
		temp, err := template.ParseFiles(fullPath)
		if err != nil {
			return nil, err
		}

		templateMap[filename] = temp
	}

	return templateMap, nil
}
