package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"wbts/internal/domain/entity"
)

type OrderRepo struct {
	pgPool *pgxpool.Pool
	cache  map[string]entity.OrderInfo
}

func NewOrderRepo(pgPool *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{pgPool, make(map[string]entity.OrderInfo)}
}

func (or *OrderRepo) Upsert(ctx context.Context, orderInfo entity.OrderInfo) error {
	tx, err := or.pgPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := or.upsertPayment(ctx, tx, orderInfo.Payment); err != nil {
		log.Println("payment")
		return err
	}

	if err := or.upsertItems(ctx, tx, orderInfo.Items); err != nil {
		log.Println("items")
		return err
	}

	if err := or.upsertOrder(ctx, tx, orderInfo.Order); err != nil {
		log.Println("order")
		return err
	}

	return tx.Commit(ctx)
}

func (or *OrderRepo) upsertPayment(ctx context.Context, tx pgx.Tx, payment entity.Payment) error {
	query := `
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

func (or *OrderRepo) upsertItems(ctx context.Context, tx pgx.Tx, items []entity.Item) error {
	query := `
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

func (or *OrderRepo) upsertOrder(ctx context.Context, tx pgx.Tx, order entity.Order) error {
	query := `
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
