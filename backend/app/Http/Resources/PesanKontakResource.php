<?php

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin \App\Models\PesanKontak */
class PesanKontakResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'nama' => $this->nama,
            'email' => $this->email,
            'topik' => $this->topik,
            'pesan' => $this->pesan,
            'dibaca' => $this->dibaca,
            'tanggal' => $this->created_at?->toIso8601String(),
        ];
    }
}
