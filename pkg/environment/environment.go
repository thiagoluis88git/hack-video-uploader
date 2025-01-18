package environment

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	RedocFolderPath *string = flag.String("PATH_REDOC_FOLDER", "/docs/swagger.json", "Swagger docs folder")

	localDev = flag.String("localDev", "false", "local development")
)

const (
	DBHost     = "DB_HOST"
	DBUser     = "DB_USER"
	DBPassword = "DB_PASSWORD"
	DBPort     = "DB_PORT"
	DBName     = "DB_NAME"
	Region     = "AWS_REGION"
	S3Bucket   = "AWS_S3_BUCKET"
)

type Environment struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	Region     string
	S3Bucket   string
}

func LoadEnvironmentVariables() Environment {
	flag.Parse()

	if localFlag := *localDev; localFlag != "false" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file", err.Error())
		}
	}

	dbHost := getEnvironmentVariable(DBHost)
	dbPort := getEnvironmentVariable(DBPort)
	dbUser := getEnvironmentVariable(DBUser)
	dbPassword := getEnvironmentVariable(DBPassword)
	dbName := getEnvironmentVariable(DBName)
	region := getEnvironmentVariable(Region)
	s3Bucket := getEnvironmentVariable(S3Bucket)

	return Environment{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		Region:     region,
		S3Bucket:   s3Bucket,
	}

}

func getEnvironmentVariable(key string) string {
	value, hashKey := os.LookupEnv(key)

	if !hashKey {
		log.Fatalf("There is no %v environment variable", key)
	}

	return value
}
