package mocks

import "io/fs"

type MockTmplWriter struct {
	history  []string
	existing map[string]string
}

func (m *MockTmplWriter) WriteTemplatedFS(templatePath string, targetPath string, templateFiles fs.FS, p interface{}) error {
	m.history = append(m.history, "WriteTemplatedFS: "+targetPath)
	return nil
}

func (m *MockTmplWriter) WriteExactFS(templatePath string, targetPath string, templateFiles fs.FS) error {
	m.history = append(m.history, "WriteExactFS: "+targetPath)
	return nil
}

func (m *MockTmplWriter) CreateDir(path string) error {
	m.history = append(m.history, "CreateDir: "+path)
	return nil
}

func (m *MockTmplWriter) Exists(path string) bool {
	_, exists := m.existing[path]
	return exists
}

func (m *MockTmplWriter) History() []string {
	return m.history
}

func (m *MockTmplWriter) AddExists(filename string) {
	m.existing[filename] = filename
}

func NewMockTmplWriter() MockTmplWriter {
	return MockTmplWriter{
		history:  []string{},
		existing: map[string]string{},
	}
}
