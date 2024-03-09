package structs

type Env struct {
  DBHost string `env:"DB_HOST" envDefault:"127.0.0.1"`
  DBPort int `env:"DB_PORT" envDefault:"27017"`
  DBName string `env:"DB_NAME" envDefault:"turbex"`
  DBUser string `env:"DB_USER" envDefault:""`
  DBPass string `env:"DB_PASS" envDefault:""`
}

