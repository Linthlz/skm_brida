<?php

namespace App\Http\Requests\Publik;

use App\Models\Pengaturan;
use App\Models\Pertanyaan;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;
use Illuminate\Validation\Validator;
use Symfony\Component\HttpKernel\Exception\AccessDeniedHttpException;

class KirimSurveiRequest extends FormRequest
{
    private ?Pengaturan $pengaturan = null;

    public function pengaturan(): Pengaturan
    {
        return $this->pengaturan ??= Pengaturan::sekarang();
    }

    /** Survei hanya diterima selama periode aktif belum ditutup. */
    public function authorize(): bool
    {
        return ! $this->pengaturan()->periodeAktif->ditutup();
    }

    protected function failedAuthorization(): void
    {
        throw new AccessDeniedHttpException('Periode survei sudah ditutup. Terima kasih atas minat Anda.');
    }

    protected function prepareForValidation(): void
    {
        // Teks bebas disimpan tanpa tag HTML. Tampilan tetap meng-escape keluaran.
        foreach (['nama', 'saran', 'apresiasi'] as $k) {
            if (is_string($this->input($k))) {
                $this->merge([$k => trim(strip_tags($this->input($k)))]);
            }
        }
    }

    public function rules(): array
    {
        $maks = $this->pengaturan()->periodeAktif->skala_maks;

        return [
            'nama' => [$this->pengaturan()->anonim ? 'nullable' : 'required', 'string', 'max:100'],
            'jk' => ['required', Rule::in(config('skm.jk'))],
            'usia' => ['required', 'integer', 'between:10,99'],
            'pendidikan' => ['required', Rule::in(config('skm.pendidikan'))],
            'pekerjaan' => ['required', Rule::in(config('skm.pekerjaan'))],
            'frekuensi' => ['required', Rule::in(config('skm.frekuensi'))],
            'kecamatan' => ['nullable', Rule::in(config('skm.kecamatan'))],
            'layanan_id' => ['required', 'integer', Rule::exists('layanan', 'id')->where('aktif', true)],
            'jawaban' => ['required', 'array', 'min:1'],
            'jawaban.*.pertanyaan_id' => [
                'required', 'integer', 'distinct',
                Rule::exists('pertanyaan', 'id')->where('aktif', true)->whereNull('deleted_at'),
            ],
            'jawaban.*.nilai' => ['required', 'integer', "between:1,$maks"],
            'saran' => ['nullable', 'string', 'max:600'],
            'apresiasi' => ['nullable', 'string', 'max:300'],
            'setuju' => ['accepted'],
            // Honeypot: kolom tersembunyi yang hanya diisi bot.
            'website' => ['prohibited'],
        ];
    }

    /** Seluruh pertanyaan aktif yang wajib harus ikut dijawab. */
    public function after(): array
    {
        return [function (Validator $v) {
            if ($v->errors()->has('jawaban') || $v->errors()->has('jawaban.*')) {
                return;
            }
            $dijawab = collect($this->input('jawaban', []))->pluck('pertanyaan_id')->map(fn ($id) => (int) $id);
            $belum = Pertanyaan::aktif()->where('wajib', true)->pluck('id')->diff($dijawab)->count();
            if ($belum > 0) {
                $v->errors()->add('jawaban', "Silakan pilih salah satu jawaban sebelum melanjutkan. $belum pertanyaan belum dinilai.");
            }
        }];
    }
}
