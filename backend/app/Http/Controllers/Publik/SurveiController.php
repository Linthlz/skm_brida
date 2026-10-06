<?php

namespace App\Http\Controllers\Publik;

use App\Http\Controllers\Controller;
use App\Http\Requests\Publik\KirimSurveiRequest;
use App\Services\SurveiService;
use App\Support\Skm;
use Illuminate\Http\JsonResponse;

class SurveiController extends Controller
{
    public function __construct(private SurveiService $survei) {}

    public function form(): JsonResponse
    {
        return $this->sukses($this->survei->form());
    }

    public function store(KirimSurveiRequest $request): JsonResponse
    {
        $periode = $request->pengaturan()->periodeAktif;
        $r = $this->survei->kirim($request->validated(), $periode, $request->ip());
        $indeks = $r->rata_rata * $periode->pengali();

        return $this->sukses([
            'nomor' => $r->nomor,
            'rata' => $r->rata_rata,
            'skala_maks' => $periode->skala_maks,
            'indeks' => round($indeks, 2),
            'mutu' => Skm::mutu($indeks),
        ], 'Jawaban survei berhasil dikirim', 201);
    }
}
