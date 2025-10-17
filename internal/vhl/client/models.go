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

type ICVPQrValidationRequest struct {
	IncludeRaw bool   `json:"include_raw,required"`
	QRData     string `json:"qr_data,required"`
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

// ----------------
// ICVP Validation response strucs
type ICVPQRValidationResponse struct {
	COSE        COSEData        `json:"cose"`
	Diagnostics DiagnosticsData `json:"diagnostics"`
	HCERT       interface{}     `json:"hcert"` // Can be null, use interface{}
	Payload     PayloadData     `json:"payload"`
}

type COSEData struct {
	Raw         RawData                `json:"_raw"`
	KidB64      string                 `json:"kid_b64"`
	KidHex      string                 `json:"kid_hex"`
	Protected   ProtectedData          `json:"protected"`
	Signature   string                 `json:"signature"`
	Unprotected map[string]interface{} `json:"unprotected"` // Empty object, use map[string]interface{}
}

type RawData struct {
	PayloadBstr   string `json:"payload_bstr"`
	ProtectedBstr string `json:"protected_bstr"`
	Signature     string `json:"signature"`
}

type ProtectedData struct {
	Key1 int      `json:"1"` // The key is the number '1'
	Key4 Key4Data `json:"4"` // The key is the number '4'
}

type Key4Data struct {
	B64 string `json:"_b64"`
}

type DiagnosticsData struct {
	Base45DecodedLen    int `json:"base45_decoded_len"`
	ZlibDecompressedLen int `json:"zlib_decompressed_len"`
}

type PayloadData struct {
	Key260 Payload260Data `json:"-260"` // The key is the number '-260'
	Key1   string         `json:"1"`    // The key is the number '1'
	Key6   int            `json:"6"`    // The key is the number '6'
}

type Payload260Data struct {
	Key6 InnerPayloadData `json:"-6"` // The key is the number '-6'
}

type InnerPayloadData struct {
	DOB string          `json:"dob"`
	N   string          `json:"n"`
	NDT string          `json:"ndt"`
	NID string          `json:"nid"`
	NTL string          `json:"ntl"`
	S   string          `json:"s"`
	V   VaccinationData `json:"v"`
}

type VaccinationData struct {
	BO  string `json:"bo"`
	CN  string `json:"cn"`
	DT  string `json:"dt"`
	IS  string `json:"is"`
	VLE string `json:"vle"`
	VLS string `json:"vls"`
	VP  string `json:"vp"`
}

//--------------------
