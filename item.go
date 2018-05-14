package main

import "fmt"

// List of template UUID from 1Password.
const (
	TemplateUUIDLogin      = "001"
	TemplateUUIDSecureNote = "003"
)

// Item contains information on 1Password item.
type Item struct {
	UUID         string       `json:"uuid"`
	TemplateUUID string       `json:"templateUuid"`
	Overview     ItemOverview `json:"overview"`
	Details      *ItemDetails `json:"details,omitempty"`
}

func (item *Item) String() string {
	str := fmt.Sprintf(
		"{uuid: %q, title: %q",
		item.UUID,
		item.Overview.Title,
	)
	if item.Details != nil {
		username, password := item.Details.FindLogin()
		str += fmt.Sprintf(
			", username: %q, password: %q, sectionCount: %d",
			username,
			password,
			len(item.Details.Sections),
		)
	}
	return str + "}"
}

// ItemOverview contains overview information for an item.
type ItemOverview struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// ItemDetails contains detailed information on an Item.
type ItemDetails struct {
	Password string      `json:"password"`
	Fields   []ItemField `json:"fields"`
	Sections []Section   `json:"sections"`
}

// FindLogin finds username and passeword in the fields array.
func (details *ItemDetails) FindLogin() (username, password string) {
	for _, f := range details.Fields {
		switch f.Designation {
		case "username":
			username = f.Value
		case "password":
			password = f.Value
		}
	}
	if details.Password != "" {
		password = details.Password
	}
	return
}

// ItemField contains field information for an Item.
type ItemField struct {
	ID          string `json:"id"`
	Designation string `json:"designation"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value"`
}

// Section contains section information for an Item.
type Section struct {
	Name   string         `json:"name"`
	Title  string         `json:"title"`
	Fields []SectionField `json:"fields"`
}

// SectionField contains field information for a Section.
type SectionField struct {
	Title string `json:"t"`
	Value string `json:"v"`
}
