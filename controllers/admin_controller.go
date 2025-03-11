package controllers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"bimbel-absensi/config"
	"bimbel-absensi/models"

	"github.com/gin-gonic/gin"
)

// SiswaInput untuk request penambahan siswa
type SiswaInput struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Nama      string `json:"nama" binding:"required"`
	Email     string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Alamat    string `json:"alamat"`
}

// PengajarInput untuk request penambahan pengajar
type PengajarInput struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Nama      string `json:"nama" binding:"required"`
	Email     string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Alamat    string `json:"alamat"`
}

// MataPelajaranInput untuk request penambahan mata pelajaran
type MataPelajaranInput struct {
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
}

// TambahSiswa handler untuk menambahkan siswa baru
func TambahSiswa(c *gin.Context) {
	var input SiswaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat hash password"})
		return
	}

	// Buat user baru
	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     "siswa",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user"})
		return
	}

	// Buat siswa baru
	siswa := models.Siswa{
		UserID:    user.ID,
		Nama:      input.Nama,
		Email:     input.Email,
		NoTelepon: input.NoTelepon,
		Alamat:    input.Alamat,
	}

	if err := config.DB.Create(&siswa).Error; err != nil {
		// Rollback jika gagal
		config.DB.Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat siswa"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Siswa berhasil ditambahkan",
		"data":    siswa,
	})
}

// TambahPengajar handler untuk menambahkan pengajar baru
func TambahPengajar(c *gin.Context) {
	var input PengajarInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat hash password"})
		return
	}

	// Buat user baru
	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     "pengajar",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user"})
		return
	}

	// Buat pengajar baru
	pengajar := models.Pengajar{
		UserID:    user.ID,
		Nama:      input.Nama,
		Email:     input.Email,
		NoTelepon: input.NoTelepon,
		Alamat:    input.Alamat,
	}

	if err := config.DB.Create(&pengajar).Error; err != nil {
		// Rollback jika gagal
		config.DB.Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengajar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengajar berhasil ditambahkan",
		"data":    pengajar,
	})
}

// TambahMataPelajaran handler untuk menambahkan mata pelajaran baru
func TambahMataPelajaran(c *gin.Context) {
	var input MataPelajaranInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mataPelajaran := models.MataPelajaran{
		Nama:      input.Nama,
		Deskripsi: input.Deskripsi,
	}

	if err := config.DB.Create(&mataPelajaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat mata pelajaran"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Mata pelajaran berhasil ditambahkan",
		"data":    mataPelajaran,
	})
}

// TambahPengajarMataPelajaran handler untuk menghubungkan pengajar dengan mata pelajaran
type PengajarMataPelajaranInput struct {
	PengajarID      uint `json:"pengajar_id" binding:"required"`
	MataPelajaranID uint `json:"mata_pelajaran_id" binding:"required"`
}

func TambahPengajarMataPelajaran(c *gin.Context) {
	var input PengajarMataPelajaranInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	relation := models.PengajarMataPelajaran{
		PengajarID:      input.PengajarID,
		MataPelajaranID: input.MataPelajaranID,
	}

	if err := config.DB.Create(&relation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungkan pengajar dengan mata pelajaran"})
		return
	}

	// **Preload data pengajar dan mata pelajaran sebelum mengembalikan respons**
	config.DB.Preload("Pengajar").Preload("MataPelajaran").First(&relation, relation.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengajar berhasil dihubungkan dengan mata pelajaran",
		"data":    relation,
	})
}

// TambahSiswaMataPelajaran handler untuk menghubungkan siswa dengan mata pelajaran
type SiswaMataPelajaranInput struct {
	SiswaID         uint `json:"siswa_id" binding:"required"`
	MataPelajaranID uint `json:"mata_pelajaran_id" binding:"required"`
}

func TambahSiswaMataPelajaran(c *gin.Context) {
	var input SiswaMataPelajaranInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	relation := models.SiswaMataPelajaran{
		SiswaID:         input.SiswaID,
		MataPelajaranID: input.MataPelajaranID,
	}

	if err := config.DB.Create(&relation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungkan siswa dengan mata pelajaran"})
		return
	}

	// **Preload data siswa dan mata pelajaran sebelum mengembalikan respons**
	config.DB.Preload("Siswa").Preload("MataPelajaran").First(&relation, relation.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Siswa berhasil dihubungkan dengan mata pelajaran",
		"data":    relation,
	})
}

func GetAllPengajar(c *gin.Context) {
	var pengajarList []models.Pengajar

	if err := config.DB.Find(&pengajarList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan daftar absensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": pengajarList,
	})
}

func GetAllSiswa(c *gin.Context) {
	var siswaList []models.Siswa

	if err := config.DB.Find(&siswaList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan daftar absensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": siswaList,
	})
}

func GetAllMataPelajaran(c *gin.Context) {
	var mataPelajaranList []models.MataPelajaran

	if err := config.DB.Find(&mataPelajaranList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan daftar absensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": mataPelajaranList,
	})
}
