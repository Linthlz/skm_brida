<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\RespondenIndexRequest;
use App\Http\Resources\RespondenDetailResource;
use App\Http\Resources\RespondenResource;
use App\Models\Pengaturan;
use App\Models\Periode;
use App\Models\Responden;
use App\Support\Cari;
use Illuminate\Http\JsonResponse;

class RespondenController extends Controller
{
    public function index(RespondenIndexRequest $request): JsonResponse
    {
        $periode = $request->filled('periode')
            ? Periode::where('tahun', $request->integer('periode'))->firstOrFail()
            : Pengaturan::sekarang()->periodeAktif;

        $sort = RespondenIndexRequest::SORT[$request->input('sort', 'tanggal')];
        $dir = $request->input('dir', 'desc');

        $data = Responden::query()
            ->select('responden.*')
            ->join('layanan', 'layanan.id', '=', 'responden.layanan_id')
            ->with(['layanan:id,nama', 'periode:id,tahun,skala_maks'])
            ->where('responden.periode_id', $periode->id)
            ->when($request->integer('layanan_id'), fn ($q, $id) => $q->where('responden.layanan_id', $id))
            ->when($request->input('search'), fn ($q, $s) => Cari::di($q, $s, [
                'responden.nomor', 'responden.jk', 'responden.pekerjaan', 'responden.pendidikan', 'layanan.nama',
            ]))
            ->orderBy($sort, $dir)
            ->orderBy('responden.id', $dir)
            ->paginate($request->integer('per_page', 10))
            ->withQueryString();

        return $this->halaman(RespondenResource::collection($data));
    }

    public function show(Responden $responden): JsonResponse
    {
        $responden->load(['layanan:id,nama', 'periode', 'jawaban.pertanyaan', 'feedback:id,responden_id,isi']);

        return $this->sukses(new RespondenDetailResource($responden), 'Data berhasil ditemukan');
    }
}
