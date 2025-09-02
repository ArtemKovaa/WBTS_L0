import React, { useState } from "react";

const App = () => {
  const [orderUid, setOrderUid] = useState("");
  const [orderData, setOrderData] = useState(null);
  const [error, setError] = useState("");

  const fetchOrder = async () => {
    setError("");
    setOrderData(null);
    if (!orderUid.trim()) {
      setError("Пожалуйста, введите номер заказа");
      return;
    }
    try {
      const response = await fetch(`/order/${orderUid}`);
      if (!response.ok) {
        throw new Error(`Ошибка сервера: ${response.status}`);
      }
      const data = await response.json();
      setOrderData(data);
    } catch (err) {
      setError(`Не удалось получить данные: ${err.message}`);
    }
  };

  return (
    <div style={styles.container}>
      <h1 style={styles.title}>Поиск заказа</h1>
      <div style={styles.inputRow}>
        <input
          type="text"
          placeholder="Введите номер заказа (order_uid)"
          value={orderUid}
          onChange={(e) => setOrderUid(e.target.value)}
          style={styles.input}
        />
        <button onClick={fetchOrder} style={styles.button}>
          Найти
        </button>
      </div>
      {error && <div style={styles.error}>{error}</div>}
      {orderData && (
        <div style={styles.orderContainer}>
          <h2 style={styles.sectionTitle}>Информация о заказе</h2>
          <div style={styles.row}>
            <div style={styles.label}>Номер заказа:</div>
            <div style={styles.value}>{orderData.order_uid}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Трек-номер:</div>
            <div style={styles.value}>{orderData.track_number}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Дата создания:</div>
            <div style={styles.value}>{new Date(orderData.date_created).toLocaleString()}</div>
          </div>

          <h2 style={styles.sectionTitle}>Данные доставки</h2>
          <div style={styles.row}>
            <div style={styles.label}>Получатель:</div>
            <div style={styles.value}>{orderData.delivery.name}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Телефон:</div>
            <div style={styles.value}>{orderData.delivery.phone}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Город:</div>
            <div style={styles.value}>{orderData.delivery.city}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Адрес:</div>
            <div style={styles.value}>{orderData.delivery.address}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Регион / Область:</div>
            <div style={styles.value}>{orderData.delivery.region}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Почтовый индекс:</div>
            <div style={styles.value}>{orderData.delivery.zip}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>E-mail:</div>
            <div style={styles.value}>{orderData.delivery.email}</div>
          </div>

          <h2 style={styles.sectionTitle}>Информация по оплате</h2>
          <div style={styles.row}>
            <div style={styles.label}>Номер транзакции:</div>
            <div style={styles.value}>{orderData.payment.transaction}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Сумма оплаты:</div>
            <div style={styles.value}>
              {orderData.payment.amount} {orderData.payment.currency}
            </div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Поставщик платежа:</div>
            <div style={styles.value}>{orderData.payment.provider}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Стоимость доставки:</div>
            <div style={styles.value}>{orderData.payment.delivery_cost} {orderData.payment.currency}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Общая стоимость товаров:</div>
            <div style={styles.value}>{orderData.payment.goods_total} {orderData.payment.currency}</div>
          </div>
          <div style={styles.row}>
            <div style={styles.label}>Банк:</div>
            <div style={styles.value}>{orderData.payment.bank}</div>
          </div>

          <h2 style={styles.sectionTitle}>Состав заказа</h2>
          <table style={styles.table}>
            <thead>
              <tr>
                <th style={{ textAlign: "left" }}>Название</th>
                <th style={{ textAlign: "left" }}>Бренд</th>
                <th style={{ textAlign: "left" }}>Цена (за шт.)</th>
                <th style={{ textAlign: "left" }}>Скидка (%)</th>
                <th style={{ textAlign: "left" }}>Итоговая цена</th>
              </tr>
            </thead>
            <tbody>
              {orderData.items.map((item) => (
                <tr key={item.chrt_id}>
                  <td>{item.name}</td>
                  <td>{item.brand}</td>
                  <td>{item.price} {orderData.payment.currency}</td>
                  <td>{item.sale}</td>
                  <td>{item.total_price} {orderData.payment.currency}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

const styles = {
  container: {
    maxWidth: 900,
    margin: "30px auto",
    padding: 20,
    fontFamily: "'Helvetica Neue', Arial, sans-serif",
    backgroundColor: "#ffffff",
    borderRadius: 10,
    boxShadow: "0 4px 12px rgba(0,0,0,0.1)",
    color: "#333",
  },
  title: {
    textAlign: "center",
    marginBottom: 25,
    fontSize: 28,
    fontWeight: "bold",
    color: "#2c3e50",
  },
  inputRow: {
    display: "flex",
    marginBottom: 20,
    justifyContent: "center",
  },
  input: {
    width: "350px",
    padding: 12,
    fontSize: 16,
    borderRadius: 6,
    border: "1px solid #ccc",
    marginRight: 12,
    boxSizing: "border-box",
  },
  button: {
    padding: "12px 28px",
    fontSize: 16,
    borderRadius: 6,
    border: "none",
    backgroundColor: "#2980b9",
    color: "#fff",
    cursor: "pointer",
    transition: "background-color 0.3s ease",
  },
  buttonHover: {
    backgroundColor: "#1f618d",
  },
  error: {
    color: "#e74c3c",
    marginBottom: 20,
    textAlign: "center",
    fontWeight: "600",
  },
  orderContainer: {
    backgroundColor: "#f9f9f9",
    padding: 25,
    borderRadius: 8,
    boxShadow: "0 2px 8px rgba(0,0,0,0.05)",
  },
  sectionTitle: {
    borderBottom: "2px solid #2980b9",
    paddingBottom: 6,
    marginBottom: 16,
    color: "#2980b9",
    fontWeight: "700",
    fontSize: 22,
  },
  row: {
    display: "flex",
    padding: "6px 0",
    borderBottom: "1px solid #ddd",
  },
  label: {
    width: 180,
    fontWeight: "600",
    color: "#555",
  },
  value: {
    flexGrow: 1,
    color: "#222",
  },
  table: {
    width: "100%",
    borderCollapse: "collapse",
    marginTop: 10,
  },
  "table th, table td": {
    border: "1px solid #bbb",
    padding: "10px 15px",
    textAlign: "left",
  },
  "table th": {
    backgroundColor: "#2980b9",
    color: "white",
  }
};

export default App;
