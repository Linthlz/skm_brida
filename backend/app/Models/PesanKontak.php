<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Hidden;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Model;

#[Table('pesan_kontak')]
#[Fillable(['nama', 'email', 'topik', 'pesan'])]
#[Hidden(['ip_hash'])]
class PesanKontak extends Model
{
    protected function casts(): array
    {
        return ['dibaca' => 'boolean'];
    }
}
