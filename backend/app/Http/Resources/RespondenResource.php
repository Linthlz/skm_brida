<?php

namespace App\Http\Resources;

use App\Support\Skm;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/**
 * Baris tabel Data Responden / Data Survei. Nama responden tidak pernah
 * dikirim (lihat #[Hidden] di model).
 *
 * @mixin \App\Models\Responden
 */
class RespondenResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'nomor' => $this->nomor,
            'usia' => $this->usia,
            'jk' => $this->jk,
            'pendidikan' => $this->pendidikan,
            'pekerjaan' => $this->pekerjaan,
            'layanan' => $this->layanan?->nama,
            'nilai' => $this->rata_rata,
            'kategori' => Skm::kategori($this->rata_rata, $this->periode->skala_maks),
            'tanggal' => $this->created_at?->toIso8601String(),
        ];
    }
}
