<?php

namespace App\Exports;

use Maatwebsite\Excel\Concerns\FromArray;
use Maatwebsite\Excel\Concerns\ShouldAutoSize;
use Maatwebsite\Excel\Concerns\WithHeadings;
use Maatwebsite\Excel\Concerns\WithTitle;

/** Rekapitulasi nilai per unsur, sama dengan tabel di halaman Laporan. */
class LaporanExport implements FromArray, ShouldAutoSize, WithHeadings, WithTitle
{
    public function __construct(private array $rekap, private string $keterangan) {}

    public function title(): string
    {
        return 'SKM '.$this->rekap['periode'];
    }

    public function headings(): array
    {
        return [
            ['Rekapitulasi Survei Kepuasan Masyarakat — periode '.$this->rekap['periode']],
            [$this->keterangan.' · '.$this->rekap['total_responden'].' responden'],
            [],
            ['Kode', 'Unsur pelayanan', 'NRR', 'NRR tertimbang', 'Indeks', 'Mutu'],
        ];
    }

    public function array(): array
    {
        $baris = collect($this->rekap['unsur'])->map(fn ($u) => [
            $u['kode'], $u['nama'], $u['nrr'], $u['nrr_tertimbang'], $u['indeks'], $u['mutu']['huruf'] ?? '',
        ])->all();

        $baris[] = [
            'Total', '', $this->rekap['nrr'],
            $this->rekap['nrr'] !== null && count($this->rekap['unsur'])
                ? round($this->rekap['nrr'] / count($this->rekap['unsur']), 3) : null,
            $this->rekap['ikm'], $this->rekap['mutu']['huruf'] ?? '',
        ];

        return $baris;
    }
}
