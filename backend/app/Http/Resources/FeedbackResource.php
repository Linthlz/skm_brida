<?php

namespace App\Http\Resources;

use App\Models\Feedback;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin \App\Models\Feedback */
class FeedbackResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'nomor' => $this->nomor(),
            'tanggal' => $this->created_at?->toIso8601String(),
            'isi' => $this->isi,
            'status' => $this->status,
            'status_berikut' => Feedback::BERIKUT[$this->status] ?? null,
            'kategori' => $this->kategori ? ['id' => $this->kategori->id, 'nama' => $this->kategori->unsur] : null,
            'layanan' => $this->responden?->layanan?->nama,
            'nomor_responden' => $this->responden?->nomor,
            'catatan' => $this->catatan,
            'ditangani_oleh' => $this->penangan?->nama,
            'diperbarui' => $this->updated_at?->toIso8601String(),
        ];
    }
}
