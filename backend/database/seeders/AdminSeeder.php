<?php

namespace Database\Seeders;

use App\Models\User;
use Illuminate\Database\Seeder;

/**
 * Akun admin awal dari ADMIN_EMAIL / ADMIN_PASSWORD di .env. Tidak ada kata
 * sandi bawaan di kode; bila variabel kosong, seeder dilewati.
 */
class AdminSeeder extends Seeder
{
    public function run(): void
    {
        $cfg = config('skm.admin');

        if (blank($cfg['email']) || blank($cfg['password'])) {
            $this->command?->warn('ADMIN_EMAIL / ADMIN_PASSWORD belum diisi di .env — akun admin tidak dibuat.');

            return;
        }
        if (strlen($cfg['password']) < 12) {
            $this->command?->warn('ADMIN_PASSWORD minimal 12 karakter — akun admin tidak dibuat.');

            return;
        }

        User::firstOrCreate(['email' => $cfg['email']], [
            'nama' => $cfg['nama'],
            'password' => $cfg['password'],
            'jabatan' => $cfg['jabatan'],
            'aktif' => true,
        ]);
    }
}
