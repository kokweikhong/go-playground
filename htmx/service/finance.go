package service

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/kokweikhong/go-playground/htmx/model"
)

type FinanceService interface {
	GetExpensesCategories() ([]*model.ExpensesCategory, error)
	GetExpensesCategory(id int) (*model.ExpensesCategory, error)
	UpdateExpensesCategory(id int, name, remarks string) error
	DeleteExpensesCategory(id int) error
	CreateExpensesCategory(name, remarks string) (*model.ExpensesCategory, error)
}

type financeService struct{}

func NewFinanceService() FinanceService {
	return &financeService{}
}

var categories []*model.ExpensesCategory

func init() {
	f := &financeService{}
	f.GetExpensesCategories()
}

func (s *financeService) GetExpensesCategories() ([]*model.ExpensesCategory, error) {
	if len(categories) > 0 {
		return categories, nil
	}
	var expensesCategories []*model.ExpensesCategory

	file, err := os.Open("data/expensesCategory.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(content, &expensesCategories); err != nil {
		return nil, err
	}

	categories = expensesCategories

	return expensesCategories, nil
}

func (s *financeService) GetExpensesCategory(id int) (*model.ExpensesCategory, error) {
	for _, v := range categories {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, errors.New("Category not found")
}

func (s *financeService) UpdateExpensesCategory(id int, name, remarks string) error {
	for _, v := range categories {
		if v.ID == id {
			v.Name = name
			v.Remarks = remarks
			return nil
		}
	}

	return errors.New("Category not found")
}

func (s *financeService) DeleteExpensesCategory(id int) error {
	for i, v := range categories {
		if v.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			return nil
		}
	}

	return errors.New("Category not found")
}

func (s *financeService) CreateExpensesCategory(name, remarks string) (*model.ExpensesCategory, error) {
	category := &model.ExpensesCategory{
		ID:        len(categories) + 1,
		Name:      name,
		Remarks:   remarks,
		CreatedAt: time.Now().Format("2006-01-02"),
	}
	categories = append(categories, category)
	slog.Info("categories", categories)
	return category, nil
}
