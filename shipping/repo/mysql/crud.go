package mysql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"shipping/core"
)

func (m mysqlRepo) AddAddress(ctx context.Context, address *core.Address) (int, error) {
	s := strings.Builder{}
	s.WriteString("user = '" + strconv.Itoa(address.User) + "'")
	s.WriteString(",contact_name = '" + address.ContactName + "'")
	s.WriteString(",contact_phone = '" + address.ContactPhone + "'")
	s.WriteString(",landmark = '" + address.Landmark + "'")
	s.WriteString(",city = '" + address.City + "'")
	s.WriteString(",state = '" + address.State + "'")
	s.WriteString(",country = '" + address.Country + "'")
	s.WriteString(",zip = '" + strconv.Itoa(address.Zip) + "'")

	rows, err := m.db.ExecContext(ctx, "INSERT "+m.table+"_addresses SET "+s.String())

	if err != nil {
		return 0, err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return 0, NoRowsAffected
	}

	id, _ := rows.LastInsertId()

	return int(id), nil
}

func (m mysqlRepo) AddOrderShipping(ctx context.Context, ship *core.Shipping) (int, error) {
	s := strings.Builder{}
	s.WriteString("order_id = '" + strconv.Itoa(ship.OrderId) + "'")
	s.WriteString(",status = '" + strconv.Itoa(int(ship.Status)) + "'")
	s.WriteString(",payment_handle = '" + strconv.Itoa(ship.PaymentHandle) + "'")
	s.WriteString(",payment_status = '" + strconv.Itoa(ship.PaymentStatus) + "'")
	s.WriteString(",user = '" + strconv.Itoa(ship.User) + "'")
	s.WriteString(",contact_name = '" + ship.ContactName + "'")
	s.WriteString(",contact_phone = '" + ship.ContactPhone + "'")
	s.WriteString(",landmark = '" + ship.Landmark + "'")
	s.WriteString(",city = '" + ship.City + "'")
	s.WriteString(",state = '" + ship.State + "'")
	s.WriteString(",country = '" + ship.Country + "'")
	s.WriteString(",zip = '" + strconv.Itoa(ship.Zip) + "'")

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


func (m mysqlRepo) AddressById(ctx context.Context, i int) (*core.Address, error) {
	addr := &core.Address{}

	err := m.db.QueryRowContext(ctx, "SELECT id, user, contact_name, contact_phone, "+
		"landmark, city, state, country, zip, created_at"+
		" FROM "+m.table+"_addresses WHERE id=?", i).
		Scan(&addr.Id, &addr.User, &addr.ContactName, &addr.ContactPhone,
			&addr.Landmark, &addr.City, &addr.State, &addr.Country, &addr.Zip, &addr.CreatedAt)
	if err != nil {
		if err == rowsEmpty {
			return nil, nil
		}
		return addr, err
	}
	return addr, nil
}

func (m mysqlRepo) AllAddresses(ctx context.Context, userId int) (core.Addresses, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT id, user, contact_name, contact_phone, "+
		"landmark, city, state, country, zip, created_at"+
		" FROM "+m.table+"_addresses WHERE user=?", userId)
	if err != nil {
		return nil, err
	}

	var addrs core.Addresses
	for rows.Next() {
		addr := &core.Address{}

		if err = rows.Scan(&addr.Id, &addr.User, &addr.ContactName, &addr.ContactPhone,
			&addr.Landmark, &addr.City, &addr.State, &addr.Country, &addr.Zip, &addr.CreatedAt); err != nil {
			return addrs, err
		}

		addrs = append(addrs, addr)
	}

	return addrs, nil
}

func (m mysqlRepo) EditAddress(ctx context.Context, address *core.Address) error {
	var s []string
	if address.Id == 0 {
		return noIdUpdate
	}
	if address.ContactName != "" {
		s = append(s, "contact_name = '"+address.ContactName+"'")
	}
	if address.ContactPhone != "" {
		s = append(s, "contact_phone = '"+address.ContactPhone+"'")
	}
	if address.Landmark != "" {
		s = append(s, "landmark = '"+address.Landmark+"'")
	}
	if address.City != "" {
		s = append(s, "city = '"+address.City+"'")
	}
	if address.State != "" {
		s = append(s, "state = '"+address.State+"'")
	}
	if address.Country != "" {
		s = append(s, "country = '"+address.Country+"'")
	}
	if address.Zip != 0 {
		s = append(s, "zip = '"+strconv.Itoa(address.Zip)+"'")
	}

	rows, err := m.db.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s_addresses SET %s WHERE id = ?", m.table, strings.Join(s, ",")), address.Id,
	)

	if err != nil {
		return err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return NoRowsAffected
	}

	return err
}

func (m mysqlRepo) UpdateStatus(ctx context.Context, id int, status core.ShippingStatus) error {
	rows, err := m.db.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s SET status = ? WHERE id = ?", m.table), status, id,
	)

	if err != nil {
		return err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return NoRowsAffected
	}

	return err
}

func (m mysqlRepo) UpdatePayment(ctx context.Context, id int, paymentStatus int) error {
	rows, err := m.db.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s SET status = ? WHERE id = ?", m.table), paymentStatus, id,
	)

	if err != nil {
		return err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return NoRowsAffected
	}

	return err
}
