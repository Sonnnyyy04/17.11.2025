package handler

import (
	"encoding/json"
	"links_project/internal/dto"
	"links_project/internal/pdf"
	"links_project/internal/service"
	"net/http"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler{
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleLinks (w http.ResponseWriter, r *http.Request){
	var req dto.LinksRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	batch, err := h.service.CreateBatch(req.Links)
	if err != nil{
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := dto.LinksResponse{
		LinksStatuses:   batch.Statuses,
		LinksNum: batch.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) HandleReport(w http.ResponseWriter, r *http.Request) {
	var req dto.ReportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	batches, err := h.service.GetBatch(req.LinksList)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	pdfBytes, err := pdf.GeneratePdf(batches)
	if err != nil {
		http.Error(w, "failed to generate pdf", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"report.pdf\"")
	w.Write(pdfBytes)
}