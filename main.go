package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	htmlPath = "out.html"
	cssPath  = "out.css"
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

	view, done, err := createView(cssPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file, err := os.Create(htmlPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	err = view.RenderHTML(file, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("open " + getFileURI(htmlPath) + " in your browser")

	items, err := listItems()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	filtered := items[:0]
	for _, item := range items {
		if item.TemplateUUID != TemplateUUIDLogin {
			continue
		}
		err = getDetails(item)
		if err != nil {
			fmt.Printf("fail to get details for %q\n", item.Overview.Title)
			continue
		}
		filtered = append(filtered, item)
		file.Truncate(0)
		err = view.RenderHTML(file, filtered)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	done()
	file.Truncate(0)
	err = view.RenderHTML(file, filtered)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("wrote %d password into %s\n", len(filtered), htmlPath)
}

func createView(cssPath string) (view *View, done func(), err error) {
	cfg, done := ViewConfigAutoReload()
	view = NewView(
		"1Password Export",
		ViewConfigAddDate(),
		ViewConfigAddURL(),
		ViewConfigLinkCSS(cssPath),
		cfg,
	)

	file, err := os.Create(cssPath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	err = view.WriteCSS(file)
	if err != nil {
		return nil, nil, err
	}
	return
}

func getFileURI(ressource string) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := filepath.Join(dir, ressource)
	return "file://" + path
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
