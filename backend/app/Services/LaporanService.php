<?php

namespace App\Services;

use App\Exports\LaporanExport;
use App\Models\Layanan;
use App\Models\Pengaturan;
use App\Models\Periode;
use Barryvdh\DomPDF\Facade\Pdf;
use Maatwebsite\Excel\Facades\Excel;
use Symfony\Component\HttpFoundation\Response;

class LaporanService
{
    public function __construct(private HasilService $hasil) {}

    public function unduh(Periode $periode, array $filter, string $format): Response
    {
        $rekap = $this->hasil->rekap($periode, $filter);
        $keterangan = $this->keterangan($filter);
        $nama = 'laporan-skm-'.$periode->tahun;

        if ($format === 'xlsx') {
            return Excel::download(new LaporanExport($rekap, $keterangan), $nama.'.xlsx');
        }

        return Pdf::loadView('laporan.pdf', [
            'rekap' => $rekap,
            'keterangan' => $keterangan,
            'judul' => Pengaturan::sekarang()->judul,
        ])->setPaper('a4')->download($nama.'.pdf');
    }

    private function keterangan(array $filter): string
    {
        $layanan = isset($filter['layanan_id']) ? Layanan::find($filter['layanan_id'])?->nama : null;

        return ($layanan ?? 'Semua layanan').' · '.($filter['pendidikan'] ?? 'Semua jenjang pendidikan');
    }
}
