package utils

type EnvFile struct {
	Name string
}

// default env file
var envFile string = ".env"

func GetEnvFile() EnvFile {
	return EnvFile{Name: envFile}
}
