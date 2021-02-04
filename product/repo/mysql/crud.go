package mysql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"product/core"
)

func (m mysqlRepo) All(ctx context.Context) (core.Products, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT id, sku, name, price, stocks, created_at FROM "+m.table+" WHERE 1")
	if err != nil {
		return nil, err
	}

	var prods []core.Product
	for rows.Next() {
		var prod core.Product

		if err = rows.Scan(&prod.Id, &prod.Sku, &prod.Name, &prod.Price, &prod.Stocks, &prod.CreatedAt); err != nil {
			return prods, err
		}

		prods = append(prods, prod)
	}

	return prods, nil
}

func (m mysqlRepo) ByID(ctx context.Context, i int) (core.Product, error) {

	var prod core.Product

	err := m.db.QueryRowContext(ctx, "SELECT id, sku, name, price, stocks, created_at FROM "+m.table+" WHERE id=?", i).
		Scan(&prod.Id, &prod.Sku, &prod.Name, &prod.Price, &prod.Stocks, &prod.CreatedAt)
	if err != nil {
		return prod, err
	}
	return prod, nil
}

func (m mysqlRepo) Add(ctx context.Context, prod *core.Product) (int, error) {

	s := strings.Builder{}
	s.WriteString("sku = '" + prod.Sku + "'")
	s.WriteString(",name = '" + prod.Name + "'")
	price := strconv.FormatFloat(prod.Price, 'g', 1, 64)
	s.WriteString(",price = '" + price + "'")
	s.WriteString(",stocks = '" + strconv.Itoa(prod.Stocks) + "'")

	rows, err := m.db.ExecContext(ctx, "INSERT "+m.table+" SET "+s.String())

	if err != nil {
		return 0, err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return 0, NoRowsAffected
	}

	id, _ := rows.LastInsertId()

	return int(id), nil
}

func (m mysqlRepo) Edit(ctx context.Context, prod *core.Product) error {

	var s []string
	if prod.Id == 0 {
		return noIdUpdate
	}
	if prod.Sku != "" {
		s = append(s, "sku = '"+prod.Sku+"'")
	}
	if prod.Name != "" {
		s = append(s, "name = '"+prod.Name+"'")
	}
	if prod.Price != 0 {
		price := strconv.FormatFloat(prod.Price, 'g', 1, 64)
		s = append(s, "price = '"+price+"'")
	}
	if prod.Stocks != 0 {
		s = append(s, "stocks = '"+strconv.Itoa(prod.Stocks)+"'")
	}

	rows, err := m.db.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", m.table, strings.Join(s, ","), prod.Id),
	)

	if err != nil {
		return err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return NoRowsAffected
	}

	return err
}
