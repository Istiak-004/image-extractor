package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/istiak-004/image-extractor/internals/domain"
	"github.com/istiak-004/image-extractor/internals/service"
)

type Handler struct {
	extractorService  service.ExtractorService
	pngCreatorService service.PNGCreatorService
}

func NewHandler(es service.ExtractorService, cs service.PNGCreatorService) *Handler {
	return &Handler{
		extractorService:  es,
		pngCreatorService: cs,
	}
}

var request struct {
	ImageBase64 string `json:"imageBase64"`
}

func (h *Handler) PNGCreatorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside handler------------------")
	w.Header().Set("Content-Type", "application/json")

	var requestedData domain.ExtractedData
	jsonMapData := make(map[string]string)

	if err := json.NewDecoder(r.Body).Decode(&requestedData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	jsonMapData["name"] = requestedData.Name
	jsonMapData["organization"] = requestedData.Organization
	jsonMapData["address"] = requestedData.Address
	jsonMapData["mobile"] = requestedData.Mobile

	pngData, err := h.pngCreatorService.CreatePNGFromJson(jsonMapData)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating png from json %s", err))
		return
	}

	err = os.WriteFile("with_json.png", pngData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Convert to base64 image
	base64Str := base64.StdEncoding.EncodeToString(pngData)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": base64Str,
	})

}

func (h *Handler) ExtractDataFromImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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
