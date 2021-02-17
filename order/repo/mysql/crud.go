package mysql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"order/core"
)

func (m mysqlRepo) GetOrder(ctx context.Context, orderId int) (*core.Order, error) {

	order := &core.Order{}
	pay := &core.Payment{}

	err := m.db.QueryRowContext(ctx, "SELECT id, cart_id, buyer, amount, pay_method, status, created_at FROM "+m.table+" WHERE id = ?", orderId).
		Scan(&order.Id, &order.CartId, &order.Buyer, &pay.Total, &pay.Type, &order.Status, &order.CreatedAt)
	if err != nil {
		if err == rowsEmpty {
			return nil, nil
		}
		return nil, fmt.Errorf("GetOrder 1, %v", err)
	}

	rows, err := m.db.QueryContext(ctx, "SELECT id, name, price, qty FROM "+m.table+"_items WHERE order_id = ?", orderId)
	if err != nil {
		return nil, fmt.Errorf("GetOrder 2, %v", err)
	}

	for rows.Next() {
		item := &core.Item{}

		if err = rows.Scan(&item.Id, &item.Name, &item.Price, &item.Qty); err != nil {
			return nil, fmt.Errorf("GetOrder 3, %v", err)
		}
		order.Items = append(order.Items, item)
	}

	order.Payment = pay

	return order, nil
}

func (m mysqlRepo) AllOrders(ctx context.Context, buyer int) (core.Orders, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT id, cart_id, buyer, amount, pay_method, status, created_at FROM "+m.table+" WHERE buyer = ?", buyer)
	if err != nil {
		return nil, err
	}

	var orders core.Orders
	for rows.Next() {
		order := &core.Order{}
		pay := &core.Payment{}

		if err = rows.Scan(&order.Id, &order.CartId, &order.Buyer, &pay.Total, &pay.Type, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}

		lines, err := m.db.QueryContext(ctx, "SELECT id, name, price, qty FROM "+m.table+"_items WHERE order_id = ?", order.Id)
		if err != nil {
			return nil, err
		}

		for lines.Next() {
			item := &core.Item{}

			if err = lines.Scan(&item.Id, &item.Name, &item.Price, &item.Qty); err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}

		order.Payment = pay
		orders = append(orders, order)
	}

	return orders, nil
}

func (m mysqlRepo) PlaceOrder(ctx context.Context, order *core.Order) (id int, err error) {

	var s []string
	s = append(s, "cart_id = "+strconv.Itoa(order.CartId))
	s = append(s, "buyer = "+strconv.Itoa(order.Buyer))
	s = append(s, "amount = "+strconv.FormatFloat(order.Payment.Total, 'f', -1, 64))
	s = append(s, "pay_method = "+strconv.Itoa(int(order.Payment.Type)))
	s = append(s, "status = "+strconv.Itoa(int(order.Status)))

	rows, err := m.db.ExecContext(ctx, "INSERT "+m.table+" SET "+strings.Join(s, ","))

	if err != nil {
		return
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return 0, NoRowsAffected
	}

	oid, _ := rows.LastInsertId()
	id = int(oid)

	for _, item := range order.Items {
		var a []string
		a = append(a, "id = "+strconv.Itoa(item.Id))
		a = append(a, "order_id = "+strconv.Itoa(id))
		a = append(a, "name = '"+item.Name+"'")
		a = append(a, "price = "+strconv.FormatFloat(item.Price, 'f', -1, 64))
		a = append(a, "qty = "+strconv.Itoa(int(item.Qty)))

		rows, err = m.db.ExecContext(ctx, "INSERT "+m.table+"_items SET "+strings.Join(a, ","))
	}

	return
}
