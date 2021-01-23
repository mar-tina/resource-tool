package config

type Config struct {
	Token    string `arg:"env:POSTMAN_TOKEN"`
	URI      string `arg:"env:URI"`
	PORT     string `arg:"env:PORT"`
	USERNAME string `arg:"env:USERNAME"`
	PASSWORD string `arg:"env:PASSWORD"`
}

func DefaultConfig() *Config {
	return &Config{
		Token:    "",
		URI:      "https://api.getpostman.com",
		PORT:     "9999",
		USERNAME: "mar-tina",
		PASSWORD: "password",
	}
}
