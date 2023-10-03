package controller

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kokweikhong/go-playground/htmx/model"
	"github.com/kokweikhong/go-playground/htmx/service"
)

type BaseController interface {
	Index(w http.ResponseWriter, r *http.Request)
	NewExpensesCategory(w http.ResponseWriter, r *http.Request)
	CreateExpensesCategory(w http.ResponseWriter, r *http.Request)
	UpdateExpensesCategory(w http.ResponseWriter, r *http.Request)
	DeleteExpensesCategory(w http.ResponseWriter, r *http.Request)
}

type baseController struct{}

func NewBaseController() BaseController {
	return &baseController{}
}

func (c *baseController) Index(w http.ResponseWriter, r *http.Request) {
	slog.Info("Index page")
	tpl := template.Must(template.ParseFiles(
		"views/index.html",
		"views/layouts/css.html",
		"views/layouts/js.html",
		"views/blocks/expenses-category-form.html",
	))

	for _, v := range tpl.Templates() {
		slog.Info(v.Name())
	}

	financeService := service.NewFinanceService()

	financeExpensesCategories, err := financeService.GetExpensesCategories()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	category := new(model.ExpensesCategory)

	pageData := map[string]interface{}{
		"Title":      "HTMX Playground",
		"Categories": financeExpensesCategories,
		"Category":   category,
	}

	tpl.ExecuteTemplate(w, "index.html", pageData)
}

func (c *baseController) NewExpensesCategory(w http.ResponseWriter, r *http.Request) {
	slog.Info("NewExpensesCategory page")

	formType := chi.URLParam(r, "formType")
	slog.Info("formType: ", formType)

	var (
		parseString string
		category    *model.ExpensesCategory
	)

	switch formType {
	case "create":
		parseString = `{{define "expenses-category-form"}}{{template "blocks/expenses-category-form" .}}{{end}}`
		category = new(model.ExpensesCategory)
	case "update":
		parseString = `{{define "expenses-category-form"}}{{template "blocks/expenses-category-form" .}}{{end}}
		 {{define "blocks/expenses-category-form-submit"}}{{template "blocks/expenses-category-form-submit/update" .Category}}{{end}}`
		idParam := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		financeService := service.NewFinanceService()
		category, err = financeService.GetExpensesCategory(id)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	slog.Info("NewExpensesCategory page")
	tpl, err := template.Must(template.ParseFiles(
		"views/blocks/expenses-category-form.html",
	)).New("blocks/expenses-category-form").Parse(parseString)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := map[string]interface{}{
		"Title":    "HTMX Playground",
		"Category": category,
	}

	tpl.ExecuteTemplate(w, "blocks/expenses-category-form", pageData)
}

func (c *baseController) CreateExpensesCategory(w http.ResponseWriter, r *http.Request) {
	slog.Info("CreateExpensesCategory page")

	r.ParseForm()

	name := r.FormValue("name")
	remarks := r.FormValue("remarks")

	financeService := service.NewFinanceService()

	category, err := financeService.CreateExpensesCategory(name, remarks)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tpl := template.Must(template.ParseFiles(
		"views/blocks/expenses-category-form.html",
	))

	tpl.ExecuteTemplate(w, "blocks/expenses-category-form-submit/card", category)
}

func (c *baseController) UpdateExpensesCategory(w http.ResponseWriter, r *http.Request) {
	slog.Info("UpdateExpensesCategory page")

	r.ParseForm()
	name := r.FormValue("name")
	remarks := r.FormValue("remarks")
	idParam := r.FormValue("id")
	slog.Info("id: ", idParam)
	slog.Info("name: ", name)
	slog.Info("remarks: ", remarks)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	financeService := service.NewFinanceService()

	err = financeService.UpdateExpensesCategory(id, name, remarks)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	category, err := financeService.GetExpensesCategory(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tpl := template.Must(template.ParseFiles(
		"views/blocks/expenses-category-form.html",
	))

	tpl.ExecuteTemplate(w, "blocks/expenses-category-form-submit/card", category)
}

func (c *baseController) DeleteExpensesCategory(w http.ResponseWriter, r *http.Request) {
	slog.Info("DeleteExpensesCategory page")

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	financeService := service.NewFinanceService()

	err = financeService.DeleteExpensesCategory(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// http.Redirect(w, r, "/", http.StatusSeeOther)
}
