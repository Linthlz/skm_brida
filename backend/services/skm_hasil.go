package services

import (
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

// Perhitungan agregat SKM. Indeks = rata-rata NRR unsur × (100 / skala_maks),
// sesuai Permen PANRB 14/2017 (bobot tiap unsur sama).

var itoa = strconv.Itoa

// Periode: tahun yang diminta, atau periode aktif bila tahun = 0.
func (s *SkmServices) Periode(tahun int) (*models.SkmPeriode, error) {
	if tahun == 0 {
		p, err := s.repo.Pengaturan()
		if err != nil {
			return nil, err
		}
		return &p.PeriodeAktif, nil
	}
	p, err := s.repo.PeriodeByTahun(tahun)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, skmGagal("periode", "Periode tidak ditemukan.")
	}
	return p, nil
}

func (s *SkmServices) Filter(layananID int64, pendidikan string) (repositories.SkmFilter, error) {
	f := repositories.SkmFilter{LayananID: layananID, Pendidikan: pendidikan}
	if layananID != 0 {
		l, err := s.repo.LayananByID(layananID)
		if err != nil {
			return f, err
		}
		if l == nil {
			return f, skmGagal("layanan_id", "Layanan tidak ditemukan.")
		}
	}
	if pendidikan != "" && !utils.SkmAda(utils.SkmPendidikan, pendidikan) {
		return f, skmGagal("pendidikan", "Pendidikan tidak valid.")
	}
	return f, nil
}

// unsur: NRR tiap unsur yang punya jawaban pada periode itu, urut sesuai bank pertanyaan.
func (s *SkmServices) unsur(p *models.SkmPeriode, f repositories.SkmFilter) ([]SkmUnsur, error) {
	agregat, err := s.repo.NRRUnsur(p.ID, f)
	if err != nil {
		return nil, err
	}
	out := make([]SkmUnsur, 0, len(agregat))
	if len(agregat) == 0 {
		return out, nil
	}
	per := map[int64]repositories.SkmNRR{}
	ids := make([]int64, 0, len(agregat))
	for _, a := range agregat {
		per[a.PertanyaanID] = a
		ids = append(ids, a.PertanyaanID)
	}
	pertanyaan, err := s.repo.PertanyaanByIDs(ids)
	if err != nil {
		return nil, err
	}
	for _, q := range pertanyaan {
		a := per[q.ID]
		out = append(out, SkmUnsur{ID: q.ID, Kode: q.Kode, Nama: q.Unsur, NRR: a.NRR, JumlahJawaban: a.N})
	}
	return out, nil
}

func skmRataNRR(u []SkmUnsur) *float64 {
	if len(u) == 0 {
		return nil
	}
	total := 0.0
	for _, x := range u {
		total += x.NRR
	}
	v := total / float64(len(u))
	return &v
}

func skmKali(n *float64, k float64) *float64 {
	if n == nil {
		return nil
	}
	v := *n * k
	return &v
}

// Rekap: rekap utama untuk halaman Hasil (publik & admin) dan Laporan.
func (s *SkmServices) Rekap(p *models.SkmPeriode, f repositories.SkmFilter) (*SkmRekap, error) {
	unsur, err := s.unsur(p, f)
	if err != nil {
		return nil, err
	}
	nrr := skmRataNRR(unsur)
	ikm := skmKali(nrr, p.Pengali())
	total, err := s.repo.JumlahResponden(p.ID, f)
	if err != nil {
		return nil, err
	}
	rentang, err := s.rentangBulan(p, f)
	if err != nil {
		return nil, err
	}
	distribusi, err := s.distribusi(p, f)
	if err != nil {
		return nil, err
	}
	tren, err := s.tren(f)
	if err != nil {
		return nil, err
	}
	for i := range unsur {
		u := &unsur[i]
		indeks := u.NRR * p.Pengali()
		u.NRRTertimbang = utils.SkmBulat(u.NRR/float64(len(unsur)), 3)
		u.Indeks = utils.SkmBulat(indeks, 2)
		u.Mutu = utils.SkmKategoriMutu(&indeks)
		u.NRR = utils.SkmBulat(u.NRR, 2)
	}
	return &SkmRekap{
		Periode: p.Tahun, SkalaMaks: p.SkalaMaks, Pengali: p.Pengali(), TotalResponden: total,
		RentangBulan: rentang, NRR: utils.SkmBulatPtr(nrr, 2), IKM: utils.SkmBulatPtr(ikm, 2),
		Mutu: utils.SkmKategoriMutu(ikm), Unsur: unsur, Distribusi: distribusi, Tren: tren,
	}, nil
}

// Ringkasan: angka ringkas untuk Hero, Tentang, dan panel login.
func (s *SkmServices) Ringkasan() (*SkmRingkasan, error) {
	pengaturan, err := s.repo.Pengaturan()
	if err != nil {
		return nil, err
	}
	p := &pengaturan.PeriodeAktif
	unsur, err := s.unsur(p, repositories.SkmFilter{})
	if err != nil {
		return nil, err
	}
	ikm := skmKali(skmRataNRR(unsur), p.Pengali())
	total, err := s.repo.JumlahResponden(p.ID, repositories.SkmFilter{})
	if err != nil {
		return nil, err
	}
	feedback, err := s.repo.JumlahFeedbackPeriode(p.ID)
	if err != nil {
		return nil, err
	}
	aktif, err := s.repo.PertanyaanAktif()
	if err != nil {
		return nil, err
	}
	layanan, err := s.repo.JumlahLayananAktif()
	if err != nil {
		return nil, err
	}
	return &SkmRingkasan{
		Periode: p.Tahun, IKM: utils.SkmBulatPtr(ikm, 2), Mutu: utils.SkmKategoriMutu(ikm),
		TotalResponden: total, TotalFeedback: feedback, JumlahUnsur: len(aktif), JumlahLayanan: layanan,
	}, nil
}

func (s *SkmServices) Dashboard(p *models.SkmPeriode) (*SkmDashboard, error) {
	rekap, err := s.Rekap(p, repositories.SkmFilter{})
	if err != nil {
		return nil, err
	}
	now := time.Now().In(utils.SkmZona)
	awalBulan := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, utils.SkmZona)
	bulanIni, err := s.repo.RespondenSejak(p.ID, awalBulan)
	if err != nil {
		return nil, err
	}
	lengkap, err := s.repo.JumlahLengkap(p.ID)
	if err != nil {
		return nil, err
	}
	var persen *int64
	if rekap.TotalResponden > 0 {
		v := int64(math.Round(float64(lengkap) / float64(rekap.TotalResponden) * 100))
		persen = &v
	}
	feedbackTotal, err := s.repo.JumlahFeedbackPeriode(p.ID)
	if err != nil {
		return nil, err
	}
	hitung, err := s.repo.HitungStatusFeedback()
	if err != nil {
		return nil, err
	}
	belumDibaca, err := s.repo.JumlahPesanBelumDibaca()
	if err != nil {
		return nil, err
	}
	perBulan, err := s.perBulan(p)
	if err != nil {
		return nil, err
	}
	prioritas := append([]SkmUnsur(nil), rekap.Unsur...)
	sort.SliceStable(prioritas, func(i, j int) bool { return prioritas[i].NRR < prioritas[j].NRR })
	if len(prioritas) > 3 {
		prioritas = prioritas[:3]
	}
	if prioritas == nil {
		prioritas = []SkmUnsur{}
	}
	return &SkmDashboard{
		Periode: p.Tahun, SkalaMaks: p.SkalaMaks, Pengali: p.Pengali(), TotalResponden: rekap.TotalResponden,
		BulanIni: bulanIni, PersenLengkap: persen, IKM: rekap.IKM, Mutu: rekap.Mutu,
		FeedbackTotal: feedbackTotal, FeedbackBaru: hitung[models.SkmFeedbackBaru], PesanBelumDibaca: belumDibaca,
		PerBulan: perBulan, Distribusi: rekap.Distribusi, Unsur: rekap.Unsur, Prioritas: prioritas,
	}, nil
}

func (s *SkmServices) Statistik(p *models.SkmPeriode) (*SkmStatistik, error) {
	total, err := s.repo.JumlahResponden(p.ID, repositories.SkmFilter{})
	if err != nil {
		return nil, err
	}
	pendidikan, err := s.hitungPer(p, "pendidikan", utils.SkmPendidikan)
	if err != nil {
		return nil, err
	}
	pekerjaan, err := s.hitungPer(p, "pekerjaan", utils.SkmPekerjaan)
	if err != nil {
		return nil, err
	}
	tren, err := s.tren(repositories.SkmFilter{})
	if err != nil {
		return nil, err
	}
	return &SkmStatistik{Periode: p.Tahun, TotalResponden: total, Pendidikan: pendidikan, Pekerjaan: pekerjaan, Tren: tren}, nil
}

// distribusi: sebaran responden per tingkat kepuasan, dari tingkat tertinggi ke terendah.
func (s *SkmServices) distribusi(p *models.SkmPeriode, f repositories.SkmFilter) ([]SkmDistribusi, error) {
	hitung, err := s.repo.Distribusi(p.ID, f)
	if err != nil {
		return nil, err
	}
	skala := utils.SkmDaftarSkala(p.SkalaMaks)
	out := make([]SkmDistribusi, 0, len(skala))
	for i := len(skala) - 1; i >= 0; i-- {
		out = append(out, SkmDistribusi{V: skala[i].V, Label: skala[i].Label, Jumlah: hitung[skala[i].V]})
	}
	return out, nil
}

// tren: IKM tiap periode yang punya data, urut tahun naik.
func (s *SkmServices) tren(f repositories.SkmFilter) ([]SkmTren, error) {
	rows, err := s.repo.NRRPerPeriode(f)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(rows))
	for id := range rows {
		ids = append(ids, id)
	}
	periode, err := s.repo.PeriodeByIDs(ids)
	if err != nil {
		return nil, err
	}
	out := make([]SkmTren, 0, len(periode))
	for _, p := range periode {
		total := 0.0
		for _, n := range rows[p.ID] {
			total += n
		}
		ikm := utils.SkmBulat(total/float64(len(rows[p.ID]))*p.Pengali(), 2)
		out = append(out, SkmTren{Periode: itoa(p.Tahun), IKM: &ikm})
	}
	return out, nil
}

// perBulan: jumlah responden per bulan, dari Januari sampai bulan terakhir yang berisi data.
func (s *SkmServices) perBulan(p *models.SkmPeriode) ([]SkmPerBulan, error) {
	hitung, err := s.repo.PerBulan(p.ID)
	if err != nil {
		return nil, err
	}
	akhir := 1
	for b := range hitung {
		if b > akhir {
			akhir = b
		}
	}
	out := make([]SkmPerBulan, 0, akhir)
	for b := 1; b <= akhir; b++ {
		out = append(out, SkmPerBulan{Bulan: b, Label: utils.SkmBulan[b-1], Jumlah: hitung[b]})
	}
	return out, nil
}

// rentangBulan: bulan pertama & terakhir pengisian, contoh "Jan – Sep".
func (s *SkmServices) rentangBulan(p *models.SkmPeriode, f repositories.SkmFilter) (*string, error) {
	awal, akhir, err := s.repo.RentangWaktu(p.ID, f)
	if err != nil || awal == nil || akhir == nil {
		return nil, err
	}
	a := utils.SkmBulan[awal.In(utils.SkmZona).Month()-1]
	b := utils.SkmBulan[akhir.In(utils.SkmZona).Month()-1]
	out := a
	if a != b {
		out = a + " – " + b
	}
	return &out, nil
}

func (s *SkmServices) hitungPer(p *models.SkmPeriode, kolom string, opsi []string) ([]SkmHitung, error) {
	hitung, err := s.repo.HitungPer(p.ID, kolom)
	if err != nil {
		return nil, err
	}
	out := make([]SkmHitung, 0, len(opsi))
	for _, o := range opsi {
		out = append(out, SkmHitung{Label: o, N: hitung[o]})
	}
	return out, nil
}
