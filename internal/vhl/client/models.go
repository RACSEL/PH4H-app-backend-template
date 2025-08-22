package client

type CreateQrRequest struct {
	ExpiresOn   string `json:"expiresOn,omitempty"`
	JsonContent string `json:"jsonContent,required"`
	PassCode    string `json:"passCode,omitempty"`
}

type QrData struct {
	Value string
}

type QrValidationRequest struct {
	QRCodeContent string `json:"qrCodeContent,required"`
}

type ValidationResponseStep struct {
	Step        string `json:"step,omitempty"`
	Status      string `json:"status,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Error       string `json:"error,omitempty"`
}

type ValidationResponseUrl struct {
	Url   string `json:"url,required"`
	Flag  string `json:"flag,omitempty"`
	Exp   int    `json:"exp,omitempty"`
	Key   string `json:"key,omitempty"`
	Label string `json:"label,omitempty"`
}

type QRValidationResponse struct {
	Status        map[string]ValidationResponseStep `json:"status"`
	ShLinkContent ValidationResponseUrl             `json:"shLinkContent"`
}

type QrIpsRequest struct {
	Recipient string `json:"recipient,required"`
	PassCode  string `json:"passcode,required"`
}

type VhlManifestResponse struct {
	Files []VhlManifestResponseFile `json:"files,required"`
}

type VhlManifestResponseFile struct {
	ContentType string `json:"contentType,omitempty"`
	Location    string `json:"location,required"`
}
