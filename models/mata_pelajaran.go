package models

import (
	"time"
)

// MataPelajaran model
type MataPelajaran struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"not null" json:"nama"`
	Deskripsi string    `json:"deskripsi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PengajarMataPelajaran model (relasi many-to-many)
type PengajarMataPelajaran struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	PengajarID      uint      `gorm:"not null" json:"pengajar_id"`
	MataPelajaranID uint      `gorm:"not null" json:"mata_pelajaran_id"`
	CreatedAt       time.Time `json:"created_at"`
	
	Pengajar      Pengajar      `gorm:"foreignKey:PengajarID" json:"pengajar"`
	MataPelajaran MataPelajaran `gorm:"foreignKey:MataPelajaranID" json:"mata_pelajaran"`
}

// SiswaMataPelajaran model (relasi many-to-many)
type SiswaMataPelajaran struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SiswaID         uint      `gorm:"not null" json:"siswa_id"`
	MataPelajaranID uint      `gorm:"not null" json:"mata_pelajaran_id"`
	CreatedAt       time.Time `json:"created_at"`
	
	Siswa         Siswa         `gorm:"foreignKey:SiswaID" json:"siswa"`
	MataPelajaran MataPelajaran `gorm:"foreignKey:MataPelajaranID" json:"mata_pelajaran"`
}
