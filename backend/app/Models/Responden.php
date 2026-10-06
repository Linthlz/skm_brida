<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Hidden;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\Relations\HasOne;

/**
 * Satu pengisian survei. Kolom sistem (periode_id, urut, nomor, rata_rata,
 * ip_hash) sengaja tidak fillable dan diisi oleh SurveiService.
 */
#[Table('responden')]
#[Fillable(['layanan_id', 'nama', 'jk', 'usia', 'pendidikan', 'pekerjaan', 'frekuensi', 'kecamatan', 'apresiasi'])]
#[Hidden(['nama', 'ip_hash'])]
class Responden extends Model
{
    protected function casts(): array
    {
        return [
            'usia' => 'integer',
            'rata_rata' => 'float',
        ];
    }

    public function periode(): BelongsTo
    {
        return $this->belongsTo(Periode::class);
    }

    public function layanan(): BelongsTo
    {
        return $this->belongsTo(Layanan::class);
    }

    public function jawaban(): HasMany
    {
        return $this->hasMany(Jawaban::class);
    }

    public function feedback(): HasOne
    {
        return $this->hasOne(Feedback::class);
    }
}
