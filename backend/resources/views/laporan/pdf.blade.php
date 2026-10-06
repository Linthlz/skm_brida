<!doctype html>
<html lang="id">
<head>
<meta charset="utf-8">
<title>Laporan SKM {{ $rekap['periode'] }}</title>
<style>
  body { font-family: DejaVu Sans, sans-serif; font-size: 11px; color: #1f1a17; }
  h1 { font-size: 16px; margin: 0 0 2px; }
  .muted { color: #6b625c; }
  table { width: 100%; border-collapse: collapse; margin-top: 14px; }
  th, td { border: 1px solid #d9d0c9; padding: 6px 8px; }
  th { background: #f3ede8; text-align: left; font-size: 10px; text-transform: uppercase; }
  td.num { text-align: right; font-family: DejaVu Sans Mono, monospace; }
  tr.total td { font-weight: bold; background: #f7ecec; }
</style>
</head>
<body>
  <h1>{{ $judul }}</h1>
  <div class="muted">Rekapitulasi hasil — periode {{ $rekap['periode'] }} · {{ $keterangan }}</div>
  <div class="muted">
    {{ number_format($rekap['total_responden'], 0, ',', '.') }} responden ·
    IKM {{ $rekap['ikm'] !== null ? number_format($rekap['ikm'], 2, ',', '.') : '–' }}
    @if ($rekap['mutu']) · Mutu {{ $rekap['mutu']['huruf'] }} ({{ $rekap['mutu']['label'] }}) @endif
  </div>

  <table>
    <thead>
      <tr><th>Kode</th><th>Unsur pelayanan</th><th>NRR</th><th>NRR tertimbang</th><th>Indeks</th><th>Mutu</th></tr>
    </thead>
    <tbody>
      @forelse ($rekap['unsur'] as $u)
        <tr>
          <td>{{ $u['kode'] }}</td>
          <td>{{ $u['nama'] }}</td>
          <td class="num">{{ number_format($u['nrr'], 2, ',', '.') }}</td>
          <td class="num">{{ number_format($u['nrr_tertimbang'], 3, ',', '.') }}</td>
          <td class="num">{{ number_format($u['indeks'], 2, ',', '.') }}</td>
          <td>{{ $u['mutu']['huruf'] ?? '' }}</td>
        </tr>
      @empty
        <tr><td colspan="6" class="muted">Belum ada data pada filter ini.</td></tr>
      @endforelse
      @if (count($rekap['unsur']))
        <tr class="total">
          <td colspan="2">Total</td>
          <td class="num">{{ number_format($rekap['nrr'], 2, ',', '.') }}</td>
          <td class="num">{{ number_format($rekap['nrr'] / count($rekap['unsur']), 3, ',', '.') }}</td>
          <td class="num">{{ number_format($rekap['ikm'], 2, ',', '.') }}</td>
          <td>{{ $rekap['mutu']['huruf'] ?? '' }}</td>
        </tr>
      @endif
    </tbody>
  </table>

  <p class="muted" style="margin-top:14px">
    Skala 1–{{ $rekap['skala_maks'] }}, indeks = NRR × {{ rtrim(rtrim(number_format($rekap['pengali'], 2, ',', ''), '0'), ',') }}.
    Dicetak {{ now()->locale('id')->translatedFormat('d F Y H:i') }}.
  </p>
</body>
</html>
