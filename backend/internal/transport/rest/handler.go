package rest

import (
	"net/http"
	"encoding/json"

	"wbts/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService}
}


func (h *OrderHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
    order_uid := r.PathValue("order_uid")

	order, err := h.orderService.Get(order_uid)
	if err != nil {
		http.Error(w, "Error getting order by uid: " + err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.MarshalIndent(order, "", "    ")
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(body)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

