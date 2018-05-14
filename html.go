package main

import (
	"html/template"
	"io"
	"strings"
	"time"
)

const cssTmpl = `.main {
	margin: 0 auto;
	max-width: 768px;
}
.header {
	display: flex;
	flex-direction: row;
	flex-wrap: wrap; 
	align-items: center;
}
.header__date {
	margin-left: auto;
}
.item {
	-webkit-column-break-inside: avoid;
	page-break-inside: avoid;
	break-inside: avoid;
}
.dl {
	display: grid;
	grid-template-columns: auto auto;
	justify-content: start;
	margin-left: 2em;
}
.dl__dd {
	margin-left: 1em;
}`

const htmlTmpl = `<!doctype html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		{{ with .CSSPath -}}
		<link rel="stylesheet" href="{{.}}">
		{{- end }}
		{{ with .CSS -}}
		<style>
			{{ . }}
		</style>
		{{- end }}
		<title>{{ .Title }}</title>
	</head>
	<body>
		<main class="main">
			<header class="header">
				<h1 class="header__title">{{ .Title }}</h1>
				{{ with .Date -}}
				<p class="header__date">{{ . }}</p>
				{{- end }}
			</header>
			{{ range .Items }}
			<article class="item">
				<h2 class="item__title">{{ .Title }}</h2>
				<dl class="dl">
					{{ if $.ShowURL -}}
					<dt class="dl__dt">url:</dt><dd class="dl__dd">{{ .URL }}</dd>
					{{- end }}
					<dt class="dl__dt">username:</dt><dd class="dl__dd">{{ .Username }}</dd>
					<dt class="dl__dt">password:</dt><dd class="dl__dd">{{ .Password }}</dd>
				</dl>
				{{ range .Sections }}
				<h3 class="item__title">{{ .Title }}</h3>
				<dl class="dl">
					{{- range .Fields -}}
					<dt class="dl__dt">{{ .Name }}:</dt><dd class="dl__dd">{{ .Value }}</dd>
					{{- end }}
				</dl>
				{{ end }}
			</article>
			{{ end }}
		</main>
	</body>
</html>`

// ViewConfig is a configuration function for the View.
type ViewConfig func(*View)

// ViewConfigAddDate adds a date to the view.
func ViewConfigAddDate() ViewConfig {
	return func(v *View) {
		v.Date = time.Now().Format("2006/01/02 15:04:05")
	}
}

// ViewConfigAddURL adds an URL to the item.
func ViewConfigAddURL() ViewConfig {
	return func(v *View) {
		v.ShowURL = true
	}
}

// ViewConfigInlineCSS inlines CSS inside the HTML.
func ViewConfigInlineCSS() ViewConfig {
	return func(v *View) {
		v.CSS = cssTmpl
	}
}

// ViewConfigLinkCSS link a stylesheet with the HTML.
func ViewConfigLinkCSS(path string) ViewConfig {
	return func(v *View) {
		v.CSSPath = path
	}
}

// ViewItem contains information to display to the user.
type ViewItem struct {
	Title    string
	Username string
	Password string
	URL      string
	Sections []ViewSection
}

// ViewSection contains Section information to display to the user.
type ViewSection struct {
	Title  string
	Fields []ViewSectionField
}

// NewViewSection instanciates a new ViewSection.
func NewViewSection(section Section) ViewSection {
	v := ViewSection{
		Title:  section.Title,
		Fields: make([]ViewSectionField, len(section.Fields)),
	}
	for k, field := range section.Fields {
		v.Fields[k] = ViewSectionField{
			Name:  field.Title,
			Value: field.Value,
		}
	}
	return v
}

// ViewSectionField contains field information for the section.
type ViewSectionField struct {
	Name  string
	Value string
}

// View represents the data used to render data to the user.
type View struct {
	Title   string
	Date    string
	ShowURL bool
	CSS     template.CSS
	CSSPath string
	Items   []ViewItem
	tmpl    *template.Template
}

// NewView instanciates a new view.
func NewView(title string, opts ...ViewConfig) *View {
	v := &View{
		Title: title,
		tmpl:  template.Must(template.New("htmlTmpl").Parse(htmlTmpl)),
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// RenderHTML writes a list of Item to w.
func (v *View) RenderHTML(w io.Writer, items []*Item) error {
	v.Items = make([]ViewItem, len(items))
	for i, item := range items {
		v.Items[i] = ViewItem{
			Title: item.Overview.Title,
			URL:   item.Overview.URL,
		}
		if item.Details == nil {
			continue
		}

		username, password := item.Details.FindLogin()
		v.Items[i].Username = username
		v.Items[i].Password = password

		sections := make([]ViewSection, len(item.Details.Sections))
		for j, section := range item.Details.Sections {
			sections[j] = NewViewSection(section)
		}
		v.Items[i].Sections = sections
	}
	return v.tmpl.Execute(w, v)
}

// WriteCSS writes CSS style into w.
func (v *View) WriteCSS(w io.Writer) error {
	r := strings.NewReader(cssTmpl)
	_, err := io.Copy(w, r)
	return err
}
