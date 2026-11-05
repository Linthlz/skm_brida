package services

import (
	"time"

	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

// SkmService: Survei Kepuasan Masyarakat. Bagian publik (responden tanpa akun)
// ada di skm_survei.go, perhitungan indeks di skm_hasil.go, pengelolaan admin
// di skm_admin.go, dan ekspor laporan di skm_laporan.go.
type SkmService interface {
	// Publik
	Referensi() (*SkmReferensi, error)
	Form() (*SkmForm, error)
	Kirim(in SkmKirimSurvei, ip string) (*SkmHasilKirim, error)
	KirimKontak(in SkmKirimKontak, ip string) error
	DaftarPeriodePublik() ([]SkmPeriodePublik, error)
	PastikanPublik() error

	// Hasil
	Periode(tahun int) (*models.SkmPeriode, error)
	Filter(layananID int64, pendidikan string) (repositories.SkmFilter, error)
	Rekap(p *models.SkmPeriode, f repositories.SkmFilter) (*SkmRekap, error)
	Ringkasan() (*SkmRingkasan, error)
	Dashboard(p *models.SkmPeriode) (*SkmDashboard, error)
	Statistik(p *models.SkmPeriode) (*SkmStatistik, error)

	// Admin
	DaftarResponden(q SkmRespondenParam) ([]SkmRespondenItem, int64, error)
	DetailResponden(id int64) (*SkmRespondenDetail, error)
	DaftarPertanyaan() ([]models.SkmPertanyaan, error)
	TambahPertanyaan(in SkmPertanyaanInput) (*models.SkmPertanyaan, error)
	UbahPertanyaan(id int64, in SkmPertanyaanInput) (*models.SkmPertanyaan, error)
	HapusPertanyaan(id int64) error
	UrutkanPertanyaan(ids []int64) ([]models.SkmPertanyaan, error)
	DaftarFeedback(status, search string, limit, offset int) ([]SkmFeedbackItem, int64, error)
	HitungStatusFeedback() (map[string]int64, error)
	DetailFeedback(id int64) (*SkmFeedbackItem, error)
	UbahFeedback(id int64, in SkmFeedbackInput, olehUserID int64) (*SkmFeedbackItem, error)
	DaftarPesan(dibaca *bool, limit, offset int) ([]SkmPesanItem, int64, error)
	JumlahPesanBelumDibaca() (int64, error)
	UbahPesan(id int64, dibaca bool) (*SkmPesanItem, error)
	Pengaturan() (*SkmPengaturanView, error)
	SimpanPengaturan(in SkmPengaturanInput) (*SkmPengaturanView, error)
	Laporan(p *models.SkmPeriode, f repositories.SkmFilter, format string) (*SkmBerkas, error)
}

type SkmServices struct {
	repo repositories.SkmRepository
}

func NewSkmService(repo repositories.SkmRepository) SkmService {
	return &SkmServices{repo: repo}
}

// ---------- Bentuk respons ----------

type SkmReferensi struct {
	Layanan        []models.SkmLayanan `json:"layanan"`
	JK             []string            `json:"jk"`
	Pendidikan     []string            `json:"pendidikan"`
	Pekerjaan      []string            `json:"pekerjaan"`
	Frekuensi      []string            `json:"frekuensi"`
	Kecamatan      []string            `json:"kecamatan"`
	TopikKontak    []string            `json:"topik_kontak"`
	StatusFeedback []string            `json:"status_feedback"`
}

type SkmFormPeriode struct {
	Tahun        int     `json:"tahun"`
	SkalaMaks    int     `json:"skala_maks"`
	TanggalTutup *string `json:"tanggal_tutup"`
	Ditutup      bool    `json:"ditutup"`
}

type SkmFormPertanyaan struct {
	ID    int64  `json:"id"`
	Kode  string `json:"kode"`
	Unsur string `json:"unsur"`
	Teks  string `json:"teks"`
	Wajib bool   `json:"wajib"`
}

type SkmForm struct {
	Judul      string              `json:"judul"`
	Anonim     bool                `json:"anonim"`
	Periode    SkmFormPeriode      `json:"periode"`
	Skala      []utils.SkmSkala    `json:"skala"`
	Pertanyaan []SkmFormPertanyaan `json:"pertanyaan"`
}

type SkmHasilKirim struct {
	Nomor     string         `json:"nomor"`
	Rata      float64        `json:"rata"`
	SkalaMaks int            `json:"skala_maks"`
	Indeks    float64        `json:"indeks"`
	Mutu      *utils.SkmMutu `json:"mutu"`
}

type SkmPeriodePublik struct {
	Tahun int  `json:"tahun"`
	Aktif bool `json:"aktif"`
}

type SkmUnsur struct {
	ID            int64          `json:"id"`
	Kode          string         `json:"kode"`
	Nama          string         `json:"nama"`
	NRR           float64        `json:"nrr"`
	JumlahJawaban int64          `json:"jumlah_jawaban"`
	NRRTertimbang float64        `json:"nrr_tertimbang"`
	Indeks        float64        `json:"indeks"`
	Mutu          *utils.SkmMutu `json:"mutu"`
}

type SkmDistribusi struct {
	V      int    `json:"v"`
	Label  string `json:"label"`
	Jumlah int64  `json:"jumlah"`
}

type SkmTren struct {
	Periode string   `json:"periode"`
	IKM     *float64 `json:"ikm"`
}

type SkmRekap struct {
	Periode        int             `json:"periode"`
	SkalaMaks      int             `json:"skala_maks"`
	Pengali        float64         `json:"pengali"`
	TotalResponden int64           `json:"total_responden"`
	RentangBulan   *string         `json:"rentang_bulan"`
	NRR            *float64        `json:"nrr"`
	IKM            *float64        `json:"ikm"`
	Mutu           *utils.SkmMutu  `json:"mutu"`
	Unsur          []SkmUnsur      `json:"unsur"`
	Distribusi     []SkmDistribusi `json:"distribusi"`
	Tren           []SkmTren       `json:"tren"`
}

type SkmRingkasan struct {
	Periode        int            `json:"periode"`
	IKM            *float64       `json:"ikm"`
	Mutu           *utils.SkmMutu `json:"mutu"`
	TotalResponden int64          `json:"total_responden"`
	TotalFeedback  int64          `json:"total_feedback"`
	JumlahUnsur    int            `json:"jumlah_unsur"`
	JumlahLayanan  int64          `json:"jumlah_layanan"`
}

type SkmPerBulan struct {
	Bulan  int    `json:"bulan"`
	Label  string `json:"label"`
	Jumlah int64  `json:"jumlah"`
}

type SkmDashboard struct {
	Periode          int             `json:"periode"`
	SkalaMaks        int             `json:"skala_maks"`
	Pengali          float64         `json:"pengali"`
	TotalResponden   int64           `json:"total_responden"`
	BulanIni         int64           `json:"bulan_ini"`
	PersenLengkap    *int64          `json:"persen_lengkap"`
	IKM              *float64        `json:"ikm"`
	Mutu             *utils.SkmMutu  `json:"mutu"`
	FeedbackTotal    int64           `json:"feedback_total"`
	FeedbackBaru     int64           `json:"feedback_baru"`
	PesanBelumDibaca int64           `json:"pesan_belum_dibaca"`
	PerBulan         []SkmPerBulan   `json:"per_bulan"`
	Distribusi       []SkmDistribusi `json:"distribusi"`
	Unsur            []SkmUnsur      `json:"unsur"`
	Prioritas        []SkmUnsur      `json:"prioritas"`
}

type SkmHitung struct {
	Label string `json:"label"`
	N     int64  `json:"n"`
}

type SkmStatistik struct {
	Periode        int         `json:"periode"`
	TotalResponden int64       `json:"total_responden"`
	Pendidikan     []SkmHitung `json:"pendidikan"`
	Pekerjaan      []SkmHitung `json:"pekerjaan"`
	Tren           []SkmTren   `json:"tren"`
}

type SkmRespondenItem struct {
	ID         int64   `json:"id"`
	Nomor      string  `json:"nomor"`
	Usia       int     `json:"usia"`
	JK         string  `json:"jk"`
	Pendidikan string  `json:"pendidikan"`
	Pekerjaan  string  `json:"pekerjaan"`
	Layanan    string  `json:"layanan"`
	Nilai      float64 `json:"nilai"`
	Kategori   string  `json:"kategori"`
	Tanggal    string  `json:"tanggal"`
}

type SkmJawabanItem struct {
	Kode  string `json:"kode"`
	Unsur string `json:"unsur"`
	Nilai int    `json:"nilai"`
}

type SkmRespondenDetail struct {
	SkmRespondenItem
	Kecamatan *string          `json:"kecamatan"`
	Frekuensi string           `json:"frekuensi"`
	Apresiasi *string          `json:"apresiasi"`
	SkalaMaks int              `json:"skala_maks"`
	Indeks    float64          `json:"indeks"`
	Jawaban   []SkmJawabanItem `json:"jawaban"`
	Saran     *string          `json:"saran"`
}

type SkmKategori struct {
	ID   int64  `json:"id"`
	Nama string `json:"nama"`
}

type SkmFeedbackItem struct {
	ID             int64        `json:"id"`
	Nomor          string       `json:"nomor"`
	Tanggal        string       `json:"tanggal"`
	Isi            string       `json:"isi"`
	Status         string       `json:"status"`
	StatusBerikut  *string      `json:"status_berikut"`
	Kategori       *SkmKategori `json:"kategori"`
	Layanan        *string      `json:"layanan"`
	NomorResponden *string      `json:"nomor_responden"`
	Catatan        *string      `json:"catatan"`
	DitanganiOleh  *string      `json:"ditangani_oleh"`
	Diperbarui     string       `json:"diperbarui"`
}

type SkmPesanItem struct {
	ID      int64  `json:"id"`
	Nama    string `json:"nama"`
	Email   string `json:"email"`
	Topik   string `json:"topik"`
	Pesan   string `json:"pesan"`
	Dibaca  bool   `json:"dibaca"`
	Tanggal string `json:"tanggal"`
}

type SkmPengelolaItem struct {
	ID      string `json:"id"`
	Nama    string `json:"nama"`
	Jabatan string `json:"jabatan"`
	Inisial string `json:"inisial"`
}

type SkmPengaturanView struct {
	Judul         string             `json:"judul"`
	Periode       int                `json:"periode"`
	SkalaMaks     int                `json:"skala_maks"`
	SkalaTerkunci bool               `json:"skala_terkunci"`
	TanggalTutup  *string            `json:"tanggal_tutup"`
	Anonim        bool               `json:"anonim"`
	Publik        bool               `json:"publik"`
	OpsiPeriode   []int              `json:"opsi_periode"`
	Pengelola     []SkmPengelolaItem `json:"pengelola"`
}

type SkmBerkas struct {
	Nama        string
	ContentType string
	Isi         []byte
}

// ---------- Input ----------

type SkmJawabanInput struct {
	PertanyaanID int64 `json:"pertanyaan_id"`
	Nilai        int   `json:"nilai"`
}

type SkmKirimSurvei struct {
	Nama       *string           `json:"nama"`
	JK         string            `json:"jk"`
	Usia       int               `json:"usia"`
	Pendidikan string            `json:"pendidikan"`
	Pekerjaan  string            `json:"pekerjaan"`
	Frekuensi  string            `json:"frekuensi"`
	Kecamatan  *string           `json:"kecamatan"`
	LayananID  int64             `json:"layanan_id"`
	Jawaban    []SkmJawabanInput `json:"jawaban"`
	Saran      *string           `json:"saran"`
	Apresiasi  *string           `json:"apresiasi"`
	Setuju     bool              `json:"setuju"`
	Website    string            `json:"website"` // honeypot: hanya diisi bot
}

type SkmKirimKontak struct {
	Nama    string `json:"nama"`
	Email   string `json:"email"`
	Topik   string `json:"topik"`
	Pesan   string `json:"pesan"`
	Website string `json:"website"` // honeypot
}

type SkmRespondenParam struct {
	Tahun     int
	LayananID int64
	Search    string
	Sort      string
	Dir       string
	Limit     int
	Offset    int
}

// SkmPertanyaanInput: nil berarti field tidak dikirim (PATCH sebagian).
type SkmPertanyaanInput struct {
	Unsur *string `json:"unsur"`
	Teks  *string `json:"teks"`
	Wajib *bool   `json:"wajib"`
	Aktif *bool   `json:"aktif"`
}

type SkmFeedbackInput struct {
	Status        *string `json:"status"`
	Catatan       *string `json:"catatan"`
	KategoriID    *int64  `json:"kategori_id"`
	CatatanAda    bool    `json:"-"`
	KategoriIDAda bool    `json:"-"`
}

type SkmPengaturanInput struct {
	Judul        string  `json:"judul"`
	Periode      int     `json:"periode"`
	SkalaMaks    int     `json:"skala_maks"`
	TanggalTutup *string `json:"tanggal_tutup"`
	Anonim       *bool   `json:"anonim"`
	Publik       *bool   `json:"publik"`
}

// ---------- Bantuan ----------

// skmGagal: galat validasi berbentuk "field: pesan" agar frontend dapat
// menandai field yang salah.
func skmGagal(field, pesan string) error {
	return utils.ErrValidation(field + ": " + pesan)
}

func skmTanggal(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

// skmDitutup: survei ditutup setelah hari tanggal_tutup berakhir (WITA).
func skmDitutup(p models.SkmPeriode) bool {
	if p.TanggalTutup == nil {
		return false
	}
	now := time.Now().In(utils.SkmZona)
	hariIni := now.Format("2006-01-02")
	return hariIni > p.TanggalTutup.Format("2006-01-02")
}
