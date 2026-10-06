<?php

namespace App\Support;

/**
 * Aturan domain SKM yang dipakai bersama oleh service, resource, dan ekspor.
 * Padanannya di frontend: src/utils/mutu.js dan src/data/skala.js.
 */
final class Skm
{
    /** Kategori mutu pelayanan untuk sebuah nilai indeks (0–100). */
    public static function mutu(?float $indeks): ?array
    {
        if ($indeks === null) {
            return null;
        }
        foreach (config('skm.mutu') as $m) {
            if ($indeks >= $m['min']) {
                return ['huruf' => $m['huruf'], 'label' => $m['label']];
            }
        }

        return null;
    }

    /** Daftar tingkat skala, contoh [['v' => 1, 'label' => 'Sangat Tidak Puas'], …]. */
    public static function skala(int $maks): array
    {
        return collect(config("skm.label_skala.$maks"))
            ->map(fn ($label, $v) => ['v' => $v, 'label' => $label])
            ->values()
            ->all();
    }

    /** Label kepuasan seorang responden: nilai rata-ratanya dibulatkan ke tingkat terdekat. */
    public static function kategori(float $rata, int $maks): string
    {
        $v = max(1, min($maks, (int) round($rata)));

        return config("skm.label_skala.$maks.$v");
    }

    public static function ipHash(?string $ip): ?string
    {
        return $ip ? hash('sha256', config('skm.ip_salt').'|'.$ip) : null;
    }

    public static function bulat(?float $n, int $d = 2): ?float
    {
        return $n === null ? null : round($n, $d);
    }
}
