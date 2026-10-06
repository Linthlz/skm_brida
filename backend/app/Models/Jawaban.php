<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Attributes\WithoutTimestamps;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

#[Table('jawaban')]
#[WithoutTimestamps]
#[Fillable(['pertanyaan_id', 'nilai'])]
class Jawaban extends Model
{
    protected function casts(): array
    {
        return ['nilai' => 'integer'];
    }

    public function responden(): BelongsTo
    {
        return $this->belongsTo(Responden::class);
    }

    public function pertanyaan(): BelongsTo
    {
        return $this->belongsTo(Pertanyaan::class)->withTrashed();
    }
}
