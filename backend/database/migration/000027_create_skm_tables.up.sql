-- Survei Kepuasan Masyarakat (SKM) BRIDA Kabupaten Buleleng.
-- Responden mengisi tanpa akun; pengelolaan oleh users dengan role 'admin'.

CREATE TABLE IF NOT EXISTS skm_layanan (
    layanan_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    aktif BOOLEAN NOT NULL DEFAULT TRUE,
    urutan SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT skm_layanan_nama_unique UNIQUE (nama)
);

CREATE TABLE IF NOT EXISTS skm_periode (
    periode_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tahun SMALLINT NOT NULL,
    -- 4 = Permen PANRB 14/2017, 5 = skala yang dipakai survei ini.
    skala_maks SMALLINT NOT NULL DEFAULT 5,
    tanggal_tutup DATE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT skm_periode_tahun_unique UNIQUE (tahun),
    CONSTRAINT skm_periode_skala_check CHECK (skala_maks IN (4, 5))
);

-- Hanya berisi satu baris.
CREATE TABLE IF NOT EXISTS skm_pengaturan (
    pengaturan_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    judul VARCHAR(150) NOT NULL,
    periode_aktif_id BIGINT NOT NULL REFERENCES skm_periode (periode_id) ON DELETE RESTRICT,
    anonim BOOLEAN NOT NULL DEFAULT TRUE,
    publik BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Soft delete: jawaban periode lama tetap merujuk ke pertanyaan yang dihapus.
CREATE TABLE IF NOT EXISTS skm_pertanyaan (
    pertanyaan_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    kode VARCHAR(10) NOT NULL,
    unsur VARCHAR(100) NOT NULL,
    teks TEXT NOT NULL,
    wajib BOOLEAN NOT NULL DEFAULT TRUE,
    aktif BOOLEAN NOT NULL DEFAULT TRUE,
    urutan SMALLINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT skm_pertanyaan_kode_unique UNIQUE (kode)
);

CREATE INDEX IF NOT EXISTS idx_skm_pertanyaan_aktif_urutan ON skm_pertanyaan (aktif, urutan);
CREATE INDEX IF NOT EXISTS idx_skm_pertanyaan_deleted_at ON skm_pertanyaan (deleted_at);

-- Satu pengisian survei (responden tanpa akun). Nomor tampilan SKM-<tahun>-<urut>.
CREATE TABLE IF NOT EXISTS skm_responden (
    responden_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    periode_id BIGINT NOT NULL REFERENCES skm_periode (periode_id) ON DELETE RESTRICT,
    urut INT NOT NULL,
    nomor VARCHAR(20) NOT NULL,
    layanan_id BIGINT NOT NULL REFERENCES skm_layanan (layanan_id) ON DELETE RESTRICT,
    nama VARCHAR(100),
    jk VARCHAR(20) NOT NULL,
    usia SMALLINT NOT NULL,
    pendidikan VARCHAR(40) NOT NULL,
    pekerjaan VARCHAR(40) NOT NULL,
    frekuensi VARCHAR(40) NOT NULL,
    kecamatan VARCHAR(40),
    apresiasi VARCHAR(300),
    rata_rata NUMERIC(4, 2) NOT NULL,
    -- IP mentah tidak disimpan, hanya sha256(salt|ip).
    ip_hash CHAR(64),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT skm_responden_nomor_unique UNIQUE (nomor),
    CONSTRAINT skm_responden_periode_urut_unique UNIQUE (periode_id, urut),
    CONSTRAINT skm_responden_jk_check CHECK (jk IN ('Laki-laki', 'Perempuan'))
);

CREATE INDEX IF NOT EXISTS idx_skm_responden_periode_created ON skm_responden (periode_id, created_at);
CREATE INDEX IF NOT EXISTS idx_skm_responden_periode_layanan ON skm_responden (periode_id, layanan_id);
CREATE INDEX IF NOT EXISTS idx_skm_responden_periode_pendidikan ON skm_responden (periode_id, pendidikan);
CREATE INDEX IF NOT EXISTS idx_skm_responden_periode_pekerjaan ON skm_responden (periode_id, pekerjaan);
CREATE INDEX IF NOT EXISTS idx_skm_responden_ip_created ON skm_responden (ip_hash, created_at);

CREATE TABLE IF NOT EXISTS skm_jawaban (
    jawaban_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    responden_id BIGINT NOT NULL REFERENCES skm_responden (responden_id) ON DELETE CASCADE,
    pertanyaan_id BIGINT NOT NULL REFERENCES skm_pertanyaan (pertanyaan_id) ON DELETE RESTRICT,
    nilai SMALLINT NOT NULL,
    CONSTRAINT skm_jawaban_responden_pertanyaan_unique UNIQUE (responden_id, pertanyaan_id)
);

CREATE INDEX IF NOT EXISTS idx_skm_jawaban_pertanyaan_nilai ON skm_jawaban (pertanyaan_id, nilai);

-- Saran responden beserta status tindak lanjutnya oleh admin.
CREATE TABLE IF NOT EXISTS skm_feedback (
    feedback_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    responden_id BIGINT NOT NULL REFERENCES skm_responden (responden_id) ON DELETE CASCADE,
    kategori_id BIGINT REFERENCES skm_pertanyaan (pertanyaan_id) ON DELETE SET NULL,
    isi VARCHAR(600) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'Baru',
    catatan TEXT,
    ditangani_oleh BIGINT REFERENCES users (internal_id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT skm_feedback_responden_unique UNIQUE (responden_id),
    CONSTRAINT skm_feedback_status_check CHECK (status IN ('Baru', 'Ditinjau', 'Ditindaklanjuti', 'Selesai'))
);

CREATE INDEX IF NOT EXISTS idx_skm_feedback_status ON skm_feedback (status);
CREATE INDEX IF NOT EXISTS idx_skm_feedback_created ON skm_feedback (created_at);

CREATE TABLE IF NOT EXISTS skm_pesan_kontak (
    pesan_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    topik VARCHAR(40) NOT NULL,
    pesan TEXT NOT NULL,
    dibaca BOOLEAN NOT NULL DEFAULT FALSE,
    ip_hash CHAR(64),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_skm_pesan_kontak_dibaca_created ON skm_pesan_kontak (dibaca, created_at);

CREATE TRIGGER update_skm_layanan_updated_at BEFORE UPDATE ON skm_layanan FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_skm_periode_updated_at BEFORE UPDATE ON skm_periode FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_skm_pengaturan_updated_at BEFORE UPDATE ON skm_pengaturan FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_skm_pertanyaan_updated_at BEFORE UPDATE ON skm_pertanyaan FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_skm_responden_updated_at BEFORE UPDATE ON skm_responden FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_skm_feedback_updated_at BEFORE UPDATE ON skm_feedback FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_skm_pesan_kontak_updated_at BEFORE UPDATE ON skm_pesan_kontak FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
