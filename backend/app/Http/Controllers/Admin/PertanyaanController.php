<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\PertanyaanRequest;
use App\Http\Requests\Admin\UrutanPertanyaanRequest;
use App\Http\Resources\PertanyaanResource;
use App\Models\Pertanyaan;
use App\Services\PertanyaanService;
use Illuminate\Http\JsonResponse;

class PertanyaanController extends Controller
{
    public function __construct(private PertanyaanService $bank) {}

    public function index(): JsonResponse
    {
        return $this->sukses(PertanyaanResource::collection(Pertanyaan::berurutan()->get()));
    }

    public function store(PertanyaanRequest $request): JsonResponse
    {
        $p = $this->bank->tambah($request->validated());

        return $this->sukses(new PertanyaanResource($p), 'Pertanyaan ditambahkan', 201);
    }

    public function update(PertanyaanRequest $request, Pertanyaan $pertanyaan): JsonResponse
    {
        $pertanyaan->update($request->validated());

        return $this->sukses(new PertanyaanResource($pertanyaan), 'Pertanyaan diperbarui');
    }

    /** Soft delete: jawaban lama tetap terhubung ke pertanyaan ini untuk rekap periode sebelumnya. */
    public function destroy(Pertanyaan $pertanyaan): JsonResponse
    {
        $pertanyaan->delete();

        return $this->sukses(null, 'Pertanyaan dihapus');
    }

    public function urutan(UrutanPertanyaanRequest $request): JsonResponse
    {
        $this->bank->urutkan($request->validated('ids'));

        return $this->sukses(PertanyaanResource::collection(Pertanyaan::berurutan()->get()), 'Urutan pertanyaan disimpan');
    }
}
