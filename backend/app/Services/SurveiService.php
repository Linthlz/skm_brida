<?php

namespace App\Services;

use App\Models\Pengaturan;
use App\Models\Periode;
use App\Models\Pertanyaan;
use App\Models\Responden;
use App\Support\Skm;
use Illuminate\Support\Arr;
use Illuminate\Support\Facades\DB;

class SurveiService
{
    /** Data yang dibutuhkan halaman /survei: pengaturan, skala, dan pertanyaan aktif. */
    public function form(): array
    {
        $pengaturan = Pengaturan::sekarang();
        $periode = $pengaturan->periodeAktif;

        return [
            'judul' => $pengaturan->judul,
            'anonim' => $pengaturan->anonim,
            'periode' => [
                'tahun' => $periode->tahun,
                'skala_maks' => $periode->skala_maks,
                'tanggal_tutup' => $periode->tanggal_tutup?->toDateString(),
                'ditutup' => $periode->ditutup(),
            ],
            'skala' => Skm::skala($periode->skala_maks),
            'pertanyaan' => Pertanyaan::aktif()->berurutan()->get(['id', 'kode', 'unsur', 'teks', 'wajib']),
        ];
    }

    /**
     * Menyimpan satu pengisian survei beserta jawaban dan sarannya dalam satu
     * transaksi. Data sudah tervalidasi oleh KirimSurveiRequest.
     */
    public function kirim(array $data, Periode $periode, ?string $ip): Responden
    {
        return DB::transaction(function () use ($data, $periode, $ip) {
            // Kunci baris periode agar nomor urut tidak bentrok saat pengiriman bersamaan.
            Periode::whereKey($periode->id)->lockForUpdate()->first();
            $urut = (int) Responden::where('periode_id', $periode->id)->max('urut') + 1;

            $nilai = collect($data['jawaban'])->mapWithKeys(fn ($j) => [(int) $j['pertanyaan_id'] => (int) $j['nilai']]);

            $responden = new Responden(Arr::only($data, [
                'layanan_id', 'nama', 'jk', 'usia', 'pendidikan', 'pekerjaan', 'frekuensi', 'kecamatan', 'apresiasi',
            ]));
            $responden->periode_id = $periode->id;
            $responden->urut = $urut;
            $responden->nomor = sprintf('SKM-%d-%04d', $periode->tahun, $urut);
            $responden->rata_rata = round($nilai->avg(), 2);
            $responden->ip_hash = Skm::ipHash($ip);
            $responden->save();

            $responden->jawaban()->createMany(
                $nilai->map(fn ($v, $id) => ['pertanyaan_id' => $id, 'nilai' => $v])->values()->all()
            );

            if (filled($data['saran'] ?? null)) {
                $responden->feedback()->create([
                    'isi' => $data['saran'],
                    'kategori_id' => $this->unsurTerendah($nilai->all()),
                    'status' => 'Baru',
                ]);
            }

            return $responden->setRelation('periode', $periode);
        });
    }

    /**
     * Kategori awal feedback = unsur dengan nilai terendah dari responden itu.
     * Bila ada yang sama rendah, dipilih yang urutannya paling awal.
     *
     * @param  array<int,int>  $nilai  pertanyaan_id => nilai
     */
    private function unsurTerendah(array $nilai): ?int
    {
        $min = min($nilai);
        $kandidat = array_keys(array_filter($nilai, fn ($v) => $v === $min));

        return Pertanyaan::whereIn('id', $kandidat)->berurutan()->value('id');
    }
}
