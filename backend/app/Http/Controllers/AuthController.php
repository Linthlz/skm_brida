<?php

namespace App\Http\Controllers;

use App\Http\Requests\Auth\LoginRequest;
use App\Http\Resources\UserResource;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Log;

/**
 * Autentikasi Sanctum mode SPA: sesi disimpan di cookie httpOnly, sehingga
 * tidak ada token yang perlu disimpan frontend.
 */
class AuthController extends Controller
{
    public function login(LoginRequest $request): JsonResponse
    {
        $request->authenticate();
        $request->session()->regenerate();

        Log::info('Admin masuk', ['user_id' => $request->user()->id, 'ip' => $request->ip()]);

        return $this->sukses(new UserResource($request->user()), 'Berhasil masuk');
    }

    public function logout(Request $request): JsonResponse
    {
        Auth::guard('web')->logout();
        $request->session()->invalidate();
        $request->session()->regenerateToken();

        return $this->sukses(null, 'Berhasil keluar');
    }

    public function me(Request $request): JsonResponse
    {
        return $this->sukses(new UserResource($request->user()));
    }
}
