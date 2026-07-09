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

// GetListHandler godoc
// @Summary      Получить список подписок
// @Description  Возвращает список всех подписок с пагинацией
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        limit  query  int  false  "Количество записей на страницу"  default(10)
// @Param        offset query  int  false  "Смещение"                        default(0)
// @Success      200  {object}  map[string]interface{}  "список подписок"
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions [get]
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

// DeletSubscriptionHandler godoc
// @Summary      Удалить подписку
// @Description  Удаляет подписку по ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "UUID подписки"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeletSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	log.Debug("вызов DeletSubscriptionHandler")
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.sendError(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.service.Delete(ctx, id); err != nil {
		log.Error("failed to delete subscription", zap.Error(err))
		h.sendError(w, "failed to delete subscription", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

// UpdateSubscriptionHandler godoc
// @Summary      Обновить подписку
// @Description  Обновляет существующую подписку по ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id       path   string                           true  "UUID подписки"
// @Param        request  body   dto.UpdateSubscriptionRequest    true  "Данные для обновления"
// @Success      200      {object}  dto.SubscriptionResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	log.Debug("вызов UpdateSubscriptionHandler")
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.sendError(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req service.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "invadlid reques body", http.StatusBadRequest)
		return

	}
	sub, err := h.service.Update(ctx, id, &req)
	if err != nil {
		if err.Error() == "subscription not found" {
			h.sendError(w, "subscription not found0.", http.StatusNotFound)
			return
		}
		log.Error("failed update subscription", zap.Error(err))
		return
	}
	h.sendJson(w, sub, http.StatusOK)

}

// GetTotalPriceHandler godoc
// @Summary      Подсчитать сумму подписок за период
// @Description  Возвращает суммарную стоимость подписок за выбранный период с опциональной фильтрацией
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        start_date    query  string  true  "Дата начала (MM-YYYY)"
// @Param        end_date      query  string  true  "Дата окончания (MM-YYYY)"
// @Param        user_id       query  string  false "ID пользователя"
// @Param        service_name  query  string  false "Название сервиса"
// @Success      200  {object}  map[string]interface{}  "total_price, currency"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/total [get]
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
	if sn := r.URL.Query().Get("service_name"); sn != "" {
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

// GetByIDHandler godoc
// @Summary      Получить подписку по ID
// @Description  Возвращает подписку по её UUID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "UUID подписки"
// @Success      200  {object}  dto.SubscriptionResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /subscriptions/{id} [get]
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

// CreateSubscription godoc
// @Summary      Создать новую подписку
// @Description  Создаёт запись о подписке для пользователя
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSubscriptionRequest true "Данные подписки"
// @Success      201  {object}  dto.SubscriptionResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions [post]
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
			Path:    "/subscriptions",
			Handler: h.CreateSubscription,
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions/{id}",
			Handler: h.GetByIDHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions",
			Handler: h.GetListHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions/total",
			Handler: h.GetTotalPriceHandler,
		},
		{
			Method:  http.MethodPut,
			Path:    "/subscriptions/{id}",
			Handler: h.UpdateSubscriptionHandler,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/subscriptions/{id}",
			Handler: h.DeletSubscriptionHandler,
		},
	}
}
