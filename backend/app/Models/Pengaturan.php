<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

/** Pengaturan survei. Tabelnya hanya berisi satu baris. */
#[Table('pengaturan')]
#[Fillable(['judul', 'periode_aktif_id', 'anonim', 'publik'])]
class Pengaturan extends Model
{
    protected function casts(): array
    {
        return [
            'anonim' => 'boolean',
            'publik' => 'boolean',
        ];
    }

    public function periodeAktif(): BelongsTo
    {
        return $this->belongsTo(Periode::class, 'periode_aktif_id');
    }

    public static function sekarang(): self
    {
        return static::with('periodeAktif')->firstOrFail();
    }
}
