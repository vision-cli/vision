package form

import "html/template"

type Button struct {
	Text string
}

type CharacterCount struct {
	Max      int
	Rows     int
	Message  string
	HintText HintText
}

type HintText struct {
	ID   string
	Text string
}

type Checkboxes struct {
	Heading  string
	HintText HintText
	Items    []CheckboxesItem
}

type CheckboxesItem struct {
	ID    string
	Name  string
	Value string
	Text  string
}

type DateInput struct {
	ID       string
	Heading  string
	HintText HintText
}

type ErrorSummary struct {
	Title string
	Error template.HTML
}

type FieldSet struct {
	Heading string
	Fields  []FieldSetField
}

type FieldSetField struct {
	ID          string
	Name        string
	Label       string
	HiddenLabel string
}

type Radio struct {
	Heading string
	Items   []RadioItem
}

type RadioItem struct {
	ID    string
	Name  string
	Label string
	Value string
}

type Select struct {
	ID      string
	Name    string
	Label   string
	Options []SelectOption
}

type SelectOption struct {
	Value string
	Text  string
}

type TextInput struct {
	ID    string
	Name  string
	Label string
}

type TextArea struct {
	ID       string
	Name     string
	Label    string
	Rows     int
	HintText HintText
}
