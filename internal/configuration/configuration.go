package configuration

type ServerConfig struct {
	APISecret     string `env:"API_SECRET,required,notEmpty"`
	DatabaseURL   string `env:"DATABASE_URL,required,notEmpty"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	SignKey       string `env:"SIGN_KEY,required,notEmpty"`
	Debug         bool   `env:"DEBUG"`
}
