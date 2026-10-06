<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\PengaturanRequest;
use App\Services\PengaturanService;
use Illuminate\Http\JsonResponse;

class PengaturanController extends Controller
{
    public function __construct(private PengaturanService $pengaturan) {}

    public function show(): JsonResponse
    {
        return $this->sukses($this->pengaturan->tampil());
    }

    public function update(PengaturanRequest $request): JsonResponse
    {
        return $this->sukses($this->pengaturan->simpan($request->validated()), 'Pengaturan disimpan');
    }
}
