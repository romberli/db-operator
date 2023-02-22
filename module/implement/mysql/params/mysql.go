package params

const (
	DefaultPrompt  = `[\\u@\\h:\\p][\\d]>`
	DefaultCharSet = "utf8mb4"
)

type MySQL struct {
	Prompt              string `json:"prompt" config:"prompt"`
	DefaultCharacterSet string `json:"default_character_set" config:"default-character-set"`
}

// NewMySQL returns a new *MySQL
func NewMySQL(prompt, defaultCharacterSet string) *MySQL {
	return &MySQL{
		Prompt:              prompt,
		DefaultCharacterSet: defaultCharacterSet,
	}
}

// NewMySQLWithDefault returns a new *MySQL with default values
func NewMySQLWithDefault() *MySQL {
	return &MySQL{
		Prompt:              DefaultPrompt,
		DefaultCharacterSet: DefaultCharSet,
	}
}

// GetPrompt returns the prompt
func (m *MySQL) GetPrompt() string {
	return m.Prompt
}

// GetDefaultCharacterSet returns the default character set
func (m *MySQL) GetDefaultCharacterSet() string {
	return m.DefaultCharacterSet
}
