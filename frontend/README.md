# SKM BRIDA Kabupaten Buleleng

Frontend **Survei Kepuasan Masyarakat** untuk Badan Riset dan Inovasi Daerah
Kabupaten Buleleng. React + Vite + Tailwind CSS, React Router. Seluruh data diambil dari
API Laravel di folder `../backend`.

## Menjalankan

Jalankan backend lebih dulu (lihat `../backend`), lalu:

```bash
cp .env.example .env   # Windows: copy .env.example .env
npm install
npm run dev      # http://localhost:5173
npm run build    # keluaran ke dist/
npm run preview
```

## Rute

| Rute | Halaman |
| --- | --- |
| `/` | Beranda (bagian Tentang, Hasil, Kontak: `/#tentang`, `/#hasil`, `/#kontak`) |
| `/survei` | Pengisian survei (4 langkah, draf di sessionStorage) |
| `/survei/selesai` | Konfirmasi pengiriman |
| `/admin/login` | Login admin |
| `/admin` | Dashboard |
| `/admin/responden` | Data responden |
| `/admin/survei` | Data survei |
| `/admin/pertanyaan` | Kelola pertanyaan |
| `/admin/hasil` | Hasil survei |
| `/admin/statistik` | Statistik |
| `/admin/feedback` | Feedback masyarakat |
| `/admin/laporan` | Laporan |
| `/admin/pengaturan` | Pengaturan |

Rute `/admin/*` dijaga `RequireAuth`: sesi diperiksa ke `GET /auth/me`, tamu dialihkan ke
`/admin/login`. Hak akses tetap ditegakkan backend; frontend hanya mengikuti hasilnya.

## Koneksi API

`VITE_API_URL` di `.env` (contoh: `http://localhost:8000/api/v1`). Semua permintaan lewat
`src/services/api/` (`client.js`, `publik.js`, `auth.js`, `admin.js`). Autentikasi memakai
Sanctum mode SPA: cookie sesi httpOnly + header `X-XSRF-TOKEN`, tidak ada token di
`localStorage`. Galat dinormalisasi menjadi `ApiError` (`status`, `message`, `errors`).
Hook `useApi()` mengelola status memuat/galat; `Memuat` dan `GagalMuat` menampilkannya.

## Struktur

```
src/
├── components/
│   ├── Icon.jsx            ikon garis inline (gaya Lucide)
│   ├── Logo.jsx            lambang BRIDA
│   ├── ui/                 Button, Card, Badge, StatCard, Modal, DataTable, …
│   └── charts/             BarUnsur, DonutKepuasan, TrenIKM, KolomResponden, MeterIKM
├── layouts/                Navbar, Footer, PublicLayout, AdminLayout, Shell
├── pages/
│   ├── publik/             Beranda, Tentang, Survei, Selesai, Hasil, Kontak
│   └── admin/              Login, Dashboard, Responden, DataSurvei, Pertanyaan,
│                           HasilSurvei, Statistik, Feedback, Laporan, Pengaturan
├── data/                   skala (warna kepuasan), feedback (status & warna badge)
├── hooks/                  useApi, useAuth, useSize (grafik responsif), useToast
├── services/api/           klien HTTP dan fungsi per endpoint
├── utils/                  format angka/tanggal, konversi mutu IKM
├── App.jsx                 definisi rute
└── index.css               token warna + dasar Tailwind
```

## Palet warna

Diambil dari identitas SINGA RISET BULELENG.

| Peran | Terang | Keterangan |
| --- | --- | --- |
| Merah utama | `#8E1B1B` | tombol, tautan aktif, batang grafik |
| Maroon | `#6B1414` | sidebar admin, pita CTA, panel login |
| Emas | `#F9C74F` | garis aksen, tombol emas, titik aktif |
| Emas teks | `#A87511` | label kecil di atas latar terang |
| Krem | `#FAF6F3` | latar halaman |

Seluruh warna adalah CSS custom property di `src/index.css` dan dipetakan ke kelas
Tailwind di `tailwind.config.js`. Mode gelap punya nilai tersendiri, bukan pembalikan
otomatis. Untuk mengubah tema, cukup sunting token di `:root` — tidak ada warna yang
ditulis langsung di komponen kecuali skala kepuasan.

## Skala penilaian

Skala diatur per periode di **Admin → Pengaturan** (1–5 dengan pengali 20, atau 1–4 sesuai
Permen PANRB 14/2017 dengan pengali 25). Backend menghitung indeks; frontend menerima
daftar tingkat skala dari API dan mencocokkan warnanya dengan `SKALA` lewat label
(`skalaDari()` di `src/data/skala.js`). Skala terkunci setelah periode punya responden.

Warna skala kepuasan divergen (merah → toska) sengaja dibedakan dari merah merek agar
"Sangat Tidak Puas" tidak tertukar dengan warna identitas. Pasangan warna bersebelahan
sudah diuji terpisah ΔE ≥ 18 pada penglihatan normal maupun defisiensi warna deutan.

## Data

Tidak ada lagi data contoh di frontend. Untuk pengembangan, seeder backend mengisi
database dengan data contoh (hanya pada `APP_ENV=local`).
