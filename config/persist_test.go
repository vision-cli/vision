package config

type MockPersist struct {
	values map[string]string
}

func NewMockPersist() *MockPersist {
	return &MockPersist{
		values: make(map[string]string),
	}
}

func (m *MockPersist) Set(key, value string) {
	m.values[key] = value
}

func (m *MockPersist) GetString(key string) string {
	return m.values[key]
}

func (m *MockPersist) WriteConfig() error {
	return nil
}

func (m *MockPersist) WriteConfigAs(filename string) error {
	return nil
}

func (m *MockPersist) SetConfigName(name string) {
}

func (m *MockPersist) AddConfigPath(path string) {
}

func (m *MockPersist) ReadInConfig() error {
	return nil
}

func (m *MockPersist) Values() map[string]string {
	return m.values
}
