package configuration

type ServerConfig struct {
	APISecret      string `env:"API_SECRET,required,notEmpty"`
	DatabaseURL    string `env:"DATABASE_URL,required,notEmpty"`
	ServerAddress  string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	PrivateKeyFile string `env:"PRIVATE_KEY,required,notEmpty"`
	PublicKeyFile  string `env:"PUBLIC_KEY,required,notEmpty"`
	Debug          bool   `env:"DEBUG"`
}
