package configs

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ConnString         string
	CustomerfePath     string
	EZEngageCallPath   string
	HirarkiPath        string
	MasterMaterialPath string
	MasterOutletPath   string
	SalesFEPath        string
	DBBatchSize        int
}

func Load() (*Config, error) {
	// Get the directory of the current source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("Warning: unable to get caller info")
	}

	basePath := filepath.Dir(filename)
	// Adjust if your .env is inside "configs" folder
	envPath := filepath.Join(basePath, ".env")

	// Load the .env file
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: .env not found at %s: %v", envPath, err)
	}

	return &Config{
		ConnString:         getEnv("CONNECTION_STRING"),
		CustomerfePath:     getEnv("CustomerfePath"),
		EZEngageCallPath:   getEnv("EZEngageCallPath"),
		HirarkiPath:        getEnv("HirarkiPath"),
		MasterMaterialPath: getEnv("MasterMaterialPath"),
		MasterOutletPath:   getEnv("MasterOutletPath"),
		SalesFEPath:        getEnv("SalesFEPath"),
		DBBatchSize:        toInt(getEnv("DB_BATCHSIZE")),
	}, nil
}

func getEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return ""
}

func toInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}
