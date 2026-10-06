<?php

namespace App\Http\Controllers\Publik;

use App\Http\Controllers\Controller;
use App\Models\Layanan;
use Illuminate\Http\JsonResponse;

class ReferensiController extends Controller
{
    /** Pilihan isian form survei, form kontak, dan filter. */
    public function __invoke(): JsonResponse
    {
        return $this->sukses([
            'layanan' => Layanan::tersedia()->get(['id', 'nama']),
            'jk' => config('skm.jk'),
            'pendidikan' => config('skm.pendidikan'),
            'pekerjaan' => config('skm.pekerjaan'),
            'frekuensi' => config('skm.frekuensi'),
            'kecamatan' => config('skm.kecamatan'),
            'topik_kontak' => config('skm.topik_kontak'),
            'status_feedback' => config('skm.status_feedback'),
        ]);
    }
}
