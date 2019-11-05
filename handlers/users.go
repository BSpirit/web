package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"web/models"
)

func FillDBHandler(env *Env, w http.ResponseWriter, r *http.Request) *StatusError {
	csvFile, err := os.Open("data.csv")
	if err != nil {
		return &StatusError{Code: 500, Err: fmt.Errorf("FillDB:\n\t%s", err)}
	}

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return &StatusError{Code: 500, Err: fmt.Errorf("FillDB:\n\t%s", err)}
	}

	for _, record := range records {
		user := models.User{
			Username: record[0],
			Age:      models.NewNullInt64(record[1]),
		}

		err := user.Create(env.DB)
		if err != nil {
			return &StatusError{Code: 500, Err: fmt.Errorf("FillDB:\n\t%s", err)}
		}

		products := make([]*models.Product, 0)
		err = json.Unmarshal([]byte(record[2]), &products)
		if err != nil {
			return &StatusError{Code: 500, Err: fmt.Errorf("FillDB:\n\t%s", err)}
		}

		for _, product := range products {
			product.UserID = user.ID
			err := product.Create(env.DB)
			if err != nil {
				return &StatusError{Code: 500, Err: fmt.Errorf("FillDB:\n\t%s", err)}
			}
		}
	}

	return nil
}

func CreateUserHandler(env *Env, w http.ResponseWriter, r *http.Request) *StatusError {
	if r.Method == http.MethodGet {
		env.Templates.ExecuteTemplate(w, "user_form", nil)
	} else if r.Method == http.MethodPost {
		user := models.User{
			Username: r.FormValue("username"),
			Age:      models.NewNullInt64(r.FormValue("age")),
		}

		err := user.Create(env.DB)
		if err != nil {
			return &StatusError{Code: 500, Err: fmt.Errorf("CreateUserHandler:\n\t%s", err)}
		}

		http.Redirect(w, r, "/users", http.StatusSeeOther)
	}

	return nil
}

func UsersHandler(env *Env, w http.ResponseWriter, r *http.Request) *StatusError {
	users, err := models.GetAllUsers(env.DB)
	if err != nil {
		return &StatusError{Code: 500, Err: fmt.Errorf("UsersHandler:\n\t%s", err)}
	}

	env.Templates.ExecuteTemplate(w, "users", users)

	return nil
}

func UserHandler(env *Env, w http.ResponseWriter, r *http.Request) *StatusError {
	s := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return &StatusError{Code: 500, Err: fmt.Errorf("UserHandler:\n\t%s", err)}
	}

	user, err := models.GetUser(id, env.DB)
	if err != nil {
		return &StatusError{Code: 500, Err: fmt.Errorf("UserHandler:\n\t%s", err)}
	}

	products, err := user.GetRelatedProducts(env.DB)
	if err != nil {
		return &StatusError{Code: 500, Err: fmt.Errorf("UserHandler:\n\t%s", err)}
	}

	env.Templates.ExecuteTemplate(w, "user", struct {
		User     *models.User
		Products []*models.Product
	}{user, products})

	return nil
}