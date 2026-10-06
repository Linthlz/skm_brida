<?php

use App\Http\Controllers\Admin;
use App\Http\Controllers\AuthController;
use App\Http\Controllers\Publik;
use Illuminate\Support\Facades\Route;

/*
 * API v1 SKM BRIDA Kabupaten Buleleng.
 * Batas laju (throttle) didefinisikan di AppServiceProvider.
 */
Route::prefix('v1')->group(function () {

    // ── Publik ────────────────────────────────────────────────────────────
    Route::middleware('throttle:publik')->group(function () {
        Route::get('referensi', Publik\ReferensiController::class);
        Route::get('survei/form', [Publik\SurveiController::class, 'form']);
        Route::get('periode', [Publik\HasilController::class, 'periode']);
        Route::get('hasil', [Publik\HasilController::class, 'index']);
        Route::get('ringkasan', [Publik\HasilController::class, 'ringkasan']);
    });
    Route::post('survei', [Publik\SurveiController::class, 'store'])->middleware('throttle:survei');
    Route::post('kontak', Publik\KontakController::class)->middleware('throttle:kontak');

    // ── Autentikasi ───────────────────────────────────────────────────────
    // Batas 5 percobaan/menit per email+IP ditangani di LoginRequest.
    Route::post('auth/login', [AuthController::class, 'login'])->middleware('throttle:publik');
    Route::middleware(['auth:sanctum', 'aktif'])->group(function () {
        Route::post('auth/logout', [AuthController::class, 'logout']);
        Route::get('auth/me', [AuthController::class, 'me']);
    });

    // ── Admin (satu peran: admin, akses penuh) ────────────────────────────
    Route::prefix('admin')->middleware(['auth:sanctum', 'aktif', 'throttle:admin'])->group(function () {
        Route::get('dashboard', [Admin\HasilController::class, 'dashboard']);
        Route::get('hasil', [Admin\HasilController::class, 'hasil']);
        Route::get('statistik', [Admin\HasilController::class, 'statistik']);

        Route::get('responden', [Admin\RespondenController::class, 'index']);
        Route::get('responden/{responden}', [Admin\RespondenController::class, 'show']);

        Route::put('pertanyaan/urutan', [Admin\PertanyaanController::class, 'urutan']);
        Route::apiResource('pertanyaan', Admin\PertanyaanController::class)
            ->except('show')
            ->parameters(['pertanyaan' => 'pertanyaan']);

        Route::get('feedback', [Admin\FeedbackController::class, 'index']);
        Route::get('feedback/{feedback}', [Admin\FeedbackController::class, 'show']);
        Route::patch('feedback/{feedback}', [Admin\FeedbackController::class, 'update']);

        Route::get('pesan-kontak', [Admin\PesanKontakController::class, 'index']);
        Route::patch('pesan-kontak/{pesan}', [Admin\PesanKontakController::class, 'update']);

        Route::get('laporan/export', [Admin\LaporanController::class, 'export']);

        Route::get('pengaturan', [Admin\PengaturanController::class, 'show']);
        Route::put('pengaturan', [Admin\PengaturanController::class, 'update']);
    });
});
