package configs

import (
	"log"
	"os"
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
	ArchivePath        string
	DBBatchSize        int
	SFT_Host           string
	SFT_Port           int
	SFT_Username       string
	SFTP_Password      string
}

func Load() (*Config, error) {

	if err := godotenv.Load("configs/.env"); err != nil {
		log.Printf("Warning: .env not %v", err)
	}

	return &Config{
		ConnString:         getEnv("CONNECTION_STRING"),
		CustomerfePath:     getEnv("CustomerfePath"),
		EZEngageCallPath:   getEnv("EZEngageCallPath"),
		HirarkiPath:        getEnv("HirarkiPath"),
		MasterMaterialPath: getEnv("MasterMaterialPath"),
		MasterOutletPath:   getEnv("MasterOutletPath"),
		SalesFEPath:        getEnv("SalesFEPath"),
		ArchivePath:        getEnv("ArchivePath"),
		DBBatchSize:        toInt(getEnv("DB_BATCHSIZE")),
		SFT_Host:           getEnv("SFT_Host"),
		SFT_Port:           toInt(getEnv("SFT_Port")),
		SFT_Username:       getEnv("SFT_Username"),
		SFTP_Password:      getEnv("SFTP_Password"),
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
