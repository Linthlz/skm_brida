<?php

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;

/** Dipakai untuk POST (semua wajib) dan PATCH (field boleh sebagian). */
class PertanyaanRequest extends FormRequest
{
    protected function prepareForValidation(): void
    {
        foreach (['unsur', 'teks'] as $k) {
            if (is_string($this->input($k))) {
                $this->merge([$k => trim(strip_tags($this->input($k)))]);
            }
        }
    }

    public function rules(): array
    {
        $wajib = $this->isMethod('post') ? 'required' : 'sometimes';

        return [
            'unsur' => [$wajib, 'string', 'max:100'],
            'teks' => [$wajib, 'string', 'max:500'],
            'wajib' => ['sometimes', 'boolean'],
            'aktif' => ['sometimes', 'boolean'],
        ];
    }
}
