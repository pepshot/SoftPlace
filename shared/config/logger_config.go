package config

type LoggerConfig struct {
	Level     string
	Format    string
	Outputs   []string
	FilePath  string
	AddSource bool
}
