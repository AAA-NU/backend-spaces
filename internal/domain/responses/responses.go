package responses

type CoworkingMetaResponse struct {
	ID            string   `json:"id"`
	AvailableTime []string `json:"available_times"`
}
