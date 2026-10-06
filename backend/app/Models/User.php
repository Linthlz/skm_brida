<?php

namespace App\Models;

use Database\Factories\UserFactory;
use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Hidden;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Foundation\Auth\User as Authenticatable;
use Illuminate\Notifications\Notifiable;

/**
 * Pengelola panel SKM. Hanya ada satu peran (admin, akses penuh); `jabatan`
 * sekadar label tampilan. Akun dengan `aktif = false` ditolak middleware.
 */
#[Fillable(['nama', 'email', 'password', 'jabatan', 'aktif'])]
#[Hidden(['password', 'remember_token'])]
class User extends Authenticatable
{
    /** @use HasFactory<UserFactory> */
    use HasFactory, Notifiable;

    protected function casts(): array
    {
        return [
            'aktif' => 'boolean',
            'password' => 'hashed',
        ];
    }

    /** Dua huruf awal nama untuk avatar, contoh "Ketut Pramana" -> "KP". */
    public function inisial(): string
    {
        $kata = preg_split('/\s+/', trim($this->nama)) ?: [];

        return mb_strtoupper(collect($kata)->take(2)->map(fn ($k) => mb_substr($k, 0, 1))->implode(''));
    }
}
