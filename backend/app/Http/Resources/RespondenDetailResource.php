<?php

namespace App\Http\Resources;

use Illuminate\Http\Request;

/** @mixin \App\Models\Responden */
class RespondenDetailResource extends RespondenResource
{
    public function toArray(Request $request): array
    {
        return [
            ...parent::toArray($request),
            'kecamatan' => $this->kecamatan,
            'frekuensi' => $this->frekuensi,
            'apresiasi' => $this->apresiasi,
            'skala_maks' => $this->periode->skala_maks,
            'indeks' => round($this->rata_rata * $this->periode->pengali(), 2),
            'jawaban' => $this->jawaban->map(fn ($j) => [
                'kode' => $j->pertanyaan->kode,
                'unsur' => $j->pertanyaan->unsur,
                'nilai' => $j->nilai,
            ])->sortBy('kode', SORT_NATURAL)->values(),
            'saran' => $this->feedback?->isi,
        ];
    }
}
