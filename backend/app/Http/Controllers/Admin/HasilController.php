<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\FilterHasilRequest;
use App\Services\HasilService;
use Illuminate\Http\JsonResponse;

/** Dashboard, Hasil Survei, dan Statistik. Tidak terpengaruh pengaturan publikasi. */
class HasilController extends Controller
{
    public function __construct(private HasilService $hasil) {}

    public function dashboard(FilterHasilRequest $request): JsonResponse
    {
        return $this->sukses($this->hasil->dashboard($request->periode()));
    }

    public function hasil(FilterHasilRequest $request): JsonResponse
    {
        return $this->sukses($this->hasil->rekap($request->periode(), $request->filter()));
    }

    public function statistik(FilterHasilRequest $request): JsonResponse
    {
        return $this->sukses($this->hasil->statistik($request->periode()));
    }
}
