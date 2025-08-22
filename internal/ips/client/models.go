package client

type Bundle struct {
	ResourceType string        `json:"resourceType"`
	ID           string        `json:"id,omitempty"`
	Meta         *Meta         `json:"meta,omitempty"`
	Identifier   *Identifier   `json:"identifier,omitempty"`
	Type         string        `json:"type"`
	Timestamp    string        `json:"timestamp,omitempty"`
	Total        int           `json:"total,omitempty"`
	Link         []BundleLink  `json:"link,omitempty"`
	Entry        []BundleEntry `json:"entry,omitempty"`
	Signature    *Signature    `json:"signature,omitempty"`
}

type BundleLink struct {
	Relation string `json:"relation"`
	URL      string `json:"url"`
}

type BundleEntry struct {
	FullURL  string         `json:"fullUrl"`
	Resource *EntryResource `json:"resource,omitempty"`
	Search   *BundleSearch  `json:"search,omitempty"`
}

type BundleSearch struct {
	Mode string `json:"mode"`
}

type EntryResource struct {
	ResourceType     string            `json:"resourceType"`
	ID               string            `json:"id,omitempty"`
	Meta             *Meta             `json:"meta,omitempty"`
	Text             *ResourceText     `json:"text,omitempty"`
	MasterIdentifier *Identifier       `json:"masterIdentifier,omitempty"`
	Identifier       []Identifier      `json:"identifier,omitempty"`
	Status           string            `json:"status"`
	Type             interface{}       `json:"type,omitempty"` // This is an any type, can be a string or an object
	Subject          *Reference        `json:"subject"`
	Author           []Reference       `json:"author,omitempty"`
	Title            string            `json:"title,omitempty"`
	Date             string            `json:"date,omitempty"`
	Confidentiality  string            `json:"confidentiality,omitempty"`
	Custodian        *Reference        `json:"custodian,omitempty"`
	Content          []DocumentContent `json:"content,omitempty"`
	Section          []EntrySection    `json:"section,omitempty"`
}

type ResourceText struct {
	Status string `json:"status"`
	Div    string `json:"div,omitempty"`
}

type EntryType struct {
	Coding []Coding `json:"coding"`
}

type Coding struct {
	System  string `json:"system"`
	Code    string `json:"code"`
	Display string `json:"display,omitempty"`
}

type EntrySection struct {
	Title string        `json:"title,omitempty"`
	Code  *EntryType    `json:"code,omitempty"`
	Text  *ResourceText `json:"text,omitempty"`
	Entry []Reference   `json:"entry,omitempty"`
}

type Meta struct {
	VersionID   string   `json:"versionId,omitempty"`
	LastUpdated string   `json:"lastUpdated,omitempty"`
	Source      string   `json:"source,omitempty"`
	Profile     []string `json:"profile,omitempty"`
}

type Identifier struct {
	System string     `json:"system,omitempty"`
	Value  string     `json:"value,omitempty"`
	Use    string     `json:"use,omitempty"`
	Type   *EntryType `json:"type,omitempty"`
}

type Reference struct {
	Reference string `json:"reference"`
	Type      string `json:"type,omitempty"`
	Display   string `json:"display,omitempty"`
}

type DocumentContent struct {
	Attachment Attachment `json:"attachment"`
}

type Attachment struct {
	ContentType string `json:"contentType"`
	URL         string `json:"url"`
}

type Signature struct {
	Type []SignatureType `json:"type"`
	When string          `json:"when"`
	Who  *SignatureWho   `json:"who"`
	Data string          `json:"data"`
}

type SignatureType struct {
	System string `json:"system"`
	Code   string `json:"code"`
}

type SignatureWho struct {
	Identifier Identifier `json:"identifier"`
}
