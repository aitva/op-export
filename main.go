package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// List some template UUID from 1Password.
const (
	TemplateUUIDLogin      = "001"
	TemplateUUIDSecureNote = "003"
)

var (
	errListItems = errors.New("fail to list items, try to run \"op list items\"")
)

// Item contains information on a 1Password item.
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
			"username: %q, password: %q, sectionCount: %d",
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

// ItemField contains detailed information on a form field.
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

func main() {
	_, err := exec.LookPath("op")
	if err != nil {
		fmt.Println("\"op\" is not in the PATH")
		os.Exit(1)
	}

	items, err := listItems()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filtered := items[:]
	for _, item := range items {
		if item.TemplateUUID != TemplateUUIDLogin {
			continue
		}
		err = getDetails(item)
		if err != nil {
			fmt.Printf("fail to get details for %q\n", item.Overview.Title)
			continue
		}
		fmt.Println(item)
		filtered = append(filtered, item)
	}
}

func listItems() (items []*Item, err error) {
	cmd := exec.Command("op", "list", "items")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	err = cmd.Start()
	if err != nil {
		return nil, errListItems
	}

	err = json.NewDecoder(stdout).Decode(&items)
	if err != nil {
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}
	return
}

func getDetails(item *Item) error {
	cmd := exec.Command("op", "get", "item", item.UUID)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return errListItems
	}

	err = json.NewDecoder(stdout).Decode(item)
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
