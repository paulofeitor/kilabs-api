package config

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Driver  string
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	Charset string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Driver:  "mysql",
			Host:    "localhost",
			Port:    "8889",
			User:    "root",
			Pass:    "root",
			Name:    "kilabs",
			Charset: "utf8",
		},
	}
}
