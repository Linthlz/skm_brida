<?php

namespace App\Http\Controllers\Publik;

use App\Http\Controllers\Controller;
use App\Http\Requests\Publik\KirimKontakRequest;
use App\Models\PesanKontak;
use App\Support\Skm;
use Illuminate\Http\JsonResponse;

class KontakController extends Controller
{
    public function __invoke(KirimKontakRequest $request): JsonResponse
    {
        $pesan = new PesanKontak($request->safe()->only(['nama', 'email', 'topik', 'pesan']));
        $pesan->ip_hash = Skm::ipHash($request->ip());
        $pesan->save();

        return $this->sukses(null, 'Pesan terkirim ke tim pelayanan BRIDA', 201);
    }
}
