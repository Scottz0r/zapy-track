package mvc

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

func currencyToString(value int) string {
	floatVal := float64(value) / 100
	return fmt.Sprintf("%.2f", floatVal)
}

func ServerMain(da *zeroaprlib.DataAccess) {
	templateMap, err := loadTemplates()
	if err != nil {
		fmt.Print(err)
		return
	}

	c := Controller{da, templateMap}

	http.HandleFunc("/", c.homeHandler)
	http.HandleFunc("/pay", c.payHandler)
	http.HandleFunc("/details", c.detailsHandler)

	fmt.Println(http.ListenAndServe("localhost:8080", nil))
}

func loadTemplates() (TemplateMap, error) {
	templateMap := make(TemplateMap)

	templateCofnig := []string{
		"home.html",
		"pay.html",
		"details.html",
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
