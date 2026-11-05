-- Data induk SKM yang dibutuhkan di semua lingkungan, termasuk produksi.

INSERT INTO skm_layanan (nama, urutan) VALUES
    ('Konsultasi Riset & Inovasi', 1),
    ('Rekomendasi Izin Penelitian', 2),
    ('Fasilitasi Kajian Kebijakan', 3),
    ('Pendampingan Inovasi Daerah', 4),
    ('Layanan Data & Informasi Riset', 5),
    ('Publikasi & Diseminasi Hasil Riset', 6)
ON CONFLICT (nama) DO NOTHING;

-- Lima periode terakhir; periode tahun berjalan ditutup 31 Desember.
INSERT INTO skm_periode (tahun, skala_maks, tanggal_tutup)
SELECT t, 5, CASE WHEN t = EXTRACT(YEAR FROM CURRENT_DATE)::INT THEN make_date(t, 12, 31) END
FROM generate_series(EXTRACT(YEAR FROM CURRENT_DATE)::INT - 4, EXTRACT(YEAR FROM CURRENT_DATE)::INT) AS t
ON CONFLICT (tahun) DO NOTHING;

-- 9 unsur pelayanan Permen PANRB 14/2017.
INSERT INTO skm_pertanyaan (kode, unsur, teks, wajib, aktif, urutan) VALUES
    ('U1', 'Persyaratan Pelayanan', 'Bagaimana pendapat Anda tentang kesesuaian persyaratan pelayanan dengan jenis layanan yang Anda terima?', TRUE, TRUE, 1),
    ('U2', 'Sistem, Mekanisme & Prosedur', 'Bagaimana pemahaman Anda tentang kemudahan prosedur pelayanan di BRIDA Kabupaten Buleleng?', TRUE, TRUE, 2),
    ('U3', 'Waktu Penyelesaian', 'Bagaimana pendapat Anda tentang kecepatan waktu penyelesaian layanan?', TRUE, TRUE, 3),
    ('U4', 'Biaya / Tarif', 'Bagaimana pendapat Anda tentang kewajaran biaya/tarif dalam pelayanan?', TRUE, TRUE, 4),
    ('U5', 'Produk Spesifikasi Layanan', 'Bagaimana pendapat Anda tentang kesesuaian hasil layanan dengan ketentuan yang telah ditetapkan?', TRUE, TRUE, 5),
    ('U6', 'Kompetensi Pelaksana', 'Bagaimana pendapat Anda tentang kompetensi/kemampuan petugas dalam memberikan pelayanan?', TRUE, TRUE, 6),
    ('U7', 'Perilaku Pelaksana', 'Bagaimana pendapat Anda tentang keramahan dan kesopanan petugas dalam memberikan pelayanan?', TRUE, TRUE, 7),
    ('U8', 'Penanganan Pengaduan', 'Bagaimana pendapat Anda tentang penanganan pengaduan, saran, dan masukan?', TRUE, TRUE, 8),
    ('U9', 'Sarana & Prasarana', 'Bagaimana pendapat Anda tentang kualitas sarana dan prasarana layanan?', TRUE, TRUE, 9)
ON CONFLICT (kode) DO NOTHING;

INSERT INTO skm_pengaturan (judul, periode_aktif_id, anonim, publik)
SELECT 'Survei Kepuasan Masyarakat BRIDA Kabupaten Buleleng', periode_id, TRUE, TRUE
FROM skm_periode
WHERE tahun = EXTRACT(YEAR FROM CURRENT_DATE)::INT
  AND NOT EXISTS (SELECT 1 FROM skm_pengaturan);
