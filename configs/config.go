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
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filename)
	envPath := filepath.Join(basePath, ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Failed to load .env: %v", err)
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
