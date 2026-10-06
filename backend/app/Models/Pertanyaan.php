<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\SoftDeletes;

#[Table('pertanyaan')]
#[Fillable(['kode', 'unsur', 'teks', 'wajib', 'aktif', 'urutan'])]
class Pertanyaan extends Model
{
    use SoftDeletes;

    protected function casts(): array
    {
        return [
            'wajib' => 'boolean',
            'aktif' => 'boolean',
            'urutan' => 'integer',
        ];
    }

    public function jawaban(): HasMany
    {
        return $this->hasMany(Jawaban::class);
    }

    public function scopeBerurutan(Builder $q): Builder
    {
        return $q->orderBy('urutan')->orderBy('id');
    }

    public function scopeAktif(Builder $q): Builder
    {
        return $q->where('aktif', true);
    }
}
