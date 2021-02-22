package mysql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"cart/core"
)

func (m mysqlRepo) Create(ctx context.Context, buyer int) (id int, err error) {

	err = m.db.QueryRowContext(ctx, "SELECT id FROM "+m.table+" WHERE buyer = '"+ strconv.Itoa(buyer)+"'").
		Scan(&id)
	if err == nil && id != 0 {
		return
	}

	rows, err := m.db.ExecContext(ctx, "INSERT "+m.table+" SET buyer = '"+ strconv.Itoa(buyer)+"'")

	if err != nil {
		return 0, err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return 0, NoRowsAffected
	}

	idx, _ := rows.LastInsertId()

	return int(idx), nil
}

func (m mysqlRepo) ByID(ctx context.Context, cartId int) (core.Cart, error) {
	cart := core.Cart{Id: cartId}

	err := m.db.QueryRowContext(ctx, "SELECT buyer, created_at FROM "+m.table+" WHERE id = ?", cartId).
		Scan(&cart.Buyer, &cart.CreatedAt)
	if err != nil {
		return cart, err
	}

	rows, err := m.db.QueryContext(ctx, "SELECT id, name, price, old_price, stocks, qty FROM "+m.table+"_items WHERE cart_id = ?", cartId)
	if err != nil {
		return cart, err
	}

	for rows.Next() {
		item := &core.Item{}

		if err = rows.Scan(&item.Id, &item.Name, &item.Price, &item.OldPrice, &item.Stocks, &item.Qty); err != nil {
			return cart, err
		}
		cart.Items = append(cart.Items, item)
	}

	return cart, nil
}

func (m mysqlRepo) AddItem(ctx context.Context, cartId int, item *core.Item) (err error) {
	s := strings.Builder{}
	s.WriteString("id = '"+strconv.Itoa(item.Id)+"'")
	s.WriteString(",cart_id = '"+strconv.Itoa(cartId)+"'")
	s.WriteString(",name = '"+item.Name+"'")
	price := strconv.FormatFloat(item.Price, 'f', -1, 64)
	s.WriteString(",price = '"+price+"'")
	s.WriteString(",stocks = '"+strconv.Itoa(item.Stocks)+"'")
	s.WriteString(",qty = '"+strconv.Itoa(item.Qty)+"'")

	rows, err := m.db.ExecContext(ctx, "INSERT "+m.table+"_items SET "+s.String())

	if err != nil {
		return
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return NoRowsAffected
	}
	return
}

func (m mysqlRepo) UpdateItems(ctx context.Context, item *core.Item) error {
	var s []string
	if item.Id == 0 {
		return noIdUpdate
	}

	if item.Name != "" {
		s = append(s, "name = '"+item.Name+"'")
	}
	if item.Price != 0 {
		oldPrice := "0"
		m.db.QueryRowContext(ctx, "SELECT price FROM "+m.table+"_items WHERE id = ?", item.Id).
			Scan(&oldPrice)

		price := strconv.FormatFloat(item.Price, 'f', -1, 64)
		s = append(s, "old_price = '"+oldPrice+"'")
		s = append(s, "price = '"+price+"'")
	}
	if item.Stocks != 0 {
		s = append(s, "stocks = '"+strconv.Itoa(item.Stocks)+"'")
	}

	rows, err := m.db.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s_items SET %s WHERE id = ?", m.table, strings.Join(s, ",")), item.Id,
	)

	if err != nil {
		return err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return NoRowsAffected
	}

	return nil
}

func (m mysqlRepo) RemoveItem(ctx context.Context, cartId int, itemId int) (err error) {

	rows, err := m.db.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM %s_items WHERE cart_id = ? AND id = ?", m.table), cartId, itemId,
		)

	if err != nil {
		return
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return NoRowsAffected
	}
	return
}

func (m mysqlRepo) ResetCart(ctx context.Context, cartId int) (err error) {
	rows, err := m.db.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM %s_items WHERE cart_id = ?", m.table), cartId,
	)

	if err != nil {
		return
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return NoRowsAffected
	}
	return
}

func (m mysqlRepo) UpdateQty(ctx context.Context, cartId int, item *core.Item) error {
	if item.Id == 0 {
		return noIdUpdate
	}

	if item.Qty == 0 {
		return m.RemoveItem(ctx, cartId, item.Id)
	}

	rows, err := m.db.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s_items SET qty = ? WHERE id = ?", m.table), strconv.Itoa(item.Qty), item.Id,
	)

	if err != nil {
		return err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return NoRowsAffected
	}

	return nil
}
