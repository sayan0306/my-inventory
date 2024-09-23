package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type Product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]Product, error) {
	query := "SELECT * FROM products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Quantity, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *Product) getProductById(db *sql.DB) error {
	query := fmt.Sprintf("SELECT * FROM products WHERE id=%v", p.Id)
	rows := db.QueryRow(query)
	err := rows.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf(`INSERT INTO products (name, quantity, price) VALUES ('%v', %v, %v)`, p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.Id = int(id)
	return nil
}

func (p *Product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf(`UPDATE products SET name ='%v', quantity =%v, price =%v WHERE id =%v`, p.Name, p.Quantity, p.Price, p.Id)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("no such product exists")
	}
	return nil
}

func (p *Product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf(`DELETE from products WHERE id =%v`, p.Id)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("no such product exists")
	}
	return nil
}
