package routes

import (
	"bimbel-absensi/controllers"
	"bimbel-absensi/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Rute publik
	r.POST("/login", controllers.Login)

	// Middleware auth untuk semua rute di bawah ini
	auth := r.Group("/api")
	auth.Use(middleware.Auth())

	// Rute untuk admin
	admin := auth.Group("/admin")
	admin.Use(middleware.AdminOnly())
	{
		// Manajemen Siswa
		admin.POST("/siswa", controllers.TambahSiswa)
		admin.GET("/siswa", controllers.GetAllSiswa) // Tambahkan fungsi ini nanti

		// Manajemen Pengajar
		admin.POST("/pengajar", controllers.TambahPengajar)
		admin.GET("/pengajar", controllers.GetAllPengajar) // Tambahkan fungsi ini nanti

		// Manajemen Mata Pelajaran
		admin.POST("/mata-pelajaran", controllers.TambahMataPelajaran)
		admin.GET("/mata-pelajaran", controllers.GetAllMataPelajaran) // Tambahkan fungsi ini nanti

		// Relasi Pengajar - Mata Pelajaran
		admin.POST("/pengajar-mata-pelajaran", controllers.TambahPengajarMataPelajaran)

		// Relasi Siswa - Mata Pelajaran
		admin.POST("/siswa-mata-pelajaran", controllers.TambahSiswaMataPelajaran)
	}

	// Rute untuk pengajar
	pengajar := auth.Group("/pengajar")
	pengajar.Use(middleware.PengajarOnly())
	{
		// Absensi
		pengajar.POST("/absensi", controllers.TambahAbsensi)
		pengajar.GET("/absensi", controllers.GetAbsensiByPengajar)

		// Daftar siswa yang diajar
		pengajar.GET("/siswa", controllers.GetSiswaByPengajar)

		// Daftar mata pelajaran yang diajar
		pengajar.GET("/mata-pelajaran", controllers.GetMataPelajaranByPengajar)

		// Profil pengajar
		pengajar.GET("/profile", controllers.GetProfilePengajar) // Tambahkan fungsi ini nanti
	}

	// Rute untuk siswa
	siswa := auth.Group("/siswa")
	{
		// Lihat absensi diri sendiri
		siswa.GET("/absensi", controllers.GetAbsensiBySiswa)

		// Lihat mata pelajaran yang diambil
		siswa.GET("/mata-pelajaran", controllers.GetMataPelajaranBySiswa)

		// Profil siswa
		siswa.GET("/profile", controllers.GetProfileSiswa)
	}

	return r
}
