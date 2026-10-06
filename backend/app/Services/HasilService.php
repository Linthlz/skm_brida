<?php

namespace App\Services;

use App\Models\Feedback;
use App\Models\Jawaban;
use App\Models\Periode;
use App\Models\PesanKontak;
use App\Models\Pertanyaan;
use App\Models\Responden;
use App\Support\Skm;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;

/**
 * Seluruh perhitungan agregat SKM. Indeks = rata-rata NRR unsur × (100 / skala_maks),
 * sesuai Permen PANRB 14/2017 (bobot tiap unsur sama).
 *
 * Filter yang didukung: layanan_id, pendidikan.
 */
class HasilService
{
    private const BULAN = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'];

    /** Rekap utama untuk halaman Hasil (publik & admin) dan Laporan. */
    public function rekap(Periode $periode, array $filter = []): array
    {
        $unsur = $this->unsur($periode, $filter);
        $nrr = $unsur->isEmpty() ? null : $unsur->avg('nrr');
        $ikm = $nrr === null ? null : $nrr * $periode->pengali();
        $total = $this->responden($periode, $filter)->count();
        $rentang = $this->rentangBulan($periode, $filter);

        return [
            'periode' => $periode->tahun,
            'skala_maks' => $periode->skala_maks,
            'pengali' => $periode->pengali(),
            'total_responden' => $total,
            'rentang_bulan' => $rentang,
            'nrr' => Skm::bulat($nrr),
            'ikm' => Skm::bulat($ikm),
            'mutu' => Skm::mutu($ikm),
            'unsur' => $unsur->map(fn ($u) => [
                ...$u,
                'nrr' => Skm::bulat($u['nrr']),
                'nrr_tertimbang' => Skm::bulat($u['nrr'] / $unsur->count(), 3),
                'indeks' => Skm::bulat($u['nrr'] * $periode->pengali()),
                'mutu' => Skm::mutu($u['nrr'] * $periode->pengali()),
            ])->values()->all(),
            'distribusi' => $this->distribusi($periode, $filter),
            'tren' => $this->tren($filter),
        ];
    }

    /** Angka ringkas untuk Hero, Tentang, dan panel login. */
    public function ringkasan(Periode $periode): array
    {
        $unsur = $this->unsur($periode);
        $ikm = $unsur->isEmpty() ? null : $unsur->avg('nrr') * $periode->pengali();

        return [
            'periode' => $periode->tahun,
            'ikm' => Skm::bulat($ikm),
            'mutu' => Skm::mutu($ikm),
            'total_responden' => Responden::where('periode_id', $periode->id)->count(),
            'total_feedback' => Feedback::whereHas('responden', fn ($q) => $q->where('periode_id', $periode->id))->count(),
            'jumlah_unsur' => Pertanyaan::aktif()->count(),
            'jumlah_layanan' => DB::table('layanan')->where('aktif', true)->count(),
        ];
    }

    public function dashboard(Periode $periode): array
    {
        $rekap = $this->rekap($periode);
        $total = $rekap['total_responden'];
        $lengkap = $this->jumlahLengkap($periode);

        return [
            'periode' => $periode->tahun,
            'skala_maks' => $periode->skala_maks,
            'pengali' => $periode->pengali(),
            'total_responden' => $total,
            'bulan_ini' => Responden::where('periode_id', $periode->id)
                ->where('created_at', '>=', now()->startOfMonth())->count(),
            'persen_lengkap' => $total ? round($lengkap / $total * 100) : null,
            'ikm' => $rekap['ikm'],
            'mutu' => $rekap['mutu'],
            'feedback_total' => Feedback::whereHas('responden', fn ($q) => $q->where('periode_id', $periode->id))->count(),
            'feedback_baru' => Feedback::where('status', 'Baru')->count(),
            'pesan_belum_dibaca' => PesanKontak::where('dibaca', false)->count(),
            'per_bulan' => $this->perBulan($periode),
            'distribusi' => $rekap['distribusi'],
            'unsur' => $rekap['unsur'],
            'prioritas' => collect($rekap['unsur'])->sortBy('nrr')->take(3)->values()->all(),
        ];
    }

    public function statistik(Periode $periode): array
    {
        return [
            'periode' => $periode->tahun,
            'total_responden' => Responden::where('periode_id', $periode->id)->count(),
            'pendidikan' => $this->hitungPer($periode, 'pendidikan', config('skm.pendidikan')),
            'pekerjaan' => $this->hitungPer($periode, 'pekerjaan', config('skm.pekerjaan')),
            'tren' => $this->tren(),
        ];
    }

    /** Query dasar responden pada satu periode dengan filter. */
    private function responden(Periode $periode, array $filter = []): Builder
    {
        return Responden::query()
            ->where('responden.periode_id', $periode->id)
            ->when($filter['layanan_id'] ?? null, fn ($q, $v) => $q->where('responden.layanan_id', $v))
            ->when($filter['pendidikan'] ?? null, fn ($q, $v) => $q->where('responden.pendidikan', $v));
    }

    /** NRR tiap unsur yang punya jawaban pada periode tersebut, urut sesuai bank pertanyaan. */
    private function unsur(Periode $periode, array $filter = []): Collection
    {
        $agregat = Jawaban::query()
            ->join('responden', 'responden.id', '=', 'jawaban.responden_id')
            ->where('responden.periode_id', $periode->id)
            ->when($filter['layanan_id'] ?? null, fn ($q, $v) => $q->where('responden.layanan_id', $v))
            ->when($filter['pendidikan'] ?? null, fn ($q, $v) => $q->where('responden.pendidikan', $v))
            ->groupBy('jawaban.pertanyaan_id')
            ->selectRaw('jawaban.pertanyaan_id, AVG(jawaban.nilai) as nrr, COUNT(*) as n')
            ->get()
            ->keyBy('pertanyaan_id');

        if ($agregat->isEmpty()) {
            return collect();
        }

        return Pertanyaan::withTrashed()
            ->whereIn('id', $agregat->keys())
            ->berurutan()
            ->get(['id', 'kode', 'unsur'])
            ->map(fn ($p) => [
                'id' => $p->id,
                'kode' => $p->kode,
                'nama' => $p->unsur,
                'nrr' => (float) $agregat[$p->id]->nrr,
                'jumlah_jawaban' => (int) $agregat[$p->id]->n,
            ]);
    }

    /** Sebaran responden per tingkat kepuasan, dari tingkat tertinggi ke terendah. */
    private function distribusi(Periode $periode, array $filter = []): array
    {
        $hitung = $this->responden($periode, $filter)
            ->selectRaw('ROUND(rata_rata) as v, COUNT(*) as jumlah')
            ->groupByRaw('ROUND(rata_rata)')
            ->pluck('jumlah', 'v')
            ->mapWithKeys(fn ($n, $v) => [(int) $v => (int) $n]);

        return collect(Skm::skala($periode->skala_maks))
            ->reverse()
            ->map(fn ($s) => [...$s, 'jumlah' => $hitung[$s['v']] ?? 0])
            ->values()
            ->all();
    }

    /** IKM tiap periode yang punya data, urut tahun naik. */
    private function tren(array $filter = []): array
    {
        $rows = Jawaban::query()
            ->join('responden', 'responden.id', '=', 'jawaban.responden_id')
            ->when($filter['layanan_id'] ?? null, fn ($q, $v) => $q->where('responden.layanan_id', $v))
            ->when($filter['pendidikan'] ?? null, fn ($q, $v) => $q->where('responden.pendidikan', $v))
            ->groupBy('responden.periode_id', 'jawaban.pertanyaan_id')
            ->selectRaw('responden.periode_id, AVG(jawaban.nilai) as nrr')
            ->get()
            ->groupBy('periode_id');

        return Periode::whereIn('id', $rows->keys())->orderBy('tahun')->get()
            ->map(fn ($p) => [
                'periode' => (string) $p->tahun,
                'ikm' => Skm::bulat($rows[$p->id]->avg('nrr') * $p->pengali()),
            ])
            ->all();
    }

    /** Jumlah responden per bulan, dari Januari sampai bulan terakhir yang berisi data. */
    private function perBulan(Periode $periode): array
    {
        $bulanSql = DB::connection()->getDriverName() === 'sqlite'
            ? "CAST(strftime('%m', created_at) AS INTEGER)"
            : 'MONTH(created_at)';

        $hitung = Responden::where('periode_id', $periode->id)
            ->selectRaw("$bulanSql as bulan, COUNT(*) as jumlah")
            ->groupByRaw($bulanSql)
            ->pluck('jumlah', 'bulan')
            ->mapWithKeys(fn ($n, $b) => [(int) $b => (int) $n]);

        $akhir = $hitung->keys()->max() ?? 0;

        return collect(range(1, max(1, $akhir)))
            ->map(fn ($b) => ['bulan' => $b, 'label' => self::BULAN[$b - 1], 'jumlah' => $hitung[$b] ?? 0])
            ->all();
    }

    /** Bulan pertama & terakhir pengisian, contoh "Jan – Sep". */
    private function rentangBulan(Periode $periode, array $filter): ?string
    {
        $r = $this->responden($periode, $filter)
            ->selectRaw('MIN(created_at) as awal, MAX(created_at) as akhir')
            ->first();

        if (! $r?->awal) {
            return null;
        }
        $awal = self::BULAN[(int) date('n', strtotime($r->awal)) - 1];
        $akhir = self::BULAN[(int) date('n', strtotime($r->akhir)) - 1];

        return $awal === $akhir ? $awal : "$awal – $akhir";
    }

    /** Responden yang menjawab seluruh pertanyaan yang pernah muncul pada periode itu. */
    private function jumlahLengkap(Periode $periode): int
    {
        $jumlahPertanyaan = Jawaban::query()
            ->join('responden', 'responden.id', '=', 'jawaban.responden_id')
            ->where('responden.periode_id', $periode->id)
            ->distinct()
            ->count('jawaban.pertanyaan_id');

        if ($jumlahPertanyaan === 0) {
            return 0;
        }

        return Responden::where('periode_id', $periode->id)
            ->has('jawaban', '>=', $jumlahPertanyaan)
            ->count();
    }

    private function hitungPer(Periode $periode, string $kolom, array $opsi): array
    {
        $hitung = Responden::where('periode_id', $periode->id)
            ->selectRaw("$kolom as label, COUNT(*) as n")
            ->groupBy($kolom)
            ->pluck('n', 'label');

        return collect($opsi)->map(fn ($o) => ['label' => $o, 'n' => (int) ($hitung[$o] ?? 0)])->all();
    }
}
