package repositories

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SkmFilter: filter bersama hasil, dashboard, dan laporan.
type SkmFilter struct {
	LayananID  int64
	Pendidikan string
}

type SkmNRR struct {
	PertanyaanID int64
	NRR          float64
	N            int64
}

type SkmRespondenQuery struct {
	PeriodeID int64
	LayananID int64
	Search    string
	Sort      string // kolom DB dari daftar putih di service
	Desc      bool
	Limit     int
	Offset    int
}

type SkmRespondenRow struct {
	models.SkmResponden
	Layanan   string
	SkalaMaks int
}

type SkmJawabanRow struct {
	Kode  string
	Unsur string
	Nilai int
}

type SkmFeedbackRow struct {
	models.SkmFeedback
	Kategori       *string
	Layanan        *string
	NomorResponden *string
	Penangan       *string
}

type SkmPengelola struct {
	PublicID string
	Name     *string
	Position *string
}

type SkmPengaturanBaru struct {
	Judul        string
	Tahun        int
	SkalaMaks    int
	TanggalTutup *time.Time
	Anonim       bool
	Publik       bool
}

var ErrSkmSkalaTerkunci = errors.New("skala terkunci")

type SkmRepository interface {
	Pengaturan() (*models.SkmPengaturan, error)
	PeriodeByTahun(tahun int) (*models.SkmPeriode, error)
	DaftarPeriode() ([]models.SkmPeriode, error)
	PeriodeByIDs(ids []int64) ([]models.SkmPeriode, error)
	PeriodePunyaResponden(periodeID int64) (bool, error)
	SimpanPengaturan(in SkmPengaturanBaru) error

	LayananAktif() ([]models.SkmLayanan, error)
	LayananByID(id int64) (*models.SkmLayanan, error)
	JumlahLayananAktif() (int64, error)

	PertanyaanAktif() ([]models.SkmPertanyaan, error)
	SemuaPertanyaan() ([]models.SkmPertanyaan, error)
	PertanyaanByIDs(ids []int64) ([]models.SkmPertanyaan, error)
	PertanyaanByID(id int64) (*models.SkmPertanyaan, error)
	TambahPertanyaan(p *models.SkmPertanyaan) error
	UbahPertanyaan(id int64, updates map[string]interface{}) error
	HapusPertanyaan(id int64) error
	UrutkanPertanyaan(ids []int64) error

	SimpanSurvei(r *models.SkmResponden, jawaban []models.SkmJawaban, fb *models.SkmFeedback) error

	NRRUnsur(periodeID int64, f SkmFilter) ([]SkmNRR, error)
	NRRPerPeriode(f SkmFilter) (map[int64][]float64, error)
	JumlahResponden(periodeID int64, f SkmFilter) (int64, error)
	RespondenSejak(periodeID int64, sejak time.Time) (int64, error)
	Distribusi(periodeID int64, f SkmFilter) (map[int]int64, error)
	RentangWaktu(periodeID int64, f SkmFilter) (*time.Time, *time.Time, error)
	PerBulan(periodeID int64) (map[int]int64, error)
	JumlahLengkap(periodeID int64) (int64, error)
	HitungPer(periodeID int64, kolom string) (map[string]int64, error)

	DaftarResponden(q SkmRespondenQuery) ([]SkmRespondenRow, int64, error)
	RespondenByID(id int64) (*SkmRespondenRow, error)
	JawabanResponden(respondenID int64) ([]SkmJawabanRow, error)
	SaranResponden(respondenID int64) (*string, error)

	JumlahFeedbackPeriode(periodeID int64) (int64, error)
	HitungStatusFeedback() (map[string]int64, error)
	DaftarFeedback(status, search string, limit, offset int) ([]SkmFeedbackRow, int64, error)
	FeedbackByID(id int64) (*SkmFeedbackRow, error)
	UbahFeedback(id int64, updates map[string]interface{}) error

	SimpanPesan(p *models.SkmPesanKontak) error
	DaftarPesan(dibaca *bool, limit, offset int) ([]models.SkmPesanKontak, int64, error)
	PesanByID(id int64) (*models.SkmPesanKontak, error)
	UbahPesan(id int64, dibaca bool) error
	JumlahPesanBelumDibaca() (int64, error)

	Pengelola() ([]SkmPengelola, error)
}

type SkmRepositorys struct{}

func NewSkmRepository() SkmRepository {
	return &SkmRepositorys{}
}

func skmFirst[T any](db *gorm.DB) (*T, error) {
	var out T
	err := db.First(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Pengaturan & periode ----------

func (r *SkmRepositorys) Pengaturan() (*models.SkmPengaturan, error) {
	p, err := skmFirst[models.SkmPengaturan](config.DB.Preload("PeriodeAktif").Order("pengaturan_id"))
	if err == nil && p == nil {
		return nil, errors.New("skm_pengaturan kosong: jalankan migrasi 000028_seed_skm_master")
	}
	return p, err
}

func (r *SkmRepositorys) PeriodeByTahun(tahun int) (*models.SkmPeriode, error) {
	return skmFirst[models.SkmPeriode](config.DB.Where("tahun = ?", tahun))
}

func (r *SkmRepositorys) DaftarPeriode() ([]models.SkmPeriode, error) {
	var list []models.SkmPeriode
	err := config.DB.Order("tahun DESC").Find(&list).Error
	return list, err
}

func (r *SkmRepositorys) PeriodeByIDs(ids []int64) ([]models.SkmPeriode, error) {
	var list []models.SkmPeriode
	if len(ids) == 0 {
		return list, nil
	}
	err := config.DB.Where("periode_id IN ?", ids).Order("tahun ASC").Find(&list).Error
	return list, err
}

func (r *SkmRepositorys) PeriodePunyaResponden(periodeID int64) (bool, error) {
	var n int64
	err := config.DB.Model(&models.SkmResponden{}).Where("periode_id = ?", periodeID).Limit(1).Count(&n).Error
	return n > 0, err
}

// SimpanPengaturan: skala dan tanggal tutup melekat pada periode aktif. Skala
// tidak boleh diubah bila periode itu sudah punya responden.
func (r *SkmRepositorys) SimpanPengaturan(in SkmPengaturanBaru) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		periode := models.SkmPeriode{Tahun: in.Tahun, SkalaMaks: in.SkalaMaks}
		if err := tx.Where("tahun = ?", in.Tahun).FirstOrCreate(&periode).Error; err != nil {
			return err
		}
		if periode.SkalaMaks != in.SkalaMaks {
			var n int64
			if err := tx.Model(&models.SkmResponden{}).Where("periode_id = ?", periode.ID).Limit(1).Count(&n).Error; err != nil {
				return err
			}
			if n > 0 {
				return ErrSkmSkalaTerkunci
			}
		}
		if err := tx.Model(&models.SkmPeriode{}).Where("periode_id = ?", periode.ID).Updates(map[string]interface{}{
			"skala_maks": in.SkalaMaks, "tanggal_tutup": in.TanggalTutup,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&models.SkmPengaturan{}).Where("1 = 1").Updates(map[string]interface{}{
			"judul": in.Judul, "periode_aktif_id": periode.ID, "anonim": in.Anonim, "publik": in.Publik,
		}).Error
	})
}

// ---------- Layanan ----------

func (r *SkmRepositorys) LayananAktif() ([]models.SkmLayanan, error) {
	var list []models.SkmLayanan
	err := config.DB.Where("aktif = ?", true).Order("urutan, layanan_id").Find(&list).Error
	return list, err
}

func (r *SkmRepositorys) LayananByID(id int64) (*models.SkmLayanan, error) {
	return skmFirst[models.SkmLayanan](config.DB.Where("layanan_id = ?", id))
}

func (r *SkmRepositorys) JumlahLayananAktif() (int64, error) {
	var n int64
	err := config.DB.Model(&models.SkmLayanan{}).Where("aktif = ?", true).Count(&n).Error
	return n, err
}

// ---------- Pertanyaan ----------

func (r *SkmRepositorys) PertanyaanAktif() ([]models.SkmPertanyaan, error) {
	var list []models.SkmPertanyaan
	err := config.DB.Where("aktif = ?", true).Order("urutan, pertanyaan_id").Find(&list).Error
	return list, err
}

func (r *SkmRepositorys) SemuaPertanyaan() ([]models.SkmPertanyaan, error) {
	var list []models.SkmPertanyaan
	err := config.DB.Order("urutan, pertanyaan_id").Find(&list).Error
	return list, err
}

// PertanyaanByIDs termasuk yang sudah dihapus, untuk rekap periode lama.
func (r *SkmRepositorys) PertanyaanByIDs(ids []int64) ([]models.SkmPertanyaan, error) {
	var list []models.SkmPertanyaan
	if len(ids) == 0 {
		return list, nil
	}
	err := config.DB.Unscoped().Where("pertanyaan_id IN ?", ids).Order("urutan, pertanyaan_id").Find(&list).Error
	return list, err
}

func (r *SkmRepositorys) PertanyaanByID(id int64) (*models.SkmPertanyaan, error) {
	return skmFirst[models.SkmPertanyaan](config.DB.Where("pertanyaan_id = ?", id))
}

var angka = regexp.MustCompile(`\D`)

// TambahPertanyaan: kode U{n} berikutnya (kode terhapus tidak dipakai ulang) dan urutan paling akhir.
func (r *SkmRepositorys) TambahPertanyaan(p *models.SkmPertanyaan) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("LOCK TABLE skm_pertanyaan IN SHARE ROW EXCLUSIVE MODE").Error; err != nil {
			return err
		}
		var kode []string
		if err := tx.Unscoped().Model(&models.SkmPertanyaan{}).Pluck("kode", &kode).Error; err != nil {
			return err
		}
		terakhir := 0
		for _, k := range kode {
			if n, err := strconv.Atoi(angka.ReplaceAllString(k, "")); err == nil && n > terakhir {
				terakhir = n
			}
		}
		var urutan int
		if err := tx.Model(&models.SkmPertanyaan{}).Select("COALESCE(MAX(urutan), 0)").Scan(&urutan).Error; err != nil {
			return err
		}
		p.Kode = fmt.Sprintf("U%d", terakhir+1)
		p.Urutan = urutan + 1
		return tx.Create(p).Error
	})
}

func (r *SkmRepositorys) UbahPertanyaan(id int64, updates map[string]interface{}) error {
	return config.DB.Model(&models.SkmPertanyaan{}).Where("pertanyaan_id = ?", id).Updates(updates).Error
}

func (r *SkmRepositorys) HapusPertanyaan(id int64) error {
	return config.DB.Where("pertanyaan_id = ?", id).Delete(&models.SkmPertanyaan{}).Error
}

func (r *SkmRepositorys) UrutkanPertanyaan(ids []int64) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&models.SkmPertanyaan{}).Where("pertanyaan_id = ?", id).Update("urutan", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ---------- Survei ----------

// SimpanSurvei menyimpan responden, jawaban, dan saran dalam satu transaksi.
// Baris periode dikunci agar nomor urut tidak bentrok saat pengiriman bersamaan.
func (r *SkmRepositorys) SimpanSurvei(resp *models.SkmResponden, jawaban []models.SkmJawaban, fb *models.SkmFeedback) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var periode models.SkmPeriode
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("periode_id = ?", resp.PeriodeID).First(&periode).Error; err != nil {
			return err
		}
		var urut int
		if err := tx.Model(&models.SkmResponden{}).Where("periode_id = ?", periode.ID).Select("COALESCE(MAX(urut), 0)").Scan(&urut).Error; err != nil {
			return err
		}
		resp.Urut = urut + 1
		resp.Nomor = fmt.Sprintf("SKM-%d-%04d", periode.Tahun, resp.Urut)
		if err := tx.Create(resp).Error; err != nil {
			return err
		}
		for i := range jawaban {
			jawaban[i].RespondenID = resp.ID
		}
		if err := tx.Create(&jawaban).Error; err != nil {
			return err
		}
		if fb != nil {
			fb.RespondenID = resp.ID
			return tx.Create(fb).Error
		}
		return nil
	})
}

// ---------- Agregat ----------

func skmFilterResponden(db *gorm.DB, f SkmFilter) *gorm.DB {
	if f.LayananID != 0 {
		db = db.Where("r.layanan_id = ?", f.LayananID)
	}
	if f.Pendidikan != "" {
		db = db.Where("r.pendidikan = ?", f.Pendidikan)
	}
	return db
}

func skmResponden(periodeID int64, f SkmFilter) *gorm.DB {
	return skmFilterResponden(config.DB.Table("skm_responden AS r").Where("r.periode_id = ?", periodeID), f)
}

// NRRUnsur: rata-rata nilai tiap pertanyaan yang punya jawaban pada periode itu.
func (r *SkmRepositorys) NRRUnsur(periodeID int64, f SkmFilter) ([]SkmNRR, error) {
	var rows []SkmNRR
	err := skmFilterResponden(config.DB.Table("skm_jawaban AS j").
		Joins("JOIN skm_responden AS r ON r.responden_id = j.responden_id").
		Where("r.periode_id = ?", periodeID), f).
		Select("j.pertanyaan_id, AVG(j.nilai)::float8 AS nrr, COUNT(*) AS n").
		Group("j.pertanyaan_id").
		Scan(&rows).Error
	return rows, err
}

// NRRPerPeriode: periode_id -> NRR tiap pertanyaan, untuk tren IKM antar tahun.
func (r *SkmRepositorys) NRRPerPeriode(f SkmFilter) (map[int64][]float64, error) {
	var rows []struct {
		PeriodeID int64
		NRR       float64
	}
	err := skmFilterResponden(config.DB.Table("skm_jawaban AS j").
		Joins("JOIN skm_responden AS r ON r.responden_id = j.responden_id"), f).
		Select("r.periode_id, AVG(j.nilai)::float8 AS nrr").
		Group("r.periode_id, j.pertanyaan_id").
		Scan(&rows).Error
	out := map[int64][]float64{}
	for _, row := range rows {
		out[row.PeriodeID] = append(out[row.PeriodeID], row.NRR)
	}
	return out, err
}

func (r *SkmRepositorys) JumlahResponden(periodeID int64, f SkmFilter) (int64, error) {
	var n int64
	err := skmResponden(periodeID, f).Count(&n).Error
	return n, err
}

func (r *SkmRepositorys) RespondenSejak(periodeID int64, sejak time.Time) (int64, error) {
	var n int64
	err := skmResponden(periodeID, SkmFilter{}).Where("r.created_at >= ?", sejak).Count(&n).Error
	return n, err
}

// Distribusi: jumlah responden per tingkat kepuasan (rata-rata dibulatkan).
func (r *SkmRepositorys) Distribusi(periodeID int64, f SkmFilter) (map[int]int64, error) {
	var rows []struct {
		V      int
		Jumlah int64
	}
	err := skmResponden(periodeID, f).
		Select("ROUND(r.rata_rata)::int AS v, COUNT(*) AS jumlah").
		Group("ROUND(r.rata_rata)").
		Scan(&rows).Error
	out := map[int]int64{}
	for _, row := range rows {
		out[row.V] = row.Jumlah
	}
	return out, err
}

func (r *SkmRepositorys) RentangWaktu(periodeID int64, f SkmFilter) (*time.Time, *time.Time, error) {
	var row struct {
		Awal  *time.Time
		Akhir *time.Time
	}
	err := skmResponden(periodeID, f).Select("MIN(r.created_at) AS awal, MAX(r.created_at) AS akhir").Scan(&row).Error
	return row.Awal, row.Akhir, err
}

func (r *SkmRepositorys) PerBulan(periodeID int64) (map[int]int64, error) {
	var rows []struct {
		Bulan  int
		Jumlah int64
	}
	bulan := "EXTRACT(MONTH FROM r.created_at AT TIME ZONE '" + utils.SkmZonaSQL + "')::int"
	err := skmResponden(periodeID, SkmFilter{}).
		Select(bulan + " AS bulan, COUNT(*) AS jumlah").
		Group(bulan).
		Scan(&rows).Error
	out := map[int]int64{}
	for _, row := range rows {
		out[row.Bulan] = row.Jumlah
	}
	return out, err
}

// JumlahLengkap: responden yang menjawab seluruh pertanyaan yang pernah muncul pada periode itu.
func (r *SkmRepositorys) JumlahLengkap(periodeID int64) (int64, error) {
	var n int64
	err := config.DB.Raw(`
		WITH per_responden AS (
			SELECT j.responden_id, COUNT(*) AS n
			FROM skm_jawaban j JOIN skm_responden r ON r.responden_id = j.responden_id
			WHERE r.periode_id = ?
			GROUP BY j.responden_id
		), total AS (
			SELECT COUNT(DISTINCT j.pertanyaan_id) AS n
			FROM skm_jawaban j JOIN skm_responden r ON r.responden_id = j.responden_id
			WHERE r.periode_id = ?
		)
		SELECT COUNT(*) FROM per_responden, total WHERE total.n > 0 AND per_responden.n >= total.n`,
		periodeID, periodeID).Scan(&n).Error
	return n, err
}

var skmKolomHitung = map[string]bool{"pendidikan": true, "pekerjaan": true}

func (r *SkmRepositorys) HitungPer(periodeID int64, kolom string) (map[string]int64, error) {
	if !skmKolomHitung[kolom] {
		return nil, fmt.Errorf("kolom %q tidak diizinkan", kolom)
	}
	var rows []struct {
		Label string
		N     int64
	}
	err := skmResponden(periodeID, SkmFilter{}).
		Select("r." + kolom + " AS label, COUNT(*) AS n").
		Group("r." + kolom).
		Scan(&rows).Error
	out := map[string]int64{}
	for _, row := range rows {
		out[row.Label] = row.N
	}
	return out, err
}

// ---------- Responden ----------

const skmRespondenKolom = "r.*, l.nama AS layanan, p.skala_maks AS skala_maks"

func skmRespondenRow() *gorm.DB {
	return config.DB.Table("skm_responden AS r").
		Joins("JOIN skm_layanan AS l ON l.layanan_id = r.layanan_id").
		Joins("JOIN skm_periode AS p ON p.periode_id = r.periode_id")
}

func (r *SkmRepositorys) DaftarResponden(q SkmRespondenQuery) ([]SkmRespondenRow, int64, error) {
	var list []SkmRespondenRow
	var total int64
	db := skmRespondenRow().Where("r.periode_id = ?", q.PeriodeID)
	if q.LayananID != 0 {
		db = db.Where("r.layanan_id = ?", q.LayananID)
	}
	if q.Search != "" {
		pola := utils.SkmPolaCari(q.Search)
		db = db.Where("r.nomor ILIKE ? OR r.jk ILIKE ? OR r.pekerjaan ILIKE ? OR r.pendidikan ILIKE ? OR l.nama ILIKE ?",
			pola, pola, pola, pola, pola)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	dir := "ASC"
	if q.Desc {
		dir = "DESC"
	}
	err := db.Select(skmRespondenKolom).Order(q.Sort + " " + dir).Order("r.responden_id " + dir).
		Limit(q.Limit).Offset(q.Offset).Scan(&list).Error
	return list, total, err
}

func (r *SkmRepositorys) RespondenByID(id int64) (*SkmRespondenRow, error) {
	var list []SkmRespondenRow
	if err := skmRespondenRow().Select(skmRespondenKolom).Where("r.responden_id = ?", id).Limit(1).Scan(&list).Error; err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

func (r *SkmRepositorys) JawabanResponden(respondenID int64) ([]SkmJawabanRow, error) {
	var rows []SkmJawabanRow
	err := config.DB.Table("skm_jawaban AS j").
		Joins("JOIN skm_pertanyaan AS q ON q.pertanyaan_id = j.pertanyaan_id").
		Where("j.responden_id = ?", respondenID).
		Select("q.kode, q.unsur, j.nilai").
		Order("LENGTH(q.kode), q.kode").
		Scan(&rows).Error
	return rows, err
}

func (r *SkmRepositorys) SaranResponden(respondenID int64) (*string, error) {
	var isi []string
	err := config.DB.Model(&models.SkmFeedback{}).Where("responden_id = ?", respondenID).Limit(1).Pluck("isi", &isi).Error
	if err != nil || len(isi) == 0 {
		return nil, err
	}
	return &isi[0], nil
}

// ---------- Feedback ----------

func (r *SkmRepositorys) JumlahFeedbackPeriode(periodeID int64) (int64, error) {
	var n int64
	err := config.DB.Table("skm_feedback AS f").
		Joins("JOIN skm_responden AS r ON r.responden_id = f.responden_id").
		Where("r.periode_id = ?", periodeID).Count(&n).Error
	return n, err
}

func (r *SkmRepositorys) HitungStatusFeedback() (map[string]int64, error) {
	var rows []struct {
		Status string
		N      int64
	}
	err := config.DB.Model(&models.SkmFeedback{}).Select("status, COUNT(*) AS n").Group("status").Scan(&rows).Error
	out := map[string]int64{}
	for _, row := range rows {
		out[row.Status] = row.N
	}
	return out, err
}

func skmFeedbackRow() *gorm.DB {
	return config.DB.Table("skm_feedback AS f").
		Joins("JOIN skm_responden AS r ON r.responden_id = f.responden_id").
		Joins("JOIN skm_layanan AS l ON l.layanan_id = r.layanan_id").
		Joins("LEFT JOIN skm_pertanyaan AS q ON q.pertanyaan_id = f.kategori_id").
		Joins("LEFT JOIN users AS u ON u.internal_id = f.ditangani_oleh")
}

const skmFeedbackKolom = "f.*, q.unsur AS kategori, l.nama AS layanan, r.nomor AS nomor_responden, u.name AS penangan"

func (r *SkmRepositorys) DaftarFeedback(status, search string, limit, offset int) ([]SkmFeedbackRow, int64, error) {
	var list []SkmFeedbackRow
	var total int64
	db := skmFeedbackRow()
	if status != "" {
		db = db.Where("f.status = ?", status)
	}
	if search != "" {
		pola := utils.SkmPolaCari(search)
		db = db.Where("f.isi ILIKE ? OR f.catatan ILIKE ?", pola, pola)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Select(skmFeedbackKolom).Order("f.created_at DESC, f.feedback_id DESC").Limit(limit).Offset(offset).Scan(&list).Error
	return list, total, err
}

func (r *SkmRepositorys) FeedbackByID(id int64) (*SkmFeedbackRow, error) {
	var list []SkmFeedbackRow
	if err := skmFeedbackRow().Select(skmFeedbackKolom).Where("f.feedback_id = ?", id).Limit(1).Scan(&list).Error; err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

func (r *SkmRepositorys) UbahFeedback(id int64, updates map[string]interface{}) error {
	return config.DB.Model(&models.SkmFeedback{}).Where("feedback_id = ?", id).Updates(updates).Error
}

// ---------- Pesan kontak ----------

func (r *SkmRepositorys) SimpanPesan(p *models.SkmPesanKontak) error {
	return config.DB.Create(p).Error
}

func (r *SkmRepositorys) DaftarPesan(dibaca *bool, limit, offset int) ([]models.SkmPesanKontak, int64, error) {
	var list []models.SkmPesanKontak
	var total int64
	db := config.DB.Model(&models.SkmPesanKontak{})
	if dibaca != nil {
		db = db.Where("dibaca = ?", *dibaca)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("created_at DESC, pesan_id DESC").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *SkmRepositorys) PesanByID(id int64) (*models.SkmPesanKontak, error) {
	return skmFirst[models.SkmPesanKontak](config.DB.Where("pesan_id = ?", id))
}

func (r *SkmRepositorys) UbahPesan(id int64, dibaca bool) error {
	return config.DB.Model(&models.SkmPesanKontak{}).Where("pesan_id = ?", id).Update("dibaca", dibaca).Error
}

func (r *SkmRepositorys) JumlahPesanBelumDibaca() (int64, error) {
	var n int64
	err := config.DB.Model(&models.SkmPesanKontak{}).Where("dibaca = ?", false).Count(&n).Error
	return n, err
}

// ---------- Pengelola ----------

// Pengelola: akun Singa Riset dengan role admin yang sudah terverifikasi.
func (r *SkmRepositorys) Pengelola() ([]SkmPengelola, error) {
	var list []SkmPengelola
	err := config.DB.Table("users").
		Where("role = ? AND verified = ?", "admin", true).
		Order("internal_id").
		Select("public_id::text AS public_id, name, position").
		Scan(&list).Error
	return list, err
}
