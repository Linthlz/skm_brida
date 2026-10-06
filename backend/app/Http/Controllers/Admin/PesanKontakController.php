<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\PesanKontakIndexRequest;
use App\Http\Requests\Admin\PesanKontakUpdateRequest;
use App\Http\Resources\PesanKontakResource;
use App\Models\PesanKontak;
use Illuminate\Http\JsonResponse;

class PesanKontakController extends Controller
{
    public function index(PesanKontakIndexRequest $request): JsonResponse
    {
        $data = PesanKontak::query()
            ->when($request->has('dibaca'), fn ($q) => $q->where('dibaca', $request->boolean('dibaca')))
            ->latest()
            ->orderByDesc('id')
            ->paginate($request->integer('per_page', 10))
            ->withQueryString();

        return $this->halaman(PesanKontakResource::collection($data), 'Data berhasil diambil', [
            'belum_dibaca' => PesanKontak::where('dibaca', false)->count(),
        ]);
    }

    public function update(PesanKontakUpdateRequest $request, PesanKontak $pesan): JsonResponse
    {
        $pesan->dibaca = $request->boolean('dibaca');
        $pesan->save();

        return $this->sukses(new PesanKontakResource($pesan), 'Pesan diperbarui');
    }
}
