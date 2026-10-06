<?php

namespace App\Http\Controllers;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;
use Illuminate\Pagination\LengthAwarePaginator;

/**
 * Seluruh respons API berbentuk { success, message, data, meta? }.
 * Respons galat dibentuk terpusat di bootstrap/app.php.
 */
abstract class Controller
{
    protected function sukses(mixed $data = null, string $pesan = 'Data berhasil diambil', int $status = 200, ?array $meta = null): JsonResponse
    {
        $body = ['success' => true, 'message' => $pesan, 'data' => $data];
        if ($meta !== null) {
            $body['meta'] = $meta;
        }

        return response()->json($body, $status);
    }

    /** Koleksi resource berpaginasi; $metaTambahan digabung ke meta paginasi. */
    protected function halaman(AnonymousResourceCollection $koleksi, string $pesan = 'Data berhasil diambil', array $metaTambahan = []): JsonResponse
    {
        /** @var LengthAwarePaginator $p */
        $p = $koleksi->resource;

        return $this->sukses($koleksi->resolve(), $pesan, 200, [
            'current_page' => $p->currentPage(),
            'last_page' => $p->lastPage(),
            'per_page' => $p->perPage(),
            'total' => $p->total(),
            ...$metaTambahan,
        ]);
    }
}
