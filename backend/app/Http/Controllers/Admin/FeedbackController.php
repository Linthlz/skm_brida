<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\FeedbackIndexRequest;
use App\Http\Requests\Admin\FeedbackUpdateRequest;
use App\Http\Resources\FeedbackResource;
use App\Models\Feedback;
use App\Services\FeedbackService;
use App\Support\Cari;
use Illuminate\Http\JsonResponse;

class FeedbackController extends Controller
{
    public function __construct(private FeedbackService $feedback) {}

    public function index(FeedbackIndexRequest $request): JsonResponse
    {
        $data = Feedback::query()
            ->with(['responden:id,nomor,layanan_id', 'responden.layanan:id,nama', 'kategori:id,unsur', 'penangan:id,nama'])
            ->when($request->input('status'), fn ($q, $s) => $q->where('status', $s))
            ->when($request->input('search'), fn ($q, $s) => Cari::di($q, $s, ['isi', 'catatan']))
            ->latest()
            ->orderByDesc('id')
            ->paginate($request->integer('per_page', 10))
            ->withQueryString();

        return $this->halaman(FeedbackResource::collection($data), 'Data berhasil diambil', [
            'hitung' => $this->feedback->hitungStatus(),
        ]);
    }

    public function show(Feedback $feedback): JsonResponse
    {
        $feedback->load(['responden.layanan', 'kategori', 'penangan']);

        return $this->sukses(new FeedbackResource($feedback), 'Data berhasil ditemukan');
    }

    public function update(FeedbackUpdateRequest $request, Feedback $feedback): JsonResponse
    {
        $hasil = $this->feedback->perbarui($feedback, $request->validated(), $request->user());

        return $this->sukses(new FeedbackResource($hasil), 'Feedback diperbarui');
    }
}
