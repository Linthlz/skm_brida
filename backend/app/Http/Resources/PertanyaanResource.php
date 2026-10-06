<?php

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin \App\Models\Pertanyaan */
class PertanyaanResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'kode' => $this->kode,
            'unsur' => $this->unsur,
            'teks' => $this->teks,
            'wajib' => $this->wajib,
            'aktif' => $this->aktif,
            'urutan' => $this->urutan,
        ];
    }
}
