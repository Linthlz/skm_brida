<?php

namespace App\Providers;

use Illuminate\Cache\RateLimiting\Limit;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\RateLimiter;
use Illuminate\Support\ServiceProvider;

class AppServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        //
    }

    public function boot(): void
    {
        $this->aturBatasLaju();
    }

    private function aturBatasLaju(): void
    {
        RateLimiter::for('publik', fn (Request $r) => Limit::perMinute(60)->by($r->ip()));

        // Survei bisa diisi bergantian dari satu perangkat/IP kantor (loket layanan),
        // jadi batasnya dibuat cukup longgar (SKM_SURVEI_PER_JAM di .env).
        RateLimiter::for('survei', fn (Request $r) => Limit::perHour(config('skm.batas.survei_per_jam'))->by($r->ip()));

        RateLimiter::for('kontak', fn (Request $r) => Limit::perHour(config('skm.batas.kontak_per_jam'))->by($r->ip()));

        RateLimiter::for('admin', fn (Request $r) => Limit::perMinute(120)->by($r->user()?->id ?: $r->ip()));
    }
}
