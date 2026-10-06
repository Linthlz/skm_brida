/** Angka desimal format Indonesia. null/undefined ditampilkan sebagai tanda pisah. */
export const fmt = (n, d = 2) =>
  n === null || n === undefined ? '–' : Number(n).toLocaleString('id-ID', { minimumFractionDigits: d, maximumFractionDigits: d });

export const fmtInt = (n) => (n === null || n === undefined ? '–' : Number(n).toLocaleString('id-ID'));

/** Tanggal pendek Indonesia; menerima Date atau string ISO dari API. */
export const tglID = (d) =>
  d ? new Date(d).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }) : '–';

/** Baris unsur dari API ({kode, nama, nrr}) ke bentuk yang dipakai BarUnsur ({kode, nama, nilai}). */
export const unsurGrafik = (unsur = []) => unsur.map((u) => ({ kode: u.kode, nama: u.nama, nilai: u.nrr }));
