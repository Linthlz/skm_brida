<?php

use App\Http\Middleware\PastikanAkunAktif;
use Illuminate\Auth\AuthenticationException;
use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;
use Symfony\Component\HttpKernel\Exception\HttpExceptionInterface;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        web: __DIR__.'/../routes/web.php',
        api: __DIR__.'/../routes/api.php',
        commands: __DIR__.'/../routes/console.php',
        health: '/up',
    )
    ->withMiddleware(function (Middleware $middleware): void {
        // Sanctum mode SPA: permintaan dari domain frontend memakai sesi cookie + CSRF.
        $middleware->statefulApi();
        $middleware->alias(['aktif' => PastikanAkunAktif::class]);
        // API tidak punya halaman login; tamu cukup mendapat 401 JSON.
        $middleware->redirectGuestsTo(fn (Request $r) => $r->is('api/*') ? null : '/');
    })
    ->withExceptions(function (Exceptions $exceptions): void {
        $exceptions->shouldRenderJsonWhen(
            fn (Request $request) => $request->is('api/*') || $request->expectsJson(),
        );

        // Semua galat API berbentuk { success: false, message, errors }.
        $exceptions->render(function (Throwable $e, Request $request) {
            if (! ($request->is('api/*') || $request->expectsJson())) {
                return null;
            }

            if ($e instanceof ValidationException) {
                $status = $e->status;

                return response()->json([
                    'success' => false,
                    'message' => $status === 422 ? 'Validasi gagal' : collect($e->errors())->flatten()->first(),
                    'errors' => $e->errors(),
                ], $status);
            }

            if ($e instanceof AuthenticationException) {
                return response()->json([
                    'success' => false, 'message' => 'Silakan masuk terlebih dahulu.', 'errors' => (object) [],
                ], 401);
            }

            $status = $e instanceof HttpExceptionInterface ? $e->getStatusCode() : 500;
            $bawaan = [
                400 => 'Permintaan tidak valid.',
                403 => 'Anda tidak memiliki akses.',
                404 => 'Data tidak ditemukan.',
                405 => 'Metode HTTP tidak diizinkan.',
                409 => 'Terjadi konflik data.',
                419 => 'Sesi kedaluwarsa. Muat ulang halaman lalu coba lagi.',
                429 => 'Terlalu banyak permintaan. Coba lagi beberapa saat lagi.',
                500 => 'Terjadi kesalahan pada server.',
            ];

            // Pesan dari kode kita sendiri (abort/throw) dipakai untuk 400/403/409;
            // selain itu memakai pesan baku agar detail internal (nama model, rute) tidak bocor.
            $pesan = $bawaan[$status] ?? 'Terjadi kesalahan.';
            if (in_array($status, [400, 403, 409], true) && $e->getMessage() && $e->getMessage() !== 'This action is unauthorized.') {
                $pesan = $e->getMessage();
            }

            $body = ['success' => false, 'message' => $pesan, 'errors' => (object) []];
            if ($status === 500 && config('app.debug')) {
                $body['debug'] = $e->getMessage();
            }

            return response()->json($body, $status, $e instanceof HttpExceptionInterface ? $e->getHeaders() : []);
        });
    })->create();
