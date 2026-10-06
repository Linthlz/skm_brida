<?php

namespace App\Support;

use Illuminate\Database\Eloquent\Builder;

/** Pencarian teks sederhana dengan parameter binding dan wildcard yang di-escape. */
final class Cari
{
    /** @param  string[]  $kolom  nama kolom tetap dari kode, bukan dari input pengguna */
    public static function di(Builder $q, string $kata, array $kolom): Builder
    {
        $pola = '%'.addcslashes(trim($kata), '\\%_').'%';

        return $q->where(function (Builder $w) use ($kolom, $pola) {
            foreach ($kolom as $k) {
                $w->orWhere($k, 'like', $pola);
            }
        });
    }
}
