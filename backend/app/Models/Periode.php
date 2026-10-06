<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;

#[Table('periode')]
#[Fillable(['tahun', 'skala_maks', 'tanggal_tutup'])]
class Periode extends Model
{
    protected function casts(): array
    {
        return [
            'tahun' => 'integer',
            'skala_maks' => 'integer',
            'tanggal_tutup' => 'date',
        ];
    }

    public function responden(): HasMany
    {
        return $this->hasMany(Responden::class);
    }

    /** Pengali konversi NRR ke indeks: 20 untuk skala 1–5, 25 untuk skala 1–4. */
    public function pengali(): float
    {
        return 100 / $this->skala_maks;
    }

    /** Survei ditutup setelah hari tanggal_tutup berakhir. */
    public function ditutup(): bool
    {
        return $this->tanggal_tutup !== null && now()->startOfDay()->gt($this->tanggal_tutup);
    }
}
