package entity

// Response struct which contains an API Response
type Response struct {
	Message     string                 `json:"message,omitempty"`
	Validations []string               `json:"validations,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Successful  bool                   `json:"successful"`
}
