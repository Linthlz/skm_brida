package services

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
	"github.com/xuri/excelize/v2"
)

// Laporan: rekapitulasi nilai per unsur, sama dengan tabel di halaman Laporan.
func (s *SkmServices) Laporan(p *models.SkmPeriode, f repositories.SkmFilter, format string) (*SkmBerkas, error) {
	if format != "xlsx" && format != "pdf" {
		return nil, skmGagal("format", "Format harus xlsx atau pdf.")
	}
	rekap, err := s.Rekap(p, f)
	if err != nil {
		return nil, err
	}
	keterangan, err := s.keteranganFilter(f)
	if err != nil {
		return nil, err
	}
	nama := fmt.Sprintf("laporan-skm-%d.%s", p.Tahun, format)
	if format == "xlsx" {
		isi, err := skmXlsx(rekap, keterangan)
		if err != nil {
			return nil, err
		}
		return &SkmBerkas{Nama: nama, ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", Isi: isi}, nil
	}
	pengaturan, err := s.repo.Pengaturan()
	if err != nil {
		return nil, err
	}
	isi, err := skmPdf(rekap, keterangan, pengaturan.Judul)
	if err != nil {
		return nil, err
	}
	return &SkmBerkas{Nama: nama, ContentType: "application/pdf", Isi: isi}, nil
}

func (s *SkmServices) keteranganFilter(f repositories.SkmFilter) (string, error) {
	layanan := "Semua layanan"
	if f.LayananID != 0 {
		l, err := s.repo.LayananByID(f.LayananID)
		if err != nil {
			return "", err
		}
		if l != nil {
			layanan = l.Nama
		}
	}
	pendidikan := "Semua jenjang pendidikan"
	if f.Pendidikan != "" {
		pendidikan = f.Pendidikan
	}
	return layanan + " · " + pendidikan, nil
}

func skmNrrTotalTertimbang(r *SkmRekap) *float64 {
	if r.NRR == nil || len(r.Unsur) == 0 {
		return nil
	}
	v := utils.SkmBulat(*r.NRR/float64(len(r.Unsur)), 3)
	return &v
}

func skmHuruf(m *utils.SkmMutu) string {
	if m == nil {
		return ""
	}
	return m.Huruf
}

func skmXlsx(r *SkmRekap, keterangan string) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	sheet := fmt.Sprintf("SKM %d", r.Periode)
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return nil, err
	}
	baris := [][]interface{}{
		{"Rekapitulasi Survei Kepuasan Masyarakat — periode " + itoa(r.Periode)},
		{fmt.Sprintf("%s · %d responden", keterangan, r.TotalResponden)},
		{},
		{"Kode", "Unsur pelayanan", "NRR", "NRR tertimbang", "Indeks", "Mutu"},
	}
	for _, u := range r.Unsur {
		baris = append(baris, []interface{}{u.Kode, u.Nama, u.NRR, u.NRRTertimbang, u.Indeks, skmHuruf(u.Mutu)})
	}
	total := []interface{}{"Total", "", nil, nil, nil, skmHuruf(r.Mutu)}
	if r.NRR != nil {
		total[2] = *r.NRR
	}
	if t := skmNrrTotalTertimbang(r); t != nil {
		total[3] = *t
	}
	if r.IKM != nil {
		total[4] = *r.IKM
	}
	baris = append(baris, total)
	for i, row := range baris {
		sel, _ := excelize.CoordinatesToCellName(1, i+1)
		if err := f.SetSheetRow(sheet, sel, &row); err != nil {
			return nil, err
		}
	}
	tebal, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	_ = f.SetCellStyle(sheet, "A1", "A1", tebal)
	_ = f.SetCellStyle(sheet, "A4", "F4", tebal)
	akhir := fmt.Sprintf("F%d", len(baris))
	_ = f.SetCellStyle(sheet, fmt.Sprintf("A%d", len(baris)), akhir, tebal)
	_ = f.SetColWidth(sheet, "A", "A", 8)
	_ = f.SetColWidth(sheet, "B", "B", 34)
	_ = f.SetColWidth(sheet, "C", "F", 15)
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// skmAngka: format angka Indonesia, contoh 1.234,56.
func skmAngka(n float64, d int) string {
	s := fmt.Sprintf("%.*f", d, n)
	bulat, pecahan := s, ""
	if i := strings.IndexByte(s, '.'); i >= 0 {
		bulat, pecahan = s[:i], s[i+1:]
	}
	minus := strings.HasPrefix(bulat, "-")
	bulat = strings.TrimPrefix(bulat, "-")
	var b strings.Builder
	for i, c := range bulat {
		if i > 0 && (len(bulat)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(c)
	}
	out := b.String()
	if minus {
		out = "-" + out
	}
	if pecahan != "" {
		out += "," + pecahan
	}
	return out
}

func skmPdf(r *SkmRekap, keterangan, judul string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	pdf.SetFont("Helvetica", "B", 14)
	pdf.MultiCell(0, 7, tr(judul), "", "L", false)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(0x6b, 0x62, 0x5c)
	pdf.MultiCell(0, 5, tr(fmt.Sprintf("Rekapitulasi hasil — periode %d · %s", r.Periode, keterangan)), "", "L", false)
	ringkas := fmt.Sprintf("%s responden · IKM %s", skmAngka(float64(r.TotalResponden), 0), "–")
	if r.IKM != nil {
		ringkas = fmt.Sprintf("%s responden · IKM %s", skmAngka(float64(r.TotalResponden), 0), skmAngka(*r.IKM, 2))
	}
	if r.Mutu != nil {
		ringkas += fmt.Sprintf(" · Mutu %s (%s)", r.Mutu.Huruf, r.Mutu.Label)
	}
	pdf.MultiCell(0, 5, tr(ringkas), "", "L", false)
	pdf.Ln(4)

	lebar := []float64{16, 74, 22, 28, 22, 18}
	kepala := []string{"KODE", "UNSUR PELAYANAN", "NRR", "NRR TERTIMBANG", "INDEKS", "MUTU"}
	pdf.SetTextColor(0x1f, 0x1a, 0x17)
	pdf.SetDrawColor(0xd9, 0xd0, 0xc9)
	pdf.SetFillColor(0xf3, 0xed, 0xe8)
	pdf.SetFont("Helvetica", "B", 8)
	for i, k := range kepala {
		pdf.CellFormat(lebar[i], 7, k, "1", 0, "L", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Helvetica", "", 9)
	if len(r.Unsur) == 0 {
		pdf.CellFormat(180, 7, "Belum ada data pada filter ini.", "1", 1, "L", false, 0, "")
	}
	for _, u := range r.Unsur {
		pdf.CellFormat(lebar[0], 7, tr(u.Kode), "1", 0, "L", false, 0, "")
		pdf.CellFormat(lebar[1], 7, tr(u.Nama), "1", 0, "L", false, 0, "")
		pdf.CellFormat(lebar[2], 7, skmAngka(u.NRR, 2), "1", 0, "R", false, 0, "")
		pdf.CellFormat(lebar[3], 7, skmAngka(u.NRRTertimbang, 3), "1", 0, "R", false, 0, "")
		pdf.CellFormat(lebar[4], 7, skmAngka(u.Indeks, 2), "1", 0, "R", false, 0, "")
		pdf.CellFormat(lebar[5], 7, skmHuruf(u.Mutu), "1", 1, "L", false, 0, "")
	}
	if len(r.Unsur) > 0 && r.NRR != nil && r.IKM != nil {
		pdf.SetFont("Helvetica", "B", 9)
		pdf.SetFillColor(0xf7, 0xec, 0xec)
		pdf.CellFormat(lebar[0]+lebar[1], 7, "Total", "1", 0, "L", true, 0, "")
		pdf.CellFormat(lebar[2], 7, skmAngka(*r.NRR, 2), "1", 0, "R", true, 0, "")
		pdf.CellFormat(lebar[3], 7, skmAngka(*skmNrrTotalTertimbang(r), 3), "1", 0, "R", true, 0, "")
		pdf.CellFormat(lebar[4], 7, skmAngka(*r.IKM, 2), "1", 0, "R", true, 0, "")
		pdf.CellFormat(lebar[5], 7, skmHuruf(r.Mutu), "1", 1, "L", true, 0, "")
	}

	pdf.Ln(4)
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(0x6b, 0x62, 0x5c)
	pengali := strings.TrimSuffix(strings.TrimRight(skmAngka(r.Pengali, 2), "0"), ",")
	if r.Pengali == math.Trunc(r.Pengali) {
		pengali = skmAngka(r.Pengali, 0)
	}
	now := time.Now().In(utils.SkmZona)
	bulan := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	dicetak := fmt.Sprintf("%02d %s %d %02d:%02d", now.Day(), bulan[now.Month()-1], now.Year(), now.Hour(), now.Minute())
	pdf.MultiCell(0, 4, tr(fmt.Sprintf("Skala 1–%d, indeks = NRR × %s. Dicetak %s.", r.SkalaMaks, pengali, dicetak)), "", "L", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
