<?php

namespace Database\Seeders;

use App\Models\Layanan;
use App\Models\Periode;
use App\Models\PesanKontak;
use App\Models\Pertanyaan;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Str;

/**
 * Data contoh untuk pengembangan (BUKAN data masyarakat nyata). Hanya
 * dijalankan DatabaseSeeder pada lingkungan local. Hasilnya stabil antar
 * eksekusi karena generator acak diberi seed tetap.
 */
class DataContohSeeder extends Seeder
{
    /** Rata-rata dasar per tahun (indeks naik perlahan) dan selisih per unsur. */
    private const DASAR = [-4 => 4.06, -3 => 4.17, -2 => 4.25, -1 => 4.31, 0 => 4.37];

    private const SELISIH_UNSUR = [0.05, -0.06, -0.22, 0.34, -0.03, 0.09, 0.21, -0.27, -0.09];

    private const SARAN = [
        'Proses rekomendasi izin penelitian sudah cepat, tapi pemberitahuan kalau berkas sudah selesai sebaiknya dikirim lewat WhatsApp juga.',
        'Ruang tunggu perlu tambahan kursi pada hari Senin karena antrean cukup panjang.',
        'Saya mengajukan pertanyaan lewat email dan baru dibalas tiga hari kemudian. Mohon ada kanal pengaduan yang lebih responsif.',
        'Petugas front office sangat ramah dan menjelaskan alur pendampingan inovasi dengan sabar. Terima kasih.',
        'Alur pengajuan kajian kebijakan masih membingungkan. Mohon dibuat infografis alurnya di website.',
        'Data hasil riset yang dipublikasikan sebaiknya bisa diunduh dalam format terbuka (CSV), bukan hanya PDF.',
        'Konsultasi daring sangat membantu karena saya dari Kecamatan Gerokgak, tidak perlu ke Singaraja.',
        'Jaringan wifi di ruang layanan sering terputus saat mengunggah berkas.',
    ];

    public function run(): void
    {
        if (DB::table('responden')->exists()) {
            $this->command?->info('Data responden sudah ada — data contoh dilewati.');

            return;
        }

        mt_srand(2026);
        $pertanyaan = Pertanyaan::berurutan()->pluck('id')->all();
        $layanan = Layanan::pluck('id')->all();
        $tahunIni = (int) now()->year;

        DB::transaction(function () use ($pertanyaan, $layanan, $tahunIni) {
            foreach (Periode::orderBy('tahun')->get() as $periode) {
                $selisihTahun = $periode->tahun - $tahunIni;
                $dasar = self::DASAR[$selisihTahun] ?? 4.0;
                $akhir = $periode->tahun === $tahunIni ? now() : Carbon::create($periode->tahun, 12, 31, 15);
                $awal = Carbon::create($periode->tahun, 1, 6, 8);
                $jumlah = $periode->tahun === $tahunIni ? 140 * $akhir->month : 1500 + mt_rand(-150, 150);

                for ($i = 1; $i <= $jumlah; $i++) {
                    $this->satuResponden($periode, $i, $dasar, $pertanyaan, $layanan, $awal, $akhir);
                }
            }

            $this->pesanKontak();
        });

        $this->pengelolaContoh();
    }

    private function satuResponden(Periode $periode, int $urut, float $dasar, array $pertanyaan, array $layanan, Carbon $awal, Carbon $akhir): void
    {
        $waktu = Carbon::createFromTimestamp(mt_rand($awal->timestamp, $akhir->timestamp));
        // Kecenderungan pribadi responden agar sebaran kategori kepuasan tidak menumpuk di tengah.
        $pribadi = $this->normal(0, 0.55);
        $nilai = [];
        foreach ($pertanyaan as $k => $pid) {
            $nilai[$pid] = max(1, min(5, (int) round($this->normal($dasar + $pribadi + (self::SELISIH_UNSUR[$k] ?? 0), 0.6))));
        }
        $rata = round(array_sum($nilai) / count($nilai), 2);

        $id = DB::table('responden')->insertGetId([
            'periode_id' => $periode->id,
            'urut' => $urut,
            'nomor' => sprintf('SKM-%d-%04d', $periode->tahun, $urut),
            'layanan_id' => $layanan[array_rand($layanan)],
            'nama' => null,
            'jk' => mt_rand(0, 2) === 0 ? 'Perempuan' : 'Laki-laki',
            'usia' => mt_rand(19, 60),
            'pendidikan' => $this->acak(config('skm.pendidikan'), [1, 2, 10, 8, 22, 9]),
            'pekerjaan' => $this->acak(config('skm.pekerjaan'), [16, 14, 9, 8, 11, 3, 4]),
            'frekuensi' => $this->acak(config('skm.frekuensi'), [9, 6, 3, 2]),
            'kecamatan' => $this->acak(config('skm.kecamatan'), [30, 14, 8, 8, 9, 6, 5, 4, 6]),
            'apresiasi' => null,
            'rata_rata' => $rata,
            'ip_hash' => null,
            'created_at' => $waktu,
            'updated_at' => $waktu,
        ]);

        DB::table('jawaban')->insert(collect($nilai)->map(fn ($v, $pid) => [
            'responden_id' => $id, 'pertanyaan_id' => $pid, 'nilai' => $v,
        ])->values()->all());

        if (mt_rand(1, 100) <= 22) {
            $min = min($nilai);
            $umur = $waktu->diffInDays(now());
            $status = match (true) {
                $umur > 120 => 'Selesai',
                $umur > 45 => $this->acak(['Ditindaklanjuti', 'Selesai'], [1, 3]),
                $umur > 14 => $this->acak(['Ditinjau', 'Ditindaklanjuti'], [2, 1]),
                default => $this->acak(['Baru', 'Ditinjau'], [3, 1]),
            };
            DB::table('feedback')->insert([
                'responden_id' => $id,
                'kategori_id' => array_search($min, $nilai, true),
                'isi' => self::SARAN[array_rand(self::SARAN)],
                'status' => $status,
                'created_at' => $waktu,
                'updated_at' => $waktu,
            ]);
        }
    }

    private function pesanKontak(): void
    {
        $contoh = [
            ['Pertanyaan tentang survei', 'Apakah hasil survei semester lalu bisa diunduh dalam bentuk laporan lengkap?'],
            ['Permintaan data hasil SKM', 'Kami membutuhkan data SKM tiga tahun terakhir untuk bahan penelitian skripsi.'],
            ['Kerja sama riset & inovasi', 'Kampus kami ingin menjajaki kerja sama riset terkait inovasi pelayanan publik.'],
        ];
        foreach ($contoh as $i => [$topik, $pesan]) {
            $p = new PesanKontak(['nama' => 'Pengirim Contoh '.($i + 1), 'email' => 'contoh'.($i + 1).'@example.com', 'topik' => $topik, 'pesan' => $pesan]);
            $p->dibaca = $i === 0;
            $p->save();
        }
    }

    /** Pengelola contoh agar daftar di Pengaturan tidak kosong. Kata sandinya acak (tidak dapat dipakai masuk). */
    private function pengelolaContoh(): void
    {
        foreach ([['Ni Wayan Ardani', 'Analis Data'], ['Gede Suarnata', 'Kepala Bidang']] as $i => [$nama, $jabatan]) {
            User::firstOrCreate(['email' => 'pengelola'.($i + 1).'@example.com'], [
                'nama' => $nama, 'jabatan' => $jabatan, 'password' => Str::random(40), 'aktif' => true,
            ]);
        }
    }

    private function normal(float $mean, float $sd): float
    {
        $u = max(1e-9, mt_rand() / mt_getrandmax());
        $v = mt_rand() / mt_getrandmax();

        return $mean + $sd * sqrt(-2 * log($u)) * cos(2 * M_PI * $v);
    }

    private function acak(array $opsi, array $bobot): string
    {
        $r = mt_rand(1, array_sum($bobot));
        foreach ($opsi as $i => $o) {
            $r -= $bobot[$i] ?? 1;
            if ($r <= 0) {
                return $o;
            }
        }

        return end($opsi);
    }
}
