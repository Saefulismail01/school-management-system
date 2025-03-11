package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"bimbel-absensi/models"
)

var DB *gorm.DB

// ConnectDatabase untuk menghubungkan aplikasi ke database PostgreSQL
func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		dbHost, dbPort, dbUser, dbPassword, dbName)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// Auto Migrate models
	db.AutoMigrate(
		&models.User{},
		&models.Siswa{},
		&models.Pengajar{},
		&models.MataPelajaran{},
		&models.PengajarMataPelajaran{},
		&models.SiswaMataPelajaran{},
		&models.Absensi{},
	)
	
	DB = db
	fmt.Println("Database Connected")
}