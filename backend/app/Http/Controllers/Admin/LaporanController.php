<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\ExportLaporanRequest;
use App\Services\LaporanService;
use Symfony\Component\HttpFoundation\Response;

class LaporanController extends Controller
{
    public function __construct(private LaporanService $laporan) {}

    public function export(ExportLaporanRequest $request): Response
    {
        return $this->laporan->unduh($request->periode(), $request->filter(), $request->validated('format'));
    }
}
