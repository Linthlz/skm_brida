<?php

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;

class PengaturanRequest extends FormRequest
{
    protected function prepareForValidation(): void
    {
        if (is_string($this->input('judul'))) {
            $this->merge(['judul' => trim(strip_tags($this->input('judul')))]);
        }
    }

    public function rules(): array
    {
        return [
            'judul' => ['required', 'string', 'max:150'],
            'periode' => ['required', 'integer', 'between:2000,2100'],
            'skala_maks' => ['required', 'integer', 'in:4,5'],
            'tanggal_tutup' => ['nullable', 'date_format:Y-m-d'],
            'anonim' => ['required', 'boolean'],
            'publik' => ['required', 'boolean'],
        ];
    }
}
