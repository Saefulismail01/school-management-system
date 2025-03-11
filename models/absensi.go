package models

import (
	"time"
)

type Absensi struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SiswaID         uint      `gorm:"not null" json:"siswa_id"`
	PengajarID      uint      `gorm:"not null" json:"pengajar_id"`
	MataPelajaranID uint      `gorm:"not null" json:"mata_pelajaran_id"`
	Tanggal         time.Time `gorm:"not null" json:"tanggal"`
	Status          string    `gorm:"not null" json:"status"` // hadir, izin, sakit, alpha
	Keterangan      string    `json:"keterangan"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	
	Siswa         Siswa         `gorm:"foreignKey:SiswaID" json:"siswa"`
	Pengajar      Pengajar      `gorm:"foreignKey:PengajarID" json:"pengajar"`
	MataPelajaran MataPelajaran `gorm:"foreignKey:MataPelajaranID" json:"mata_pelajaran"`
}