<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;

#[Table('layanan')]
#[Fillable(['nama', 'aktif', 'urutan'])]
class Layanan extends Model
{
    protected function casts(): array
    {
        return ['aktif' => 'boolean'];
    }

    public function responden(): HasMany
    {
        return $this->hasMany(Responden::class);
    }

    public function scopeTersedia(Builder $q): Builder
    {
        return $q->where('aktif', true)->orderBy('urutan')->orderBy('id');
    }
}
