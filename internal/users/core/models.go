package core

type DocumentType string

const (
	Passport   DocumentType = "passport"
	Identifier DocumentType = "identifier"
)

var AllowedDocumenTypes = map[string]DocumentType{
	"identifier": Identifier,
	"passport":   Passport,
}

type User struct {
	ID           string
	Username     string
	Email        string
	FirstName    string
	LastName     string
	Locale       string
	DocumentType DocumentType
	Identifier   string
}

type UserRequest struct {
	Email           string       `json:"email" binding:"required,email"`
	Password        string       `json:"password" binding:"required"`
	PasswordConfirm string       `json:"password_confirm" binding:"required"`
	FirstName       string       `json:"first_name" binding:"required"`
	LastName        string       `json:"last_name" binding:"required"`
	Locale          string       `json:"locale" binding:"required"`
	DocumentType    DocumentType `json:"document_type" binding:"required"`
	Identifier      string       `json:"identifier" binding:"required"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}
