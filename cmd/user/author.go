package user

type Author struct {
	Type string `json:"type"          mapstructure:"type"`
	Raw  string `json:"raw,omitempty" mapstructure:"raw"`
	User User   `json:"user"          mapstructure:"user"`
}

// IsEmpty checks if this Author is empty
func (author Author) IsEmpty() bool {
	return author.Type == "" && author.Raw == "" && author.User.IsEmpty()
}
