package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata" // zona Asia/Makassar tetap tersedia di image tanpa tzdata
)

// Daftar pilihan dan aturan domain SKM. Disamakan dengan frontend
// (src/data/skala.js dan src/utils/mutu.js).

var SkmJK = []string{"Laki-laki", "Perempuan"}

var SkmPendidikan = []string{
	"SD / Sederajat", "SMP / Sederajat", "SMA / Sederajat",
	"Diploma (D1–D3)", "Sarjana (S1/D4)", "Pascasarjana (S2/S3)",
}

var SkmPekerjaan = []string{
	"Pelajar / Mahasiswa", "ASN / TNI / Polri", "Pegawai Swasta", "Wiraswasta",
	"Peneliti / Akademisi", "Petani / Nelayan", "Lainnya",
}

var SkmFrekuensi = []string{"Baru pertama kali", "2–3 kali", "4–6 kali", "Lebih dari 6 kali"}

var SkmKecamatan = []string{
	"Buleleng", "Sukasada", "Sawan", "Banjar", "Seririt",
	"Kubutambahan", "Tejakula", "Busungbiu", "Gerokgak",
}

var SkmTopikKontak = []string{
	"Pertanyaan tentang survei", "Permintaan data hasil SKM", "Pengaduan layanan",
	"Kerja sama riset & inovasi", "Lainnya",
}

var SkmStatusFeedback = []string{"Baru", "Ditinjau", "Ditindaklanjuti", "Selesai"}

// Label skala kepuasan per nilai, dikelompokkan menurut nilai maksimum skala.
var skmLabelSkala = map[int][]string{
	5: {"Sangat Tidak Puas", "Tidak Puas", "Kurang Puas", "Puas", "Sangat Puas"},
	4: {"Tidak Puas", "Kurang Puas", "Puas", "Sangat Puas"},
}

// Ambang mutu pelayanan (Permen PANRB 14/2017), batas bawah inklusif.
var skmMutu = []struct {
	Huruf string
	Min   float64
	Label string
}{
	{"A", 88.31, "Sangat Baik"},
	{"B", 76.61, "Baik"},
	{"C", 65.00, "Kurang Baik"},
	{"D", 0, "Tidak Baik"},
}

// SkmZona: tanggal tutup, rekap per bulan, dan "bulan ini" dihitung dalam WITA.
var SkmZona = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		return time.FixedZone("WITA", 8*3600)
	}
	return loc
}()

// SkmZonaSQL dipakai di query agregat per bulan (harus sama dengan SkmZona).
const SkmZonaSQL = "Asia/Makassar"

var SkmBulan = []string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}

type SkmMutu struct {
	Huruf string `json:"huruf"`
	Label string `json:"label"`
}

type SkmSkala struct {
	V     int    `json:"v"`
	Label string `json:"label"`
}

// SkmKategoriMutu: kategori mutu pelayanan untuk nilai indeks 0–100.
func SkmKategoriMutu(indeks *float64) *SkmMutu {
	if indeks == nil {
		return nil
	}
	for _, m := range skmMutu {
		if *indeks >= m.Min {
			return &SkmMutu{Huruf: m.Huruf, Label: m.Label}
		}
	}
	return nil
}

// SkmDaftarSkala: tingkat skala dari 1 sampai maks beserta labelnya.
func SkmDaftarSkala(maks int) []SkmSkala {
	labels := skmLabelSkala[maks]
	out := make([]SkmSkala, 0, len(labels))
	for i, l := range labels {
		out = append(out, SkmSkala{V: i + 1, Label: l})
	}
	return out
}

// SkmKategoriResponden: rata-rata responden dibulatkan ke tingkat skala terdekat.
func SkmKategoriResponden(rata float64, maks int) string {
	v := int(math.Round(rata))
	if v < 1 {
		v = 1
	}
	if v > maks {
		v = maks
	}
	return skmLabelSkala[maks][v-1]
}

func SkmIPHash(ip string) *string {
	if ip == "" {
		return nil
	}
	salt := os.Getenv("SKM_IP_SALT")
	if salt == "" {
		salt = os.Getenv("JWT_SECRET")
	}
	sum := sha256.Sum256([]byte(salt + "|" + ip))
	h := hex.EncodeToString(sum[:])
	return &h
}

func SkmBulat(n float64, d int) float64 {
	p := math.Pow(10, float64(d))
	return math.Round(n*p) / p
}

func SkmBulatPtr(n *float64, d int) *float64 {
	if n == nil {
		return nil
	}
	v := SkmBulat(*n, d)
	return &v
}

func SkmBatasPerJam(key string, def int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil && v > 0 {
		return v
	}
	return def
}

func SkmAda(daftar []string, v string) bool {
	for _, d := range daftar {
		if d == v {
			return true
		}
	}
	return false
}

var tagHTML = regexp.MustCompile(`<[^>]*>`)

// SkmBersih: teks bebas disimpan tanpa tag HTML. Tampilan tetap meng-escape keluaran.
func SkmBersih(s string) string {
	return strings.TrimSpace(tagHTML.ReplaceAllString(s, ""))
}

// SkmInisial: dua huruf awal nama untuk avatar, contoh "Ketut Pramana" -> "KP".
func SkmInisial(nama string) string {
	var out []rune
	for _, k := range strings.Fields(nama) {
		out = append(out, []rune(k)[0])
		if len(out) == 2 {
			break
		}
	}
	return strings.ToUpper(string(out))
}

// SkmPolaCari: pola ILIKE dengan wildcard pengguna di-escape.
func SkmPolaCari(kata string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(strings.TrimSpace(kata)) + "%"
}

func SkmISO(t time.Time) string { return t.In(SkmZona).Format(time.RFC3339) }
