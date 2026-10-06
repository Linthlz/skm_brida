<?php

namespace Database\Seeders;

use App\Models\Layanan;
use App\Models\Pengaturan;
use App\Models\Periode;
use App\Models\Pertanyaan;
use Illuminate\Database\Seeder;

/** Data induk yang dibutuhkan di semua lingkungan, termasuk production. Aman dijalankan ulang. */
class MasterSeeder extends Seeder
{
    public const LAYANAN = [
        'Konsultasi Riset & Inovasi', 'Rekomendasi Izin Penelitian', 'Fasilitasi Kajian Kebijakan',
        'Pendampingan Inovasi Daerah', 'Layanan Data & Informasi Riset', 'Publikasi & Diseminasi Hasil Riset',
    ];

    /** 9 unsur pelayanan Permen PANRB 14/2017. */
    public const UNSUR = [
        ['Persyaratan Pelayanan', 'Bagaimana pendapat Anda tentang kesesuaian persyaratan pelayanan dengan jenis layanan yang Anda terima?'],
        ['Sistem, Mekanisme & Prosedur', 'Bagaimana pemahaman Anda tentang kemudahan prosedur pelayanan di BRIDA Kabupaten Buleleng?'],
        ['Waktu Penyelesaian', 'Bagaimana pendapat Anda tentang kecepatan waktu penyelesaian layanan?'],
        ['Biaya / Tarif', 'Bagaimana pendapat Anda tentang kewajaran biaya/tarif dalam pelayanan?'],
        ['Produk Spesifikasi Layanan', 'Bagaimana pendapat Anda tentang kesesuaian hasil layanan dengan ketentuan yang telah ditetapkan?'],
        ['Kompetensi Pelaksana', 'Bagaimana pendapat Anda tentang kompetensi/kemampuan petugas dalam memberikan pelayanan?'],
        ['Perilaku Pelaksana', 'Bagaimana pendapat Anda tentang keramahan dan kesopanan petugas dalam memberikan pelayanan?'],
        ['Penanganan Pengaduan', 'Bagaimana pendapat Anda tentang penanganan pengaduan, saran, dan masukan?'],
        ['Sarana & Prasarana', 'Bagaimana pendapat Anda tentang kualitas sarana dan prasarana layanan?'],
    ];

    public function run(): void
    {
        foreach (self::LAYANAN as $i => $nama) {
            Layanan::firstOrCreate(['nama' => $nama], ['urutan' => $i + 1]);
        }

        $tahunIni = (int) now()->year;
        foreach (range($tahunIni - 4, $tahunIni) as $tahun) {
            Periode::firstOrCreate(['tahun' => $tahun], [
                'skala_maks' => 5,
                'tanggal_tutup' => $tahun === $tahunIni ? "$tahun-12-31" : null,
            ]);
        }

        if (! Pertanyaan::withTrashed()->exists()) {
            foreach (self::UNSUR as $i => [$unsur, $teks]) {
                Pertanyaan::create([
                    'kode' => 'U'.($i + 1), 'unsur' => $unsur, 'teks' => $teks,
                    'wajib' => true, 'aktif' => true, 'urutan' => $i + 1,
                ]);
            }
        }

        if (! Pengaturan::exists()) {
            Pengaturan::create([
                'judul' => 'Survei Kepuasan Masyarakat BRIDA Kabupaten Buleleng',
                'periode_aktif_id' => Periode::where('tahun', $tahunIni)->value('id'),
                'anonim' => true,
                'publik' => true,
            ]);
        }
    }
}
