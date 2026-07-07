package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/evgen6501-star/golang-subscriptions/internal/logger"
	"github.com/evgen6501-star/golang-subscriptions/internal/server"
	"github.com/evgen6501-star/golang-subscriptions/internal/service"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
	logger  *logger.Logger
}

func NewSubscriptionHandler(svc service.SubscriptionService, log *logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: svc,
		logger:  log,
	}
}
func (h *SubscriptionHandler) GetListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	log.Info("вызов GetListHandler")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	subs, err := h.service.List(ctx, limit, offset)
	if err != nil {
		log.Error("failed list subscriptions", zap.Error(err))
		h.sendError(w, "internal error", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"subscriptions": subs,
		"count":         len(subs),
		"limit":         limit,
		"offset":        offset,
	}
	h.sendJson(w, response, http.StatusOK)
}
func (h *SubscriptionHandler) GetTotalPriceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	log.Debug("вызов GetTotalPriceHandler")
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")
	log.Debug("startStr", zap.String("startStr", startStr))
	log.Debug("endStr", zap.String("endStr", endStr))
	if startStr == "" || endStr == "" {

		h.sendError(w, "invalid start_date format (MM-YYYY)", http.StatusBadRequest)
		return
	}
	start, err := time.Parse("01-2006", startStr)
	if err != nil {
		h.sendError(w, "invalid start_date format (MM-YYYY)", http.StatusBadRequest)
		return
	}
	end, err := time.Parse("01-2006", endStr)
	if err != nil {
		h.sendError(w, "invalid end_date format (MM-YYYY)", http.StatusBadRequest)
		return
	}
	end = end.AddDate(0, 1, -1)
	var userID *uuid.UUID
	if uidStr := r.URL.Query().Get("user_id"); uidStr != "" {
		uid, err := uuid.Parse(uidStr)
		if err != nil {
			h.sendError(w, "invalid user_id format ", http.StatusBadRequest)
			return
		}
		userID = &uid

	}
	var serviceName *string
	if sn := r.URL.Query().Get("sevice_name"); sn != "" {
		serviceName = &sn
	}
	req := service.TotalPriceRequest{
		UserID:      userID,
		ServiceName: serviceName,
		StartDate:   start,
		EndDate:     end,
	}
	total, err := h.service.GetTotalPrice(ctx, req)
	if err != nil {
		log.Error("failed to get total price", zap.Error(err))
		h.sendError(w, "failed to get total price", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"total_price": total,
		"currency":    "RUB",
	}
	h.sendJson(w, response, http.StatusOK)
}

func (h *SubscriptionHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	log.Debug("вызов GetByIDHandler")
	idStr := r.PathValue("id")
	// log.Debug("idStr value", zap.String("idStr", idStr))
	log.Debug("idStr value", zap.String("idStr", idStr)) // <-- добавьте эту строку
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error("Ошибка парсинга id")
		h.sendError(w, "error id", http.StatusBadRequest)
		return

	}
	sub, err := h.service.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to Get id", zap.Error(err))
		h.sendError(w, "subscription niot found", http.StatusNotFound)
		return
	}
	if sub == nil {
		h.sendError(w, "subscription niot found", http.StatusNotFound)
		return
	}
	h.sendJson(w, sub, http.StatusOK)

}
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	log.Debug("Вызов CreateSubscriptionHandler")
	var req struct {
		ServiceName string `json:"service_name"`
		Price       int    `json:"price"`
		UserID      string `json:"user_id"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		h.sendError(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		h.sendError(w, "invalid start date fomat", http.StatusBadRequest)
		return
	}
	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			h.sendError(w, "invalid end date fomat", http.StatusBadRequest)
			return
		}
		endDate = &parsed
	}
	// endDate, err := time.Parse("01-2006", req.EndDate)
	// if err != nil {

	// 	h.sendError(w, "invalid end date fomat", http.StatusBadRequest)
	// 	return SeviceName
	// }ServiceName
	createReq := &service.CreateSubscriptionRequest{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	subscription, err := h.service.Create(r.Context(), createReq)
	if err != nil {
		h.logger.Error("error create subscribe")
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.sendJson(w, subscription, http.StatusCreated)

}

func (h *SubscriptionHandler) sendError(w http.ResponseWriter, message string, status int) {
	h.logger.Warn("Send error")
	h.sendJson(w, map[string]string{"error": message}, status)
}

func (h *SubscriptionHandler) sendJson(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

//	func () Create() {
//		//...
//	}
func (h *SubscriptionHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/subscribetions",
			Handler: h.CreateSubscription,
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscribetions/{id}",
			Handler: h.GetByIDHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscribetions",
			Handler: h.GetListHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscribetions/total",
			Handler: h.GetTotalPriceHandler,
		},
	}
}
