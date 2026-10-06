<?php

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class RespondenIndexRequest extends FormRequest
{
    /** Kolom yang boleh dipakai untuk pengurutan (nama dari frontend => kolom DB). */
    public const SORT = [
        'nomor' => 'responden.urut',
        'usia' => 'responden.usia',
        'jk' => 'responden.jk',
        'pekerjaan' => 'responden.pekerjaan',
        'pendidikan' => 'responden.pendidikan',
        'layanan' => 'layanan.nama',
        'nilai' => 'responden.rata_rata',
        'tanggal' => 'responden.created_at',
    ];

    public function rules(): array
    {
        return [
            'periode' => ['nullable', 'integer', Rule::exists('periode', 'tahun')],
            'layanan_id' => ['nullable', 'integer', Rule::exists('layanan', 'id')],
            'search' => ['nullable', 'string', 'max:100'],
            'sort' => ['nullable', Rule::in(array_keys(self::SORT))],
            'dir' => ['nullable', 'in:asc,desc'],
            'page' => ['nullable', 'integer', 'min:1'],
            'per_page' => ['nullable', 'integer', 'between:5,100'],
        ];
    }
}
