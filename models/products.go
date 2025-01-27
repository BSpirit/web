package models

import (
	"database/sql"
	"errors"
	"web/utils"
)

type Product struct {
	ID     int64
	Name   string
	UserID int64
}

func (p *Product) Create(tx *sql.Tx) error {
	res, err := tx.Exec("INSERT INTO products(name, user_id) VALUES(?, ?)", p.Name, p.UserID)
	if err != nil {
		return utils.Trace(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return utils.Trace(err)
	}
	p.ID = id

	return nil
}

func (p *Product) Update(tx *sql.Tx) error {
	if p.ID == 0 {
		return errors.New("no ID")
	}
	_, err := tx.Exec("UPDATE products SET name=?, user_id=? WHERE id=?", p.Name, p.UserID, p.ID)
	if err != nil {
		return utils.Trace(err)
	}

	return nil
}

func (p *Product) Delete(tx *sql.Tx) error {
	if p.ID == 0 {
		return errors.New("no ID")
	}
	_, err := tx.Exec("DELETE FROM products WHERE id=?", p.ID)
	if err != nil {
		return utils.Trace(err)
	}
	p.ID = 0

	return nil
}

func GetProduct(pk int64, db *sql.DB) (*Product, error) {
	product := &Product{}
	err := db.QueryRow("SELECT * FROM products WHERE id = ?", pk).Scan(&product.ID, &product.Name, &product.UserID)
	if err != nil {
		return nil, utils.Trace(err)
	}

	return product, nil
}

func GetAllProducts(db *sql.DB) ([]*Product, error) {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		return nil, utils.Trace(err)
	}
	defer rows.Close()
	products := make([]*Product, 0)
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.UserID)
		if err != nil {
			return nil, utils.Trace(err)
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, utils.Trace(err)
	}

	return products, nil
}
