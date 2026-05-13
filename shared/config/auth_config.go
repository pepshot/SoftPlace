package config

type AuthConfig struct {
	JWTSecret        string
	JWTLifetimeHours int
}
