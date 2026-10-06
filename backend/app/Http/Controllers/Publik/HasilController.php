<?php

namespace App\Http\Controllers\Publik;

use App\Http\Controllers\Controller;
use App\Http\Requests\FilterHasilRequest;
use App\Models\Pengaturan;
use App\Models\Periode;
use App\Services\HasilService;
use Illuminate\Http\JsonResponse;

/** Hasil agregat untuk situs publik. Ditolak bila publikasi dimatikan di Pengaturan. */
class HasilController extends Controller
{
    public function __construct(private HasilService $hasil) {}

    public function periode(): JsonResponse
    {
        $aktif = Pengaturan::sekarang()->periode_aktif_id;

        return $this->sukses(
            Periode::orderByDesc('tahun')->get(['id', 'tahun'])
                ->map(fn ($p) => ['tahun' => $p->tahun, 'aktif' => $p->id === $aktif])
        );
    }

    public function index(FilterHasilRequest $request): JsonResponse
    {
        $this->pastikanPublik();

        return $this->sukses($this->hasil->rekap($request->periode(), $request->filter()));
    }

    public function ringkasan(): JsonResponse
    {
        $this->pastikanPublik();

        return $this->sukses($this->hasil->ringkasan(Pengaturan::sekarang()->periodeAktif));
    }

    private function pastikanPublik(): void
    {
        abort_unless(Pengaturan::sekarang()->publik, 403, 'Hasil survei belum dipublikasikan.');
    }
}
