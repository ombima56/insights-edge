package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ombima56/insights-edge/backend/models"
	"github.com/ombima56/insights-edge/backend/service"
	"github.com/ombima56/insights-edge/backend/utils"
)

type InsightHandler struct {
	insightService *service.InsightService
}

func NewInsightHandler(insightService *service.InsightService) *InsightHandler {
	return &InsightHandler{
		insightService: insightService,
	}
}

func (h *InsightHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
	insights, err := h.insightService.GetInsights()
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(insights)
}

func (h *InsightHandler) GetInsight(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.HandleError(w, models.ErrInvalidCredentials)
		return
	}

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.HandleError(w, models.ErrInvalidCredentials)
		return
	}

	insight, err := h.insightService.GetInsight(idInt64)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(insight)
}

func (h *InsightHandler) CreateInsight(w http.ResponseWriter, r *http.Request) {
	var insight models.Insight
	if err := json.NewDecoder(r.Body).Decode(&insight); err != nil {
		utils.HandleError(w, err)
		return
	}

	if err := h.insightService.CreateInsight(&insight); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insight)
}

func (h *InsightHandler) PurchaseInsight(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.HandleError(w, models.ErrInvalidCredentials)
		return
	}

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.HandleError(w, models.ErrInvalidCredentials)
		return
	}

	walletAddr := r.URL.Query().Get("walletAddr")
	if walletAddr == "" {
		utils.HandleError(w, models.ErrInvalidCredentials)
		return
	}

	if err := h.insightService.PurchaseInsight(idInt64, walletAddr); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Purchase successful"})
}

func (h *InsightHandler) GetMyInsights(w http.ResponseWriter, r *http.Request) {
	walletAddr := r.URL.Query().Get("walletAddr")
	if walletAddr == "" {
		utils.HandleError(w, models.ErrInvalidCredentials)
		return
	}

	insights, err := h.insightService.GetMyInsights(walletAddr)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(insights)
}

func (h *InsightHandler) GetMyPurchases(w http.ResponseWriter, r *http.Request) {
	walletAddr := r.URL.Query().Get("walletAddr")
	if walletAddr == "" {
		utils.HandleError(w, models.ErrInvalidCredentials)
		return
	}

	purchases, err := h.insightService.GetMyPurchases(walletAddr)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(purchases)
}
