package app

import (
	"encoding/json"
	"net/http"

	"github.com/istiak-004/image-extractor/internals/service"
)

type Handler struct {
	extractorService service.ExtractorService
}

func NewHandler(es service.ExtractorService) *Handler {
	return &Handler{
		extractorService: es,
	}
}

var request struct {
	ImageBase64 string `json:"imageBase64"`
}

func (h *Handler) ExtractDataFromImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if request.ImageBase64 == "" {
		respondWithError(w, http.StatusBadRequest, "No image provided")
		return
	}

	extractedData, err := h.extractorService.ExtractFromBase64(request.ImageBase64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    extractedData,
		"message": "Successfully extracted JSON from image",
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
