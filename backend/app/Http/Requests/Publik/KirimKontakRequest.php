<?php

namespace App\Http\Requests\Publik;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class KirimKontakRequest extends FormRequest
{
    protected function prepareForValidation(): void
    {
        foreach (['nama', 'pesan'] as $k) {
            if (is_string($this->input($k))) {
                $this->merge([$k => trim(strip_tags($this->input($k)))]);
            }
        }
    }

    public function rules(): array
    {
        return [
            'nama' => ['required', 'string', 'max:100'],
            'email' => ['required', 'string', 'email:rfc', 'max:150'],
            'topik' => ['required', Rule::in(config('skm.topik_kontak'))],
            'pesan' => ['required', 'string', 'min:10', 'max:2000'],
            'website' => ['prohibited'],
        ];
    }

    public function messages(): array
    {
        return [
            'pesan.min' => 'Tuliskan pesan minimal 10 karakter agar kami dapat menindaklanjuti.',
            'email.email' => 'Masukkan alamat email yang valid, contoh: nama@email.com',
        ];
    }
}
