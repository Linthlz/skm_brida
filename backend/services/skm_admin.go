package services

import (
	"errors"
	"fmt"
	"sort"
	"time"
	"unicode/utf8"

	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

// ---------- Responden ----------

// Kolom yang boleh dipakai untuk pengurutan (nama dari frontend => kolom DB).
var skmSortResponden = map[string]string{
	"nomor":      "r.urut",
	"usia":       "r.usia",
	"jk":         "r.jk",
	"pekerjaan":  "r.pekerjaan",
	"pendidikan": "r.pendidikan",
	"layanan":    "l.nama",
	"nilai":      "r.rata_rata",
	"tanggal":    "r.created_at",
}

func skmRespondenItem(r repositories.SkmRespondenRow) SkmRespondenItem {
	return SkmRespondenItem{
		ID: r.ID, Nomor: r.Nomor, Usia: r.Usia, JK: r.JK, Pendidikan: r.Pendidikan, Pekerjaan: r.Pekerjaan,
		Layanan: r.Layanan, Nilai: r.RataRata, Kategori: utils.SkmKategoriResponden(r.RataRata, r.SkalaMaks),
		Tanggal: utils.SkmISO(r.CreatedAt),
	}
}

func (s *SkmServices) DaftarResponden(q SkmRespondenParam) ([]SkmRespondenItem, int64, error) {
	periode, err := s.Periode(q.Tahun)
	if err != nil {
		return nil, 0, err
	}
	if q.Sort == "" {
		q.Sort = "tanggal"
	}
	kolom, ok := skmSortResponden[q.Sort]
	if !ok {
		return nil, 0, skmGagal("sort", "Kolom pengurutan tidak valid.")
	}
	if q.Dir == "" {
		q.Dir = "desc"
	}
	if q.Dir != "asc" && q.Dir != "desc" {
		return nil, 0, skmGagal("dir", "Arah pengurutan harus asc atau desc.")
	}
	if utf8.RuneCountInString(q.Search) > 100 {
		return nil, 0, skmGagal("search", "Kata kunci maksimal 100 karakter.")
	}
	rows, total, err := s.repo.DaftarResponden(repositories.SkmRespondenQuery{
		PeriodeID: periode.ID, LayananID: q.LayananID, Search: q.Search,
		Sort: kolom, Desc: q.Dir == "desc", Limit: q.Limit, Offset: q.Offset,
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]SkmRespondenItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, skmRespondenItem(r))
	}
	return out, total, nil
}

func (s *SkmServices) DetailResponden(id int64) (*SkmRespondenDetail, error) {
	r, err := s.repo.RespondenByID(id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, utils.ErrNotFound("Data responden tidak ditemukan.")
	}
	jawaban, err := s.repo.JawabanResponden(id)
	if err != nil {
		return nil, err
	}
	saran, err := s.repo.SaranResponden(id)
	if err != nil {
		return nil, err
	}
	items := make([]SkmJawabanItem, 0, len(jawaban))
	for _, j := range jawaban {
		items = append(items, SkmJawabanItem{Kode: j.Kode, Unsur: j.Unsur, Nilai: j.Nilai})
	}
	return &SkmRespondenDetail{
		SkmRespondenItem: skmRespondenItem(*r),
		Kecamatan:        r.Kecamatan, Frekuensi: r.Frekuensi, Apresiasi: r.Apresiasi, SkalaMaks: r.SkalaMaks,
		Indeks:  utils.SkmBulat(r.RataRata*100/float64(r.SkalaMaks), 2),
		Jawaban: items, Saran: saran,
	}, nil
}

// ---------- Bank pertanyaan ----------

func (s *SkmServices) DaftarPertanyaan() ([]models.SkmPertanyaan, error) {
	return s.repo.SemuaPertanyaan()
}

func skmTeksPertanyaan(p *string, field string, maks int, wajib bool) (*string, error) {
	if p == nil {
		if wajib {
			return nil, skmGagal(field, "Wajib diisi.")
		}
		return nil, nil
	}
	v := utils.SkmBersih(*p)
	if v == "" {
		return nil, skmGagal(field, "Wajib diisi.")
	}
	if utf8.RuneCountInString(v) > maks {
		return nil, skmGagal(field, fmt.Sprintf("Maksimal %d karakter.", maks))
	}
	return &v, nil
}

func (s *SkmServices) TambahPertanyaan(in SkmPertanyaanInput) (*models.SkmPertanyaan, error) {
	unsur, err := skmTeksPertanyaan(in.Unsur, "unsur", 100, true)
	if err != nil {
		return nil, err
	}
	teks, err := skmTeksPertanyaan(in.Teks, "teks", 500, true)
	if err != nil {
		return nil, err
	}
	p := &models.SkmPertanyaan{Unsur: *unsur, Teks: *teks, Wajib: true, Aktif: true}
	if in.Wajib != nil {
		p.Wajib = *in.Wajib
	}
	if err := s.repo.TambahPertanyaan(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *SkmServices) UbahPertanyaan(id int64, in SkmPertanyaanInput) (*models.SkmPertanyaan, error) {
	p, err := s.repo.PertanyaanByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, utils.ErrNotFound("Pertanyaan tidak ditemukan.")
	}
	updates := map[string]interface{}{}
	if v, err := skmTeksPertanyaan(in.Unsur, "unsur", 100, false); err != nil {
		return nil, err
	} else if v != nil {
		updates["unsur"] = *v
	}
	if v, err := skmTeksPertanyaan(in.Teks, "teks", 500, false); err != nil {
		return nil, err
	} else if v != nil {
		updates["teks"] = *v
	}
	if in.Wajib != nil {
		updates["wajib"] = *in.Wajib
	}
	if in.Aktif != nil {
		updates["aktif"] = *in.Aktif
	}
	if len(updates) > 0 {
		if err := s.repo.UbahPertanyaan(id, updates); err != nil {
			return nil, err
		}
	}
	return s.repo.PertanyaanByID(id)
}

// HapusPertanyaan: soft delete, jawaban lama tetap terhubung untuk rekap periode sebelumnya.
func (s *SkmServices) HapusPertanyaan(id int64) error {
	p, err := s.repo.PertanyaanByID(id)
	if err != nil {
		return err
	}
	if p == nil {
		return utils.ErrNotFound("Pertanyaan tidak ditemukan.")
	}
	return s.repo.HapusPertanyaan(id)
}

// UrutkanPertanyaan: ids harus memuat seluruh pertanyaan (yang belum dihapus) tepat satu kali.
func (s *SkmServices) UrutkanPertanyaan(ids []int64) ([]models.SkmPertanyaan, error) {
	semua, err := s.repo.SemuaPertanyaan()
	if err != nil {
		return nil, err
	}
	ada := map[int64]bool{}
	for _, p := range semua {
		ada[p.ID] = true
	}
	dilihat := map[int64]bool{}
	for _, id := range ids {
		if !ada[id] || dilihat[id] {
			dilihat = nil
			break
		}
		dilihat[id] = true
	}
	if len(ids) == 0 || dilihat == nil || len(dilihat) != len(semua) {
		return nil, skmGagal("ids", "Daftar urutan harus memuat seluruh pertanyaan tepat satu kali.")
	}
	if err := s.repo.UrutkanPertanyaan(ids); err != nil {
		return nil, err
	}
	return s.repo.SemuaPertanyaan()
}

// ---------- Feedback ----------

func skmFeedbackItem(f repositories.SkmFeedbackRow) SkmFeedbackItem {
	item := SkmFeedbackItem{
		ID: f.ID, Nomor: fmt.Sprintf("FB-%04d", f.ID), Tanggal: utils.SkmISO(f.CreatedAt), Isi: f.Isi,
		Status: f.Status, Layanan: f.Layanan, NomorResponden: f.NomorResponden, Catatan: f.Catatan,
		DitanganiOleh: f.Penangan, Diperbarui: utils.SkmISO(f.UpdatedAt),
	}
	if b, ok := models.SkmFeedbackBerikut[f.Status]; ok {
		item.StatusBerikut = &b
	}
	if f.KategoriID != nil && f.Kategori != nil {
		item.Kategori = &SkmKategori{ID: *f.KategoriID, Nama: *f.Kategori}
	}
	return item
}

func (s *SkmServices) DaftarFeedback(status, search string, limit, offset int) ([]SkmFeedbackItem, int64, error) {
	if status != "" && !utils.SkmAda(utils.SkmStatusFeedback, status) {
		return nil, 0, skmGagal("status", "Status tidak valid.")
	}
	if utf8.RuneCountInString(search) > 100 {
		return nil, 0, skmGagal("search", "Kata kunci maksimal 100 karakter.")
	}
	rows, total, err := s.repo.DaftarFeedback(status, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out := make([]SkmFeedbackItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, skmFeedbackItem(r))
	}
	return out, total, nil
}

// HitungStatusFeedback: hitungan per status untuk kartu filter di halaman Feedback.
func (s *SkmServices) HitungStatusFeedback() (map[string]int64, error) {
	hitung, err := s.repo.HitungStatusFeedback()
	if err != nil {
		return nil, err
	}
	out := map[string]int64{}
	for _, st := range utils.SkmStatusFeedback {
		out[st] = hitung[st]
	}
	return out, nil
}

func (s *SkmServices) DetailFeedback(id int64) (*SkmFeedbackItem, error) {
	f, err := s.repo.FeedbackByID(id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, utils.ErrNotFound("Feedback tidak ditemukan.")
	}
	item := skmFeedbackItem(*f)
	return &item, nil
}

// UbahFeedback: status hanya boleh maju satu langkah (Baru → Ditinjau → Ditindaklanjuti → Selesai).
func (s *SkmServices) UbahFeedback(id int64, in SkmFeedbackInput, olehUserID int64) (*SkmFeedbackItem, error) {
	f, err := s.repo.FeedbackByID(id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, utils.ErrNotFound("Feedback tidak ditemukan.")
	}
	updates := map[string]interface{}{"ditangani_oleh": olehUserID}
	if in.Status != nil && *in.Status != f.Status {
		if !utils.SkmAda(utils.SkmStatusFeedback, *in.Status) {
			return nil, skmGagal("status", "Status tidak valid.")
		}
		berikut, ok := models.SkmFeedbackBerikut[f.Status]
		if !ok {
			return nil, skmGagal("status", "Feedback yang sudah selesai tidak dapat diubah statusnya.")
		}
		if *in.Status != berikut {
			return nil, skmGagal("status", "Status "+f.Status+" hanya dapat diubah menjadi "+berikut+".")
		}
		updates["status"] = *in.Status
	}
	if in.CatatanAda {
		if in.Catatan == nil {
			updates["catatan"] = nil
		} else {
			v := utils.SkmBersih(*in.Catatan)
			if utf8.RuneCountInString(v) > 2000 {
				return nil, skmGagal("catatan", "Catatan maksimal 2000 karakter.")
			}
			if v == "" {
				updates["catatan"] = nil
			} else {
				updates["catatan"] = v
			}
		}
	}
	if in.KategoriIDAda {
		if in.KategoriID == nil {
			updates["kategori_id"] = nil
		} else {
			ada, err := s.repo.PertanyaanByIDs([]int64{*in.KategoriID})
			if err != nil {
				return nil, err
			}
			if len(ada) == 0 {
				return nil, skmGagal("kategori_id", "Kategori tidak ditemukan.")
			}
			updates["kategori_id"] = *in.KategoriID
		}
	}
	if err := s.repo.UbahFeedback(id, updates); err != nil {
		return nil, err
	}
	return s.DetailFeedback(id)
}

// ---------- Pesan kontak ----------

func skmPesanItem(p models.SkmPesanKontak) SkmPesanItem {
	return SkmPesanItem{ID: p.ID, Nama: p.Nama, Email: p.Email, Topik: p.Topik, Pesan: p.Pesan, Dibaca: p.Dibaca, Tanggal: utils.SkmISO(p.CreatedAt)}
}

func (s *SkmServices) DaftarPesan(dibaca *bool, limit, offset int) ([]SkmPesanItem, int64, error) {
	rows, total, err := s.repo.DaftarPesan(dibaca, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out := make([]SkmPesanItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, skmPesanItem(r))
	}
	return out, total, nil
}

func (s *SkmServices) JumlahPesanBelumDibaca() (int64, error) {
	return s.repo.JumlahPesanBelumDibaca()
}

func (s *SkmServices) UbahPesan(id int64, dibaca bool) (*SkmPesanItem, error) {
	p, err := s.repo.PesanByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, utils.ErrNotFound("Pesan tidak ditemukan.")
	}
	if err := s.repo.UbahPesan(id, dibaca); err != nil {
		return nil, err
	}
	p.Dibaca = dibaca
	item := skmPesanItem(*p)
	return &item, nil
}

// ---------- Pengaturan ----------

func (s *SkmServices) Pengaturan() (*SkmPengaturanView, error) {
	p, err := s.repo.Pengaturan()
	if err != nil {
		return nil, err
	}
	terkunci, err := s.repo.PeriodePunyaResponden(p.PeriodeAktifID)
	if err != nil {
		return nil, err
	}
	periode, err := s.repo.DaftarPeriode()
	if err != nil {
		return nil, err
	}
	pengelola, err := s.repo.Pengelola()
	if err != nil {
		return nil, err
	}
	opsi := make([]int, 0, len(periode))
	for _, x := range periode {
		opsi = append(opsi, x.Tahun)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(opsi)))
	list := make([]SkmPengelolaItem, 0, len(pengelola))
	for _, u := range pengelola {
		nama, jabatan := "", ""
		if u.Name != nil {
			nama = *u.Name
		}
		if u.Position != nil {
			jabatan = *u.Position
		}
		list = append(list, SkmPengelolaItem{ID: u.PublicID, Nama: nama, Jabatan: jabatan, Inisial: utils.SkmInisial(nama)})
	}
	return &SkmPengaturanView{
		Judul: p.Judul, Periode: p.PeriodeAktif.Tahun, SkalaMaks: p.PeriodeAktif.SkalaMaks, SkalaTerkunci: terkunci,
		TanggalTutup: skmTanggal(p.PeriodeAktif.TanggalTutup), Anonim: p.Anonim, Publik: p.Publik,
		OpsiPeriode: opsi, Pengelola: list,
	}, nil
}

func (s *SkmServices) SimpanPengaturan(in SkmPengaturanInput) (*SkmPengaturanView, error) {
	judul := utils.SkmBersih(in.Judul)
	switch {
	case judul == "":
		return nil, skmGagal("judul", "Judul survei wajib diisi.")
	case utf8.RuneCountInString(judul) > 150:
		return nil, skmGagal("judul", "Judul maksimal 150 karakter.")
	case in.Periode < 2000 || in.Periode > 2100:
		return nil, skmGagal("periode", "Periode harus tahun antara 2000 dan 2100.")
	case in.SkalaMaks != 4 && in.SkalaMaks != 5:
		return nil, skmGagal("skala_maks", "Skala harus 4 atau 5.")
	case in.Anonim == nil:
		return nil, skmGagal("anonim", "Wajib diisi.")
	case in.Publik == nil:
		return nil, skmGagal("publik", "Wajib diisi.")
	}
	var tutup *time.Time
	if in.TanggalTutup != nil && *in.TanggalTutup != "" {
		t, err := time.Parse("2006-01-02", *in.TanggalTutup)
		if err != nil {
			return nil, skmGagal("tanggal_tutup", "Format tanggal harus YYYY-MM-DD.")
		}
		tutup = &t
	}
	err := s.repo.SimpanPengaturan(repositories.SkmPengaturanBaru{
		Judul: judul, Tahun: in.Periode, SkalaMaks: in.SkalaMaks, TanggalTutup: tutup,
		Anonim: *in.Anonim, Publik: *in.Publik,
	})
	if errors.Is(err, repositories.ErrSkmSkalaTerkunci) {
		return nil, utils.ErrConflict(fmt.Sprintf("Skala penilaian periode %d tidak dapat diubah karena sudah ada responden.", in.Periode))
	}
	if err != nil {
		return nil, err
	}
	return s.Pengaturan()
}
