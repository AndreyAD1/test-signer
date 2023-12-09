package configuration

type ServerConfig struct {
	APISecret       string `env:"API_SECRET,required,notEmpty"`
	DatabaseURL     string `env:"DATABASE_URL,required,notEmpty"`
	Debug           bool   `env:"DEBUG"`
}