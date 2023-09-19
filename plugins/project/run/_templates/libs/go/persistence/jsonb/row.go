package jsonb

import "github.com/goccy/go-json"

type Row struct {
	id  string
	doc []byte
}

func parse(id string, doc any) (Row, error) {
	b, err := json.Marshal(doc)
	if err != nil {
		return Row{}, err
	}

	return Row{
		id:  id,
		doc: b,
	}, nil
}

func (r Row) decode(doc any) error {
	return json.Unmarshal(r.doc, doc)
}
