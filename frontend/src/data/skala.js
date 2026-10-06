/** Skala penilaian 1-5. Warna divergen: merah (negatif) -> toska (positif).
 *  `ink` adalah warna teks yang kontras di atas warna tersebut, untuk tema terang dan gelap. */
export const SKALA = [
  {v:1, label:'Sangat Tidak Puas', pendek:'Sangat Tidak Puas', warna:'var(--s1)', ink:'var(--s1i)'},
  {v:2, label:'Tidak Puas',        pendek:'Tidak Puas',        warna:'var(--s2)', ink:'var(--s2i)'},
  {v:3, label:'Kurang Puas',       pendek:'Kurang Puas',       warna:'var(--s3)', ink:'var(--s3i)'},
  {v:4, label:'Puas',              pendek:'Puas',              warna:'var(--s4)', ink:'var(--s4i)'},
  {v:5, label:'Sangat Puas',       pendek:'Sangat Puas',       warna:'var(--s5)', ink:'var(--s5i)'},
];

/** Tingkat skala lokal dengan label yang sama; warna dipilih berdasarkan label, bukan angka. */
const sesuaiLabel = (label) => SKALA.find((s) => s.label === label) || SKALA[2];

/** Skala dari API ([{v,label}], 4 atau 5 tingkat) dilengkapi warna lokal. */
export const skalaDari = (api) =>
  (api || []).map((s) => ({ ...sesuaiLabel(s.label), v: s.v, label: s.label, pendek: s.label }));

/** Distribusi dari API ([{v,label,jumlah}]) dilengkapi warna untuk DonutKepuasan. */
export const distribusiBerwarna = (data) =>
  (data || []).map((d) => ({ ...sesuaiLabel(d.label), ...d }));

export const warnaKategori = (label) => sesuaiLabel(label).warna;
