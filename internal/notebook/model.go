package notebook

type Customer struct {
	ID                int32               `json:"id, omitempty"`
	Birthdate         string              `json:"birthDate, omitempty"`
	Sex               string              `json:"sex, omitempty"` //есть справочник sex
	CustomerDocuments []CustomerDocuments `json:"customerDocuments, omitempty"`
	Phone             string              `json:"phone, omitempty "`
	Email             string              `json:"email, omitempty"`
}

type CustomerDocuments struct {
	ID               int64  `json:"id" validate:"required"`
	Type             string `json:"type" validate:"required"` //есть справочник DocumentType
	Number           string `json:"number" validate:"required"`
	Firstname        string `json:"firstName" validate:"required"`
	Middlename       string `json:"middleName"`
	Lastname         string `json:"lastName" validate:"required"`
	Citizenship      string `json:"citizenship" validate:"required"`
	ExpireDate       string `json:"expireDate, omitempty"`
	IssueDate        string `json:"issueDate, omitempty"`
	IssuingAuthority string `json:"issuingAuthority, omitempty"`
	IsActive         bool   `json:"isActive"`
}
