package vhl

type CreateQrRequest struct {
	ExpiresOn   string `json:"expiresOn,omitempty"`
	JsonContent string `json:"jsonContent,required"`
	PassCode    string `json:"passCode,omitempty"`
}

type QrData struct {
	Value string
}
