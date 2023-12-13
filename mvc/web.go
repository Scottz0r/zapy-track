package mvc

import (
	"fmt"
	zeroaprlib "go-zero-apr-mgr/zero-apr-lib"
	"html/template"
	"math"
	"net/http"
	"path"
)

// Alias for the map used to look up loaded templates.
type TemplateMap map[string]*template.Template

type Controller struct {
	da *zeroaprlib.DataAccess
}

func currencyToString(value int) string {
	floatVal := float64(value) / 100
	return fmt.Sprintf("%.2f", floatVal)
}

func floatToCurrency(value float64) int {
	return int(math.Round(value * 100.0))
}

func ServerMain(da *zeroaprlib.DataAccess) {
	c := Controller{da}

	// Server static files for CSS/JS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", c.homeHandler)
	http.HandleFunc("/pay", c.payHandler)
	http.HandleFunc("/details", c.detailsHandler)
	http.HandleFunc("/new", c.newHandler)
	http.HandleFunc("/calc", c.calcHandler)

	fmt.Println(http.ListenAndServe("localhost:8080", nil))
}

// Method used to show a general error page.
func (cntl *Controller) ErrorPage(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)

	_, err := fmt.Fprintf(w, "Error: %s", message)
	if err != nil {
		fmt.Println("Error when showing error page:", err)
	}
}

func (cntl *Controller) ExecuteView(w http.ResponseWriter, templateName string, data any) {
	// TODO: Could optimize and cache templates. But for now, do this for hot reloading.
	filename := templateName + ".html"
	fullPath := path.Join("templates", filename)

	// Add in the template file to share common functionality.
	templatePath := path.Join("templates", "_common.html")

	temp, err := template.ParseFiles(fullPath, templatePath)
	if err != nil {
		// TODO
		panic("Template bombed out")
	}

	err = temp.Execute(w, data)
	if err != nil {
		// TODO
		panic("Template write fail")
	}
}
