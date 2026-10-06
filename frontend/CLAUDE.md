# Catatan untuk Claude Code

Frontend SKM BRIDA Kabupaten Buleleng: React 18 + Vite + Tailwind 3 + React Router.
JavaScript, bukan TypeScript. Backend-nya API Laravel di `../backend` (Sanctum mode SPA).

## Aturan yang berlaku di repo ini

1. **Warna selalu lewat token.** Jangan menulis hex di komponen. Token ada di
   `src/index.css` (`:root`, blok `prefers-color-scheme: dark`, dan `[data-theme="dark"]`)
   dan dipetakan ke kelas Tailwind di `tailwind.config.js` (`bg-brand`, `text-muted`,
   `border-line`, dan seterusnya). Setiap token baru harus didefinisikan di ketiga blok.
   Satu-satunya pengecualian: skala kepuasan `--s1…--s5` beserta warna tintanya
   `--s1i…--s5i`, yang dipakai lewat `style={{ background, color }}`.

2. **Mode gelap bukan pembalikan otomatis.** Nilainya dipilih terpisah. Kalau menambah
   permukaan berwarna gelap dengan teks putih, pakai `var(--brand-deep)` (tetap gelap di
   kedua tema), bukan `var(--brand-strong)` (warna teks, jadi terang di mode gelap).

3. **Grafik ditulis manual dengan SVG**, tanpa pustaka chart. Pola yang dipakai:
   `useSize()` mengukur lebar container, lalu SVG digambar pada lebar piksel sebenarnya
   supaya teks tidak ikut melar. Keputusan tata letak grafik memakai lebar container
   (`w < 440`), bukan breakpoint viewport Tailwind.

4. **Item grid butuh `min-w-0`.** Kartu grafik di dalam grid akan melebarkan kolomnya
   kalau tidak diberi `min-w-0`; ini sumber overflow horizontal yang paling sering.

5. **Teks berbahasa Indonesia**, angka diformat `id-ID` lewat `src/utils/format.js`
   (`fmt`, `fmtInt`, `tglID`). Kolom angka memakai kelas `.font-num` agar rata digit.

6. **Notifikasi** lewat `useToast()` dari `src/hooks/useToast.jsx`, bukan `alert()`.
   `toast(pesan)` untuk sukses, `toast(pesan, 'galat')` untuk kegagalan.

7. **Tabel admin** memakai `DataTable` (`src/components/ui/DataTable.jsx`). Untuk data API
   pakai mode server: `fetchPage` (fungsi dari `services/api/admin.js`) + `params` filter;
   `key` kolom dikirim sebagai `sort`, kolom yang tak bisa diurutkan diberi `sortable:false`.

8. **Data selalu dari API.** Jangan menambah data dummy di `src/data/`. Panggilan HTTP
   hanya lewat `src/services/api/` (jangan `fetch` langsung atau hardcode URL; basis URL
   dari `VITE_API_URL`). Di komponen, pakai `useApi(fn, deps)` lalu tampilkan
   `Memuat` / `GagalMuat` / `EmptyState` untuk status memuat, galat, dan kosong.
   Galat 422 dipetakan ke field lewat `ex.field('nama_field')`.

9. **Autentikasi** lewat `useAuth()`; rute admin dibungkus `RequireAuth` di `App.jsx`.
   Jangan menyimpan token/sesi di storage. Menyembunyikan tombol hanyalah UX — aturan
   akses ditegakkan backend.

## Domain

- Unsur pelayanan (bank pertanyaan) dikelola admin dan diambil dari API; 9 unsur awal
  mengikuti Permen PANRB 14/2017.
- Indeks dihitung backend: rata-rata NRR unsur × (100 / skala_maks). Skala (1–4 atau
  1–5) diatur per periode; frontend memetakan warnanya dengan `skalaDari()` /
  `distribusiBerwarna()` di `src/data/skala.js`. Ambang mutu A/B/C/D (untuk warna badge)
  ada di `src/utils/mutu.js`.

## Perintah

```bash
npm run dev      # butuh backend berjalan di VITE_API_URL
npm run build
```

Belum ada test maupun linter frontend yang terpasang.
