package config

type LoggerConfig struct {
	Level      string
	Format     string
	Outputs    []string
	FilePath   string
	FileDir    string
	FilePrefix string
	AddSource  bool
}
