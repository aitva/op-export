package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var (
	errListItems = errors.New("fail to list items, try to run \"op list items\"")
)

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

	cssPath := "out.css"
	view := NewView(
		"Password Export",
		ViewConfigAddDate(),
		ViewConfigAddURL(),
		ViewConfigLinkCSS(cssPath),
	)
	{
		file, err := os.Create("out.html")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		err = view.RenderHTML(file, filtered)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	{
		file, err := os.Create(cssPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		err = view.WriteCSS(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
