package domain

type ImageRequest struct {
	ImageBase64 string `json:"imageBase64"`
}

type ExtractedData struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Address      string `json:"address"`
	Mobile       string `json:"mobile"`
}

type APIResponse struct {
	Success bool          `json:"success"`
	Data    ExtractedData `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
	Message string        `json:"message,omitempty"`
}
