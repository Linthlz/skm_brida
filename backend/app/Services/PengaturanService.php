<?php

namespace App\Services;

use App\Models\Pengaturan;
use App\Models\Periode;
use App\Models\User;
use Illuminate\Support\Facades\DB;
use Symfony\Component\HttpKernel\Exception\ConflictHttpException;

class PengaturanService
{
    public function tampil(): array
    {
        $p = Pengaturan::sekarang();
        $periode = $p->periodeAktif;

        return [
            'judul' => $p->judul,
            'periode' => $periode->tahun,
            'skala_maks' => $periode->skala_maks,
            'skala_terkunci' => $periode->responden()->exists(),
            'tanggal_tutup' => $periode->tanggal_tutup?->toDateString(),
            'anonim' => $p->anonim,
            'publik' => $p->publik,
            'opsi_periode' => Periode::orderByDesc('tahun')->pluck('tahun'),
            'pengelola' => User::where('aktif', true)->orderBy('id')->get(['id', 'nama', 'jabatan'])
                ->map(fn ($u) => ['id' => $u->id, 'nama' => $u->nama, 'jabatan' => $u->jabatan, 'inisial' => $u->inisial()]),
        ];
    }

    /**
     * Skala dan tanggal tutup melekat pada periode aktif. Skala tidak boleh diubah
     * bila periode itu sudah punya responden, karena akan mengacaukan indeks.
     */
    public function simpan(array $data): array
    {
        DB::transaction(function () use ($data) {
            $periode = Periode::firstOrCreate(['tahun' => $data['periode']], ['skala_maks' => $data['skala_maks']]);

            if ($periode->skala_maks !== (int) $data['skala_maks'] && $periode->responden()->exists()) {
                throw new ConflictHttpException(
                    'Skala penilaian periode '.$periode->tahun.' tidak dapat diubah karena sudah ada responden.'
                );
            }

            $periode->update([
                'skala_maks' => $data['skala_maks'],
                'tanggal_tutup' => $data['tanggal_tutup'] ?? null,
            ]);

            Pengaturan::sekarang()->update([
                'judul' => $data['judul'],
                'periode_aktif_id' => $periode->id,
                'anonim' => $data['anonim'],
                'publik' => $data['publik'],
            ]);
        });

        return $this->tampil();
    }
}
