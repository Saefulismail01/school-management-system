package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"bimbel-absensi/config"
	"bimbel-absensi/models"
)

// GetAbsensiBySiswa handler untuk mendapatkan daftar absensi siswa
func GetAbsensiBySiswa(c *gin.Context) {
	siswaID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID tidak ditemukan di token"})
		return
	}
	
	// Dapatkan siswa ID dari tabel siswa
	var siswa models.Siswa
	if err := config.DB.Where("user_id = ?", siswaID).First(&siswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Siswa tidak ditemukan"})
		return
	}
	
	var absensiList []struct {
		ID              uint   `json:"id"`
		Tanggal         string `json:"tanggal"`
		Status          string `json:"status"`
		Keterangan      string `json:"keterangan"`
		MataPelajaran   string `json:"mata_pelajaran"`
		NamaPengajar    string `json:"nama_pengajar"`
	}
	
	query := `
		SELECT a.id, a.tanggal, a.status, a.keterangan, 
		       mp.nama as mata_pelajaran, p.nama as nama_pengajar
		FROM absensi a
		JOIN mata_pelajaran mp ON a.mata_pelajaran_id = mp.id
		JOIN pengajar p ON a.pengajar_id = p.id
		WHERE a.siswa_id = ?
		ORDER BY a.tanggal DESC
	`
	
	if err := config.DB.Raw(query, siswa.ID).Scan(&absensiList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan daftar absensi"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": absensiList,
	})
}

// GetMataPelajaranBySiswa handler untuk mendapatkan daftar mata pelajaran yang diambil siswa
func GetMataPelajaranBySiswa(c *gin.Context) {
	siswaID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID tidak ditemukan di token"})
		return
	}
	
	// Dapatkan siswa ID dari tabel siswa
	var siswa models.Siswa
	if err := config.DB.Where("user_id = ?", siswaID).First(&siswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Siswa tidak ditemukan"})
		return
	}
	
	var mataPelajaranList []struct {
		ID            uint   `json:"id"`
		Nama          string `json:"nama"`
		Deskripsi     string `json:"deskripsi"`
		NamaPengajar  string `json:"nama_pengajar"`
	}
	
	query := `
		SELECT mp.id, mp.nama, mp.deskripsi, p.nama as nama_pengajar
		FROM mata_pelajaran mp
		JOIN siswa_mata_pelajaran smp ON mp.id = smp.mata_pelajaran_id
		JOIN pengajar_mata_pelajaran pmp ON mp.id = pmp.mata_pelajaran_id
		JOIN pengajar p ON pmp.pengajar_id = p.id
		WHERE smp.siswa_id = ?
	`
	
	if err := config.DB.Raw(query, siswa.ID).Scan(&mataPelajaranList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan daftar mata pelajaran"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": mataPelajaranList,
	})
}

// GetProfile handler untuk mendapatkan profil siswa
func GetProfileSiswa(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID tidak ditemukan di token"})
		return
	}
	
	var siswa models.Siswa
	if err := config.DB.Where("user_id = ?", userID).First(&siswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Siswa tidak ditemukan"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": siswa,
	})
}