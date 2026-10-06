<?php

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;

class PesanKontakUpdateRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'dibaca' => ['required', 'boolean'],
        ];
    }
}
