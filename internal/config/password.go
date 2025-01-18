package config

type PasswordConfig struct {
	Length        int    `json:"length"`
	UseUppercase  bool   `json:"useUppercase"`
	UseLowercase  bool   `json:"useLowercase"`
	UseNumbers    bool   `json:"useNumbers"`
	UseSymbols    bool   `json:"useSymbols"`
	CustomSymbols string `json:"customSymbols"`
}

const (
	Uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lowercase = "abcdefghijklmnopqrstuvwxyz"
	Numbers   = "0123456789"
	Symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)
