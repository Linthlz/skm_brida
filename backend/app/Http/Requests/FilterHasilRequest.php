<?php

namespace App\Http\Requests;

use App\Models\Pengaturan;
use App\Models\Periode;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

/** Filter bersama untuk endpoint hasil, dashboard, statistik, dan laporan. */
class FilterHasilRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'periode' => ['nullable', 'integer', Rule::exists('periode', 'tahun')],
            'layanan_id' => ['nullable', 'integer', Rule::exists('layanan', 'id')],
            'pendidikan' => ['nullable', Rule::in(config('skm.pendidikan'))],
        ];
    }

    /** Periode yang diminta, atau periode aktif bila tidak disebut. */
    public function periode(): Periode
    {
        return $this->filled('periode')
            ? Periode::where('tahun', $this->integer('periode'))->firstOrFail()
            : Pengaturan::sekarang()->periodeAktif;
    }

    public function filter(): array
    {
        return array_filter($this->safe()->only(['layanan_id', 'pendidikan']), fn ($v) => $v !== null && $v !== '');
    }
}
