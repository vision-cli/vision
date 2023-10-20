package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"

	"github.com/atos-digital/NHSS-scigateway/internal/auth"
	"github.com/atos-digital/NHSS-scigateway/internal/database"
	"github.com/atos-digital/NHSS-scigateway/internal/models"
	"github.com/atos-digital/NHSS-scigateway/internal/models/tables"
	"github.com/atos-digital/NHSS-scigateway/internal/server/templates/components/content"
	"github.com/atos-digital/NHSS-scigateway/internal/server/templates/components/form"
	"github.com/atos-digital/NHSS-scigateway/internal/tmplutils"
)

//go:embed templates
var templates embed.FS
var tmpl = template.Must(tmplutils.ParseFS("index", templates))

//go:embed static
var static embed.FS

func (s *Server) handleServeIndex() http.HandlerFunc {
	type VieData struct {
		Button         form.Button
		CharacterCount form.CharacterCount
		CheckBoxes     form.Checkboxes
		DateInput      form.DateInput
		ErrorSummary   form.ErrorSummary
		FieldSet       form.FieldSet
		Radio          form.Radio
		Select         form.Select
		TextInput      form.TextInput
		TextArea       form.TextArea

		Details   content.Details
		DoDont    content.DoDont
		Expander  content.Expander
		Image     content.Image
		InsetText content.InsetText
	}
	vd := VieData{}
	vd.Button.Text = "Login"

	vd.CharacterCount.Max = 300
	vd.CharacterCount.Message = "Enter your message"
	vd.CharacterCount.HintText.ID = "message-hint"
	vd.CharacterCount.HintText.Text = "This is a hint"
	vd.CharacterCount.Rows = 2

	vd.CheckBoxes.Heading = "What is your favourite colour?"
	vd.CheckBoxes.HintText.ID = "colour-hint"
	vd.CheckBoxes.HintText.Text = "Select all that apply"
	vd.CheckBoxes.Items = []form.CheckboxesItem{
		{ID: "red", Name: "colour", Value: "red", Text: "Red"},
		{ID: "green", Name: "colour", Value: "green", Text: "Green"},
		{ID: "blue", Name: "colour", Value: "blue", Text: "Blue"},
	}

	vd.DateInput.ID = "date"
	vd.DateInput.Heading = "What is your date of birth?"
	vd.DateInput.HintText.ID = "date-hint"
	vd.DateInput.HintText.Text = "For example, 31 3 1980"

	vd.ErrorSummary.Title = "There is a problem"
	vd.ErrorSummary.Error = template.HTML("<h1>ERRORS HERE</h1>")

	vd.FieldSet.Heading = "What is your address?"
	vd.FieldSet.Fields = []form.FieldSetField{
		{ID: "address-line-1", Name: "address-line-1", Label: "Address line 1", HiddenLabel: "Address line 1"},
		{ID: "address-line-2", Name: "address-line-2", HiddenLabel: "Address line 2"},
	}

	vd.Radio.Heading = "What is your favourite colour?"
	vd.Radio.Items = []form.RadioItem{
		{ID: "red", Name: "colour", Value: "red", Label: "Red"},
		{ID: "green", Name: "colour", Value: "green", Label: "Green"},
	}

	vd.Select.ID = "colour"
	vd.Select.Name = "colour"
	vd.Select.Label = "What is your favourite colour?"
	vd.Select.Options = []form.SelectOption{
		{Value: "blue", Text: "Blue"},
		{Value: "red", Text: "Red"},
	}

	vd.TextInput.ID = "name"
	vd.TextInput.Name = "name"
	vd.TextInput.Label = "What is your name?"

	vd.TextArea.ID = "message"
	vd.TextArea.Name = "message"
	vd.TextArea.Label = "What is your message?"
	vd.TextArea.Rows = 2
	vd.TextArea.HintText.ID = "message-hint"
	vd.TextArea.HintText.Text = "This is a hint"

	vd.Details.Text = "This is the text"
	vd.Details.Details = template.HTML("<h1>DETAILS HERE</h1>")

	vd.DoDont.DoLabel = "Do"
	vd.DoDont.DontLabel = "Don't"
	vd.DoDont.DoItems = []content.DoItem{
		{Text: "Do this"},
		{Text: "Do that"},
	}
	vd.DoDont.DontItems = []content.DontItem{
		{Text: "Don't do this"},
		{Text: "Don't do that"},
	}

	vd.Expander.Text = "This is the text"
	vd.Expander.Details = template.HTML("<h1>DETAILS HERE</h1>")

	vd.Image.Src = "https://assets.nhs.uk/prod/images/S_0318_Bullous_pemphigoid_lesions_.2e16d0ba.fill-320x213.jpg"
	vd.Image.Alt = "This is an image"
	vd.Image.Caption = "This is the caption"

	vd.InsetText.Text = template.HTML("<h1>INSET TEXT HERE</h1>")

	return func(w http.ResponseWriter, r *http.Request) {

		err := tmpl.ExecuteTemplate(w, "index", vd)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) handleHelloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := []byte("Hello World!\n")
		// required for compression middleware
		w.Header().Set("Content-Type", http.DetectContentType(resp))
		w.Write(resp)

		user := auth.UserFromContext(r.Context())
		fmt.Fprintln(w, "Hello", user.Username)
	}
}

func (s *Server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.db.Exec("SELECT 'DBD::Pg ping test'").Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Ping!\n"))
	}
}

func (s *Server) handleCreatePatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patient := models.Patient{
			ID:        uuid.New().String(),
			CHI:       "1234567890",
			FirstName: "John",
			LastName:  "Smith",
		}
		j, err := database.Json(patient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p := tables.Patients{ID: patient.ID, Data: j}
		err = s.db.Create(p).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleGetPatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := tables.Patients{}
		err := s.db.First(&p, "id = ?", "3d4b246e-e323-4c85-9876-e5a8d1b1a6ea").Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		patient := models.Patient{}
		err = database.UnJson(p.Data, &patient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(patient)
	}
}
