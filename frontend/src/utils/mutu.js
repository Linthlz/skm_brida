/**
 * Kategori mutu pelayanan mengacu Permen PANRB 14/2017. Dipakai untuk warna
 * badge; nilai indeksnya sendiri dihitung backend (pengali 100 / skala_maks,
 * yaitu 20 untuk skala 1–5 dan 25 untuk skala 1–4). Ambang mutu sama untuk keduanya.
 */

export function mutu(ikm) {
  if (ikm >= 88.31) return { huruf: 'A', label: 'Sangat Baik', tone: 'ok' };
  if (ikm >= 76.61) return { huruf: 'B', label: 'Baik', tone: 'ok' };
  if (ikm >= 65.0) return { huruf: 'C', label: 'Kurang Baik', tone: 'warn' };
  return { huruf: 'D', label: 'Tidak Baik', tone: 'bad' };
}

export const BAND_MUTU = [
  { huruf: 'A', rentang: '88,31 - 100', label: 'Sangat Baik', warna: 'var(--s5)', ink: 'var(--s5i)' },
  { huruf: 'B', rentang: '76,61 - 88,30', label: 'Baik', warna: 'var(--s4)', ink: 'var(--s4i)' },
  { huruf: 'C', rentang: '65,00 - 76,60', label: 'Kurang Baik', warna: 'var(--s2)', ink: 'var(--s2i)' },
  { huruf: 'D', rentang: '25,00 - 64,99', label: 'Tidak Baik', warna: 'var(--s1)', ink: 'var(--s1i)' },
];
