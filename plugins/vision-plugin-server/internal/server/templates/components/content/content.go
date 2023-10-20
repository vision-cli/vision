package content

import "html/template"

type Details struct {
	Text    string
	Details template.HTML
}

type DoDont struct {
	DoLabel   string
	DontLabel string
	DoItems   []DoItem
	DontItems []DontItem
}

type DoItem struct {
	Text string
}

type DontItem struct {
	Text string
}

type Expander struct {
	Text    string
	Details template.HTML
}

type Image struct {
	Src     string
	Alt     string
	Caption string
}

type InsetText struct {
	Text       template.HTML
	HiddenText string
}
