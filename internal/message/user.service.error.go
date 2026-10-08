package msgerr

import "errors"

var (
	AuthNotFound            = errors.New("gagal mengambil data otentikasi user")
	InvalidImageExt         = errors.New("format gambar tidak valid, hanya menerima .jpg, .jpeg, atau .png")
	ImageTooLarge           = errors.New("ukuran gambar melebihi batas maksimal 2MB")
	OpenFile                = errors.New("gagal membuka file gambar yang diunggah")
	ReadFile                = errors.New("gagal membaca data file gambar")
	WriteFile               = errors.New("gagal menyimpan file gambar ke direktori penyimpanan")
	UserNotFound            = errors.New("data user tidak ditemukan berdasarkan ID otentikasi")
	SavedEventsUserNotFound = errors.New("data events disimpan user tidak ditemukan berdasarkan ID otentikasi")
	JoinedEventsNotFound    = errors.New("data events joined user tidak ditemukan berdasarkan ID otentikasi")
	UpdateProfile           = errors.New("gagal memperbarui data profil di database")
)
