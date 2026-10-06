<?php

namespace App\Http\Requests\Admin;

use App\Http\Requests\FilterHasilRequest;

class ExportLaporanRequest extends FilterHasilRequest
{
    public function rules(): array
    {
        return [
            ...parent::rules(),
            'format' => ['required', 'in:xlsx,pdf'],
        ];
    }
}
