package models

import (
	"time"

	"gorm.io/gorm"
)

// Survei Kepuasan Masyarakat (SKM). Responden mengisi tanpa akun; seluruh
// pengelolaan dilakukan users dengan role "admin".

type SkmLayanan struct {
	ID        int64     `json:"id" gorm:"primaryKey;column:layanan_id"`
	Nama      string    `json:"nama"`
	Aktif     bool      `json:"-"`
	Urutan    int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (SkmLayanan) TableName() string { return "skm_layanan" }

type SkmPeriode struct {
	ID           int64      `json:"-" gorm:"primaryKey;column:periode_id"`
	Tahun        int        `json:"tahun"`
	SkalaMaks    int        `json:"skala_maks" gorm:"column:skala_maks"`
	TanggalTutup *time.Time `json:"-" gorm:"column:tanggal_tutup;type:date"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
}

func (SkmPeriode) TableName() string { return "skm_periode" }

// Pengali konversi NRR ke indeks: 20 untuk skala 1–5, 25 untuk skala 1–4.
func (p SkmPeriode) Pengali() float64 { return 100 / float64(p.SkalaMaks) }

// Pengaturan survei. Tabelnya hanya berisi satu baris.
type SkmPengaturan struct {
	ID             int64      `gorm:"primaryKey;column:pengaturan_id"`
	Judul          string     `gorm:"column:judul"`
	PeriodeAktifID int64      `gorm:"column:periode_aktif_id"`
	PeriodeAktif   SkmPeriode `gorm:"foreignKey:PeriodeAktifID;references:ID"`
	Anonim         bool       `gorm:"column:anonim"`
	Publik         bool       `gorm:"column:publik"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (SkmPengaturan) TableName() string { return "skm_pengaturan" }

type SkmPertanyaan struct {
	ID        int64          `json:"id" gorm:"primaryKey;column:pertanyaan_id"`
	Kode      string         `json:"kode"`
	Unsur     string         `json:"unsur"`
	Teks      string         `json:"teks"`
	Wajib     bool           `json:"wajib"`
	Aktif     bool           `json:"aktif"`
	Urutan    int            `json:"urutan"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (SkmPertanyaan) TableName() string { return "skm_pertanyaan" }

// Satu pengisian survei. Nama dan ip_hash tidak pernah dikirim ke klien.
type SkmResponden struct {
	ID         int64     `gorm:"primaryKey;column:responden_id"`
	PeriodeID  int64     `gorm:"column:periode_id"`
	Urut       int       `gorm:"column:urut"`
	Nomor      string    `gorm:"column:nomor"`
	LayananID  int64     `gorm:"column:layanan_id"`
	Nama       *string   `gorm:"column:nama"`
	JK         string    `gorm:"column:jk"`
	Usia       int       `gorm:"column:usia"`
	Pendidikan string    `gorm:"column:pendidikan"`
	Pekerjaan  string    `gorm:"column:pekerjaan"`
	Frekuensi  string    `gorm:"column:frekuensi"`
	Kecamatan  *string   `gorm:"column:kecamatan"`
	Apresiasi  *string   `gorm:"column:apresiasi"`
	RataRata   float64   `gorm:"column:rata_rata"`
	IPHash     *string   `gorm:"column:ip_hash"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (SkmResponden) TableName() string { return "skm_responden" }

type SkmJawaban struct {
	ID           int64 `gorm:"primaryKey;column:jawaban_id"`
	RespondenID  int64 `gorm:"column:responden_id"`
	PertanyaanID int64 `gorm:"column:pertanyaan_id"`
	Nilai        int   `gorm:"column:nilai"`
}

func (SkmJawaban) TableName() string { return "skm_jawaban" }

const (
	SkmFeedbackBaru            = "Baru"
	SkmFeedbackDitinjau        = "Ditinjau"
	SkmFeedbackDitindaklanjuti = "Ditindaklanjuti"
	SkmFeedbackSelesai         = "Selesai"
)

// Status feedback hanya boleh maju satu langkah.
var SkmFeedbackBerikut = map[string]string{
	SkmFeedbackBaru:            SkmFeedbackDitinjau,
	SkmFeedbackDitinjau:        SkmFeedbackDitindaklanjuti,
	SkmFeedbackDitindaklanjuti: SkmFeedbackSelesai,
}

type SkmFeedback struct {
	ID            int64     `gorm:"primaryKey;column:feedback_id"`
	RespondenID   int64     `gorm:"column:responden_id"`
	KategoriID    *int64    `gorm:"column:kategori_id"`
	Isi           string    `gorm:"column:isi"`
	Status        string    `gorm:"column:status"`
	Catatan       *string   `gorm:"column:catatan"`
	DitanganiOleh *int64    `gorm:"column:ditangani_oleh"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (SkmFeedback) TableName() string { return "skm_feedback" }

type SkmPesanKontak struct {
	ID        int64     `gorm:"primaryKey;column:pesan_id"`
	Nama      string    `gorm:"column:nama"`
	Email     string    `gorm:"column:email"`
	Topik     string    `gorm:"column:topik"`
	Pesan     string    `gorm:"column:pesan"`
	Dibaca    bool      `gorm:"column:dibaca"`
	IPHash    *string   `gorm:"column:ip_hash"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (SkmPesanKontak) TableName() string { return "skm_pesan_kontak" }
