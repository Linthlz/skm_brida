<?php

namespace App\Services;

use App\Models\Pertanyaan;
use Illuminate\Support\Facades\DB;
use Illuminate\Validation\ValidationException;

class PertanyaanService
{
    /** Pertanyaan baru mendapat kode U{n} berikutnya (tidak memakai ulang kode terhapus) dan urutan paling akhir. */
    public function tambah(array $data): Pertanyaan
    {
        return DB::transaction(function () use ($data) {
            $nomorTerakhir = Pertanyaan::withTrashed()->lockForUpdate()->pluck('kode')
                ->map(fn ($k) => (int) preg_replace('/\D/', '', $k))
                ->max() ?? 0;

            return Pertanyaan::create([
                ...$data,
                'kode' => 'U'.($nomorTerakhir + 1),
                'urutan' => (int) Pertanyaan::max('urutan') + 1,
                'aktif' => true,
            ]);
        });
    }

    /**
     * Menyimpan urutan baru. $ids harus berisi seluruh pertanyaan (yang belum
     * dihapus) tepat satu kali.
     *
     * @param  int[]  $ids
     */
    public function urutkan(array $ids): void
    {
        $semua = Pertanyaan::pluck('id')->sort()->values()->all();
        $kiriman = collect($ids)->map(fn ($id) => (int) $id)->sort()->values()->all();

        if ($semua !== $kiriman) {
            throw ValidationException::withMessages([
                'ids' => 'Daftar urutan harus memuat seluruh pertanyaan tepat satu kali.',
            ]);
        }

        DB::transaction(function () use ($ids) {
            foreach (array_values($ids) as $i => $id) {
                Pertanyaan::whereKey($id)->update(['urutan' => $i + 1]);
            }
        });
    }
}
