<?php

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class FeedbackUpdateRequest extends FormRequest
{
    protected function prepareForValidation(): void
    {
        if (is_string($this->input('catatan'))) {
            $this->merge(['catatan' => trim(strip_tags($this->input('catatan')))]);
        }
    }

    public function rules(): array
    {
        return [
            'status' => ['sometimes', Rule::in(config('skm.status_feedback'))],
            'catatan' => ['sometimes', 'nullable', 'string', 'max:2000'],
            'kategori_id' => ['sometimes', 'nullable', 'integer', Rule::exists('pertanyaan', 'id')],
        ];
    }
}
