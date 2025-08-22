package core

type Bundle struct {
	ID           string                 `json:"id"`
	Identifier   map[string]interface{} `json:"identifier,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
	ResourceType string                 `json:"resourceType"`
	Signature    map[string]interface{} `json:"signature,omitempty"`
	Timestamp    string                 `json:"timestamp"`
	Type         string                 `json:"type"`
	Entry        []Entry                `json:"entry,omitempty"`
}

type Entry struct {
	FullURL  string                 `json:"fullUrl"`
	Resource map[string]interface{} `json:"resource"` // This could be any FHIR resource, it will be treated as a map
}

type Composition struct {
	URL             string                   `json:"url,omitempty"` // This will be here to then convert to Bundle easily
	ID              string                   `json:"id"`
	ResourceType    string                   `json:"resourceType"`
	Text            map[string]interface{}   `json:"text,omitempty"`
	Meta            map[string]interface{}   `json:"meta,omitempty"`
	Status          string                   `json:"status,omitempty"`
	Subject         map[string]interface{}   `json:"subject,omitempty"`
	Code            map[string]interface{}   `json:"code,omitempty"`
	Type            CodeableConcept          `json:"type,omitempty"`
	Author          []map[string]interface{} `json:"author,omitempty"`
	Confidentiality string                   `json:"confidentiality,omitempty"`
	Custodian       map[string]interface{}   `json:"custodian,omitempty"`
	Date            string                   `json:"date,omitempty"`
	Section         []Section                `json:"section,omitempty"`
	Title           string                   `json:"title,omitempty"`
}

type CodeableConcept struct {
	Coding []map[string]interface{} `json:"coding,omitempty"`
}

type Section struct {
	Title string                   `json:"title,omitempty"`
	Code  CodeableConcept          `json:"code,omitempty"`
	Entry []map[string]interface{} `json:"entry,omitempty"`
}

type Reference struct {
	Reference string `json:"reference"`
}
