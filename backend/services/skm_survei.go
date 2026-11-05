package services

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

func (s *SkmServices) Referensi() (*SkmReferensi, error) {
	layanan, err := s.repo.LayananAktif()
	if err != nil {
		return nil, err
	}
	return &SkmReferensi{
		Layanan: layanan, JK: utils.SkmJK, Pendidikan: utils.SkmPendidikan, Pekerjaan: utils.SkmPekerjaan,
		Frekuensi: utils.SkmFrekuensi, Kecamatan: utils.SkmKecamatan, TopikKontak: utils.SkmTopikKontak,
		StatusFeedback: utils.SkmStatusFeedback,
	}, nil
}

// Form: data halaman /survei — pengaturan, skala, dan pertanyaan aktif.
func (s *SkmServices) Form() (*SkmForm, error) {
	p, err := s.repo.Pengaturan()
	if err != nil {
		return nil, err
	}
	pertanyaan, err := s.repo.PertanyaanAktif()
	if err != nil {
		return nil, err
	}
	out := &SkmForm{
		Judul:  p.Judul,
		Anonim: p.Anonim,
		Periode: SkmFormPeriode{
			Tahun: p.PeriodeAktif.Tahun, SkalaMaks: p.PeriodeAktif.SkalaMaks,
			TanggalTutup: skmTanggal(p.PeriodeAktif.TanggalTutup), Ditutup: skmDitutup(p.PeriodeAktif),
		},
		Skala:      utils.SkmDaftarSkala(p.PeriodeAktif.SkalaMaks),
		Pertanyaan: make([]SkmFormPertanyaan, 0, len(pertanyaan)),
	}
	for _, q := range pertanyaan {
		out.Pertanyaan = append(out.Pertanyaan, SkmFormPertanyaan{ID: q.ID, Kode: q.Kode, Unsur: q.Unsur, Teks: q.Teks, Wajib: q.Wajib})
	}
	return out, nil
}

func skmOpsional(p *string, field string, maks int) (*string, error) {
	if p == nil {
		return nil, nil
	}
	v := utils.SkmBersih(*p)
	if v == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(v) > maks {
		return nil, skmGagal(field, "Maksimal "+itoa(maks)+" karakter.")
	}
	return &v, nil
}

// Kirim menyimpan satu pengisian survei dari periode aktif.
func (s *SkmServices) Kirim(in SkmKirimSurvei, ip string) (*SkmHasilKirim, error) {
	pengaturan, err := s.repo.Pengaturan()
	if err != nil {
		return nil, err
	}
	periode := pengaturan.PeriodeAktif
	if skmDitutup(periode) {
		return nil, utils.ErrForbidden("Periode survei sudah ditutup. Terima kasih atas minat Anda.")
	}

	if in.Website != "" {
		return nil, skmGagal("website", "Isian tidak valid.")
	}
	nama, err := skmOpsional(in.Nama, "nama", 100)
	if err != nil {
		return nil, err
	}
	if nama == nil && !pengaturan.Anonim {
		return nil, skmGagal("nama", "Nama wajib diisi.")
	}
	switch {
	case !utils.SkmAda(utils.SkmJK, in.JK):
		return nil, skmGagal("jk", "Pilih jenis kelamin.")
	case in.Usia < 10 || in.Usia > 99:
		return nil, skmGagal("usia", "Usia harus antara 10 dan 99 tahun.")
	case !utils.SkmAda(utils.SkmPendidikan, in.Pendidikan):
		return nil, skmGagal("pendidikan", "Pilih pendidikan terakhir.")
	case !utils.SkmAda(utils.SkmPekerjaan, in.Pekerjaan):
		return nil, skmGagal("pekerjaan", "Pilih pekerjaan.")
	case !utils.SkmAda(utils.SkmFrekuensi, in.Frekuensi):
		return nil, skmGagal("frekuensi", "Pilih frekuensi penggunaan layanan.")
	}
	kecamatan, _ := skmOpsional(in.Kecamatan, "kecamatan", 40)
	if kecamatan != nil && !utils.SkmAda(utils.SkmKecamatan, *kecamatan) {
		return nil, skmGagal("kecamatan", "Kecamatan tidak valid.")
	}
	layanan, err := s.repo.LayananByID(in.LayananID)
	if err != nil {
		return nil, err
	}
	if layanan == nil || !layanan.Aktif {
		return nil, skmGagal("layanan_id", "Pilih jenis layanan yang diterima.")
	}
	saran, err := skmOpsional(in.Saran, "saran", 600)
	if err != nil {
		return nil, err
	}
	apresiasi, err := skmOpsional(in.Apresiasi, "apresiasi", 300)
	if err != nil {
		return nil, err
	}
	if !in.Setuju {
		return nil, skmGagal("setuju", "Anda perlu menyetujui pernyataan sebelum mengirim.")
	}

	// Jawaban: pertanyaan aktif, tidak ganda, nilai dalam skala, dan semua yang wajib terisi.
	aktif, err := s.repo.PertanyaanAktif()
	if err != nil {
		return nil, err
	}
	if len(in.Jawaban) == 0 {
		return nil, skmGagal("jawaban", "Silakan pilih salah satu jawaban sebelum melanjutkan.")
	}
	urutan := map[int64]int{}
	for i, q := range aktif {
		urutan[q.ID] = i
	}
	nilai := map[int64]int{}
	for _, j := range in.Jawaban {
		if _, ok := urutan[j.PertanyaanID]; !ok {
			return nil, skmGagal("jawaban", "Pertanyaan tidak valid. Muat ulang halaman lalu coba lagi.")
		}
		if _, ganda := nilai[j.PertanyaanID]; ganda {
			return nil, skmGagal("jawaban", "Setiap pertanyaan hanya boleh dijawab sekali.")
		}
		if j.Nilai < 1 || j.Nilai > periode.SkalaMaks {
			return nil, skmGagal("jawaban", "Nilai harus antara 1 dan "+itoa(periode.SkalaMaks)+".")
		}
		nilai[j.PertanyaanID] = j.Nilai
	}
	belum := 0
	for _, q := range aktif {
		if _, ok := nilai[q.ID]; q.Wajib && !ok {
			belum++
		}
	}
	if belum > 0 {
		return nil, skmGagal("jawaban", "Silakan pilih salah satu jawaban sebelum melanjutkan. "+itoa(belum)+" pertanyaan belum dinilai.")
	}

	total := 0
	jawaban := make([]models.SkmJawaban, 0, len(nilai))
	// Kategori awal feedback = unsur dengan nilai terendah; bila sama, urutan paling awal.
	var terendah *models.SkmPertanyaan
	for i := range aktif {
		q := &aktif[i]
		v, ok := nilai[q.ID]
		if !ok {
			continue
		}
		total += v
		jawaban = append(jawaban, models.SkmJawaban{PertanyaanID: q.ID, Nilai: v})
		if terendah == nil || v < nilai[terendah.ID] {
			terendah = q
		}
	}
	rata := utils.SkmBulat(float64(total)/float64(len(jawaban)), 2)

	resp := &models.SkmResponden{
		PeriodeID: periode.ID, LayananID: layanan.ID, Nama: nama, JK: in.JK, Usia: in.Usia,
		Pendidikan: in.Pendidikan, Pekerjaan: in.Pekerjaan, Frekuensi: in.Frekuensi,
		Kecamatan: kecamatan, Apresiasi: apresiasi, RataRata: rata, IPHash: utils.SkmIPHash(ip),
	}
	var fb *models.SkmFeedback
	if saran != nil {
		fb = &models.SkmFeedback{Isi: *saran, KategoriID: &terendah.ID, Status: models.SkmFeedbackBaru}
	}
	if err := s.repo.SimpanSurvei(resp, jawaban, fb); err != nil {
		return nil, err
	}
	indeks := rata * periode.Pengali()
	return &SkmHasilKirim{
		Nomor: resp.Nomor, Rata: rata, SkalaMaks: periode.SkalaMaks,
		Indeks: utils.SkmBulat(indeks, 2), Mutu: utils.SkmKategoriMutu(&indeks),
	}, nil
}

var skmPolaEmail = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[A-Za-z]{2,}$`)

func (s *SkmServices) KirimKontak(in SkmKirimKontak, ip string) error {
	in.Nama = utils.SkmBersih(in.Nama)
	in.Email = strings.TrimSpace(in.Email)
	in.Pesan = utils.SkmBersih(in.Pesan)
	switch {
	case in.Website != "":
		return skmGagal("website", "Isian tidak valid.")
	case in.Nama == "" || utf8.RuneCountInString(in.Nama) > 100:
		return skmGagal("nama", "Nama wajib diisi (maksimal 100 karakter).")
	case len(in.Email) > 150 || !skmPolaEmail.MatchString(in.Email):
		return skmGagal("email", "Masukkan alamat email yang valid, contoh: nama@email.com")
	case !utils.SkmAda(utils.SkmTopikKontak, in.Topik):
		return skmGagal("topik", "Pilih topik pesan.")
	case utf8.RuneCountInString(in.Pesan) < 10:
		return skmGagal("pesan", "Tuliskan pesan minimal 10 karakter agar kami dapat menindaklanjuti.")
	case utf8.RuneCountInString(in.Pesan) > 2000:
		return skmGagal("pesan", "Pesan maksimal 2000 karakter.")
	}
	return s.repo.SimpanPesan(&models.SkmPesanKontak{
		Nama: in.Nama, Email: in.Email, Topik: in.Topik, Pesan: in.Pesan, IPHash: utils.SkmIPHash(ip),
	})
}

func (s *SkmServices) DaftarPeriodePublik() ([]SkmPeriodePublik, error) {
	p, err := s.repo.Pengaturan()
	if err != nil {
		return nil, err
	}
	list, err := s.repo.DaftarPeriode()
	if err != nil {
		return nil, err
	}
	out := make([]SkmPeriodePublik, 0, len(list))
	for _, x := range list {
		out = append(out, SkmPeriodePublik{Tahun: x.Tahun, Aktif: x.ID == p.PeriodeAktifID})
	}
	return out, nil
}

// PastikanPublik: hasil agregat publik ditolak bila publikasi dimatikan di pengaturan.
func (s *SkmServices) PastikanPublik() error {
	p, err := s.repo.Pengaturan()
	if err != nil {
		return err
	}
	if !p.Publik {
		return utils.ErrForbidden("Hasil survei belum dipublikasikan.")
	}
	return nil
}
