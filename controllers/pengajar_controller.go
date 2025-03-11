package controllers

import (
	"net/http"
	"time"

	"bimbel-absensi/config"
	"bimbel-absensi/models"

	"github.com/gin-gonic/gin"
)

// AbsensiInput untuk request penambahan absensi
type AbsensiInput struct {
	SiswaID         uint      `json:"siswa_id" binding:"required"`
	MataPelajaranID uint      `json:"mata_pelajaran_id" binding:"required"`
	Tanggal         string    `json:"tanggal" binding:"required"` // Format: "2006-01-02"
	Status          string    `json:"status" binding:"required"`
	Keterangan      string    `json:"keterangan"`
}

// TambahAbsensi handler untuk menambahkan absensi baru
func TambahAbsensi(c *gin.Context) {
	var input AbsensiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validasi status
	if input.Status != "hadir" && input.Status != "izin" && input.Status != "sakit" && input.Status != "alpha" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus berupa hadir, izin, sakit, atau alpha"})
		return
	}
	
	// Dapatkan ID pengajar dari token
	pengajarID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID tidak ditemukan di token"})
		return
	}
	
	// Validasi bahwa pengajar ini mengajar mata pelajaran tersebut
	var pengajarMataPelajaran models.PengajarMataPelajaran
	if err := config.DB.Where("pengajar_id = ? AND mata_pelajaran_id = ?", pengajarID, input.MataPelajaranID).First(&pengajarMataPelajaran).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Pengajar tidak berwenang untuk mata pelajaran ini"})
		return
	}
	
	// Validasi bahwa siswa mengambil mata pelajaran tersebut
	var siswaMataPelajaran models.SiswaMataPelajaran
	if err := config.DB.Where("siswa_id = ? AND mata_pelajaran_id = ?", input.SiswaID, input.MataPelajaranID).First(&siswaMataPelajaran).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Siswa tidak terdaftar pada mata pelajaran ini"})
		return
	}
	
	// Parse tanggal
	tanggal, err := time.Parse("2006-01-02", input.Tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah, gunakan format YYYY-MM-DD"})
		return
	}
	
	// Buat absensi baru
	absensi := models.Absensi{
		SiswaID:         input.SiswaID,
		PengajarID:      uint(pengajarID.(uint)),
		MataPelajaranID: input.MataPelajaranID,
		Tanggal:         tanggal,
		Status:          input.Status,
		Keterangan:      input.Keterangan,
	}
	
	if err := config.DB.Create(&absensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat absensi"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Absensi berhasil ditambahkan",
		"data": absensi,
	})
}

// GetSiswaByPengajar handler untuk mendapatkan daftar siswa berdasarkan mata pelajaran pengajar
func GetSiswaByPengajar(c *gin.Context) {
	pengajarID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID tidak ditemukan di token"})
		return
	}
	
	// Dapatkan daftar mata pelajaran yang diajar oleh pengajar
	var mataPelajaranIDs []uint
	if err := config.DB.Model(&models.PengajarMataPelajaran{}).Where("pengajar_id = ?", pengajarID).Pluck("mata_pelajaran_id", &mataPelajaranIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan mata pelajaran"})
		return
	}
	
	// Dapatkan daftar siswa untuk mata pelajaran tersebut
	var siswaList []struct {
		SiswaID         uint   `json:"siswa_id"`
		Nama            string `json:"nama"`
		Email           string `json:"email"`
		MataPelajaranID uint   `json:"mata_pelajaran_id"`
		MataPelajaran   string `json:"mata_pelajaran"`
	}
	
	query := `
		SELECT s.id as siswa_id, s.nama, s.email, smp.mata_pelajaran_id, mp.nama as mata_pelajaran
		FROM siswa s
		JOIN siswa_mata_pelajaran smp ON s.id = smp.siswa_id
		JOIN mata_pelajaran mp ON smp.mata_pelajaran_id = mp.id
		WHERE smp.mata_pelajaran_id IN (?)
	`
	
	if err := config.DB.Raw(query, mataPelajaranIDs).Scan(&siswaList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan daftar siswa"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": siswaList,
	})
}

// GetMataPelajaranByPengajar handler untuk mendapatkan daftar mata pelajaran yang diajar oleh pengajar
func GetMataPelajaranByPengajar(c *gin.Context) {
	pengajarID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID tidak ditemukan di token"})
		return
	}
	
	var mataPelajaranList []struct {
		ID        uint   `json:"id"`
		Nama      string `json:"nama"`
		Deskripsi string `json:"deskripsi"`
	}
	
	query := `
		SELECT mp.id, mp.nama, mp.deskripsi
		FROM mata_pelajaran mp
		JOIN pengajar_mata_pelajaran pmp ON mp.id = pmp.mata_pelajaran_id
		WHERE pmp.pengajar_id = ?
	`
	
	if err := config.DB.Raw(query, pengajarID).Scan(&mataPelajaranList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan daftar mata pelajaran"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": mataPelajaranList,
	})
}

// GetAbsensiByPengajar handler untuk mendapatkan daftar absensi yang dibuat oleh pengajar
func GetAbsensiByPengajar(c *gin.Context) {
	pengajarID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID tidak ditemukan di token"})
		return
	}
	
	var absensiList []models.Absensi
	
	if err := config.DB.Preload("Siswa").Preload("MataPelajaran").Where("pengajar_id = ?", pengajarID).Find(&absensiList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan daftar absensi"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": absensiList,
	})
}

// GetProfile handler untuk mendapatkan profil siswa
func GetProfilePengajar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID tidak ditemukan di token"})
		return
	}
	
	var pengajar models.Pengajar
	if err := config.DB.Where("user_id = ?", userID).First(&pengajar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "pengajar tidak ditemukan"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": pengajar,
	})
}