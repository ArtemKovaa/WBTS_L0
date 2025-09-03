package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"wbts/internal/domain/entity"
	"wbts/internal/pkg"
)

type OrderRepo struct {
	pgPool *pgxpool.Pool
	cache  map[string]entity.OrderInfo
	mtx sync.RWMutex
}

func NewOrderRepo(pgPool *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{pgPool, make(map[string]entity.OrderInfo), sync.RWMutex{}}
}

func (r *OrderRepo) GetByUID(ctx context.Context, order_uid string) (*entity.OrderInfo, error) {
	startTime := time.Now()
	r.mtx.RLock()
	v, ok := r.cache[order_uid]
	r.mtx.RUnlock()
	if ok {
		log.Printf("Cache hit for order with uid=%s. Fetching time: %s", order_uid, time.Since(startTime))
		return &v, nil
	}

	order, err := r.getOrderByUID(ctx, order_uid)
	if err != nil {
		return nil, err
	}

	payment, err := r.getPaymentByTransaction(ctx, order.PaymentID)
	if err != nil {
		return nil, err
	}

	items, err := r.getItemsByOrderUID(ctx, order.OrderUID)
	if err != nil {
		return nil, err
	}

	orderInfo := entity.OrderInfo{Order: order, Payment: payment, Items: items}
	
	r.mtx.Lock()
	r.cache[order_uid] = orderInfo
	r.mtx.Unlock()

	log.Printf("Order with uid=%s was not found in cache. Fetching time: %s", order_uid, time.Since(startTime))
	return &orderInfo, nil
}

func (r *OrderRepo) Upsert(ctx context.Context, orderInfo entity.OrderInfo) error {
	tx, err := r.pgPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := r.upsertPayment(ctx, tx, orderInfo.Payment); err != nil {
		return err
	}

	if err := r.upsertItems(ctx, tx, orderInfo.Items); err != nil {
		return err
	}

	if err := r.upsertOrder(ctx, tx, orderInfo.Order); err != nil {
		return err
	}

	chrt_ids := make([]int64, len(orderInfo.Items))
    for i, v := range orderInfo.Items {
        chrt_ids[i] = v.ChrtID
    }
	if err := r.insertOrdersItems(ctx, tx, orderInfo.Order.OrderUID, chrt_ids); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *OrderRepo) upsertPayment(ctx context.Context, tx pgx.Tx, payment entity.Payment) error {
	const query = `
        INSERT INTO payments(
			transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (transaction) DO UPDATE SET
            request_id=EXCLUDED.request_id,
            currency=EXCLUDED.currency,
            provider=EXCLUDED.provider,
            amount=EXCLUDED.amount,
            payment_dt=EXCLUDED.payment_dt,
            bank=EXCLUDED.bank,
            delivery_cost=EXCLUDED.delivery_cost,
            goods_total=EXCLUDED.goods_total,
            custom_fee=EXCLUDED.custom_fee
	`
	_, err := tx.Exec(
		ctx,
		query,
		payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount,
		payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepo) upsertItems(ctx context.Context, tx pgx.Tx, items []entity.Item) error {
	const query = `
        INSERT INTO items(
			chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (chrt_id) DO UPDATE SET
            track_number=EXCLUDED.track_number,
            price=EXCLUDED.price,
            rid=EXCLUDED.rid,
            name=EXCLUDED.name,
            sale=EXCLUDED.sale,
            size=EXCLUDED.size,
            total_price=EXCLUDED.total_price,
            nm_id=EXCLUDED.nm_id,
            brand=EXCLUDED.brand,
            status=EXCLUDED.status
	`

	for _, item := range items {
		_, err := tx.Exec(
			ctx,
			query,
			item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, 
			item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderRepo) upsertOrder(ctx context.Context, tx pgx.Tx, order entity.Order) error {
	const query = `
        INSERT INTO orders(
			order_uid, track_number, entry, delivery, payment_id, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        ON CONFLICT (order_uid) DO UPDATE SET
            track_number=EXCLUDED.track_number,
            entry=EXCLUDED.entry,
            delivery=EXCLUDED.delivery,
            payment_id=EXCLUDED.payment_id,
            locale=EXCLUDED.locale,
            internal_signature=EXCLUDED.internal_signature,
            customer_id=EXCLUDED.customer_id,
            delivery_service=EXCLUDED.delivery_service,
            shardkey=EXCLUDED.shardkey,
            sm_id=EXCLUDED.sm_id,
            date_created=EXCLUDED.date_created,
            oof_shard=EXCLUDED.oof_shard
	`
	_, err := tx.Exec(
		ctx,
		query,
		order.OrderUID, order.TrackNumber, order.Entry, order.Delivery, order.PaymentID,
        order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService,
        order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepo) insertOrdersItems(ctx context.Context, tx pgx.Tx, order_uid string, chrt_ids []int64) error {
	const query = "INSERT INTO orders_items(order_uid, chrt_id) VALUES ($1, $2) ON CONFLICT (order_uid, chrt_id) DO NOTHING"

	for _, chrt_id := range chrt_ids {
		if _, err := tx.Exec(ctx, query, order_uid, chrt_id); err != nil {
			return err
		}
	}
	return nil
}

func (r *OrderRepo) getOrderByUID(ctx context.Context, order_uid string) (entity.Order, error) {
	const query = "SELECT * FROM orders WHERE order_uid = $1"

	var order entity.Order 
	err := r.pgPool.QueryRow(ctx, query, order_uid).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Delivery, &order.PaymentID, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID,
		&order.DateCreated, &order.OofShard,
	)
	if err != nil {
		return entity.Order{}, errors.New("Error getting order by UID: " + err.Error())
	}

	return order, nil
}

func (r *OrderRepo) getPaymentByTransaction(ctx context.Context, transaction string) (entity.Payment, error) {
    const query = "SELECT * FROM payments WHERE transaction = $1"

    var payment entity.Payment
    err := r.pgPool.QueryRow(ctx, query, transaction).Scan(
        &payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount,
        &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
    )
    if err != nil {
        return entity.Payment{}, errors.New("Error getting payment by transaction: " + err.Error())
    }

    return payment, nil
}

func (r *OrderRepo) getItemsByOrderUID(ctx context.Context, order_uid string) ([]entity.Item, error) {
    const queryIds = "SELECT chrt_id FROM orders_items WHERE order_uid = $1"

    rows, err := r.pgPool.Query(ctx, queryIds, order_uid)
    if err != nil {
        return nil, errors.New("Error gettings chrt_ids: " + err.Error())
    }
    defer rows.Close()

    var chrtIDs []int64
    for rows.Next() {
        var chrtID int64
        if err := rows.Scan(&chrtID); err != nil {
            return nil, errors.New("Error scanning chrt_ids: " + err.Error())
        }
        chrtIDs = append(chrtIDs, chrtID)
    }
    if err := rows.Err(); err != nil {
        return nil, errors.New("Error gettings chrt_id rows: " + err.Error())
    }

    if len(chrtIDs) == 0 {
        return make([]entity.Item, 0), nil
    }

    args := make([]interface{}, len(chrtIDs))
    for i, id := range chrtIDs {
        args[i] = id
    }

	const itemsQuery = "SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE chrt_id IN (%s)"
    queryItems := fmt.Sprintf(itemsQuery, pkg.GeneratePlaceholders(len(chrtIDs)))

    itemsRows, err := r.pgPool.Query(ctx, queryItems, args...)
    if err != nil {
        return nil, errors.New("Error gettings items: " + err.Error())
    }
    defer itemsRows.Close()

    var items []entity.Item
    for itemsRows.Next() {
        var item entity.Item
		err := itemsRows.Scan(
            &item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
            &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
        )
        if err != nil {
            return nil, errors.New("Error scanning items: " + err.Error())
        }
        items = append(items, item)
    }
    if err := itemsRows.Err(); err != nil {
        return nil, errors.New("Error gettings item rows: " + err.Error())
    }

    return items, nil
}

