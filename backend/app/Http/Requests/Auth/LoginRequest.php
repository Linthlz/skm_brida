<?php

namespace App\Http\Requests\Auth;

use Illuminate\Auth\Events\Lockout;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\RateLimiter;
use Illuminate\Support\Str;
use Illuminate\Validation\ValidationException;

class LoginRequest extends FormRequest
{
    private const MAKS_PERCOBAAN = 5;

    public function rules(): array
    {
        return [
            'email' => ['required', 'string', 'email', 'max:150'],
            'password' => ['required', 'string', 'max:255'],
            'ingat' => ['sometimes', 'boolean'],
        ];
    }

    /**
     * Mencoba login. Pesan gagal sengaja sama untuk email tak terdaftar, sandi
     * salah, maupun akun nonaktif agar tidak membocorkan keberadaan akun.
     */
    public function authenticate(): void
    {
        $kunci = $this->kunciThrottle();

        if (RateLimiter::tooManyAttempts($kunci, self::MAKS_PERCOBAAN)) {
            event(new Lockout($this));
            $detik = RateLimiter::availableIn($kunci);
            throw ValidationException::withMessages([
                'email' => "Terlalu banyak percobaan masuk. Coba lagi dalam $detik detik.",
            ])->status(429);
        }

        $kredensial = [...$this->only('email', 'password'), 'aktif' => true];

        if (! Auth::guard('web')->attempt($kredensial, $this->boolean('ingat'))) {
            RateLimiter::hit($kunci, 60);
            throw ValidationException::withMessages([
                'email' => 'Email atau kata sandi tidak sesuai.',
            ]);
        }

        RateLimiter::clear($kunci);
    }

    private function kunciThrottle(): string
    {
        return 'login|'.Str::transliterate(Str::lower($this->string('email'))).'|'.$this->ip();
    }
}
