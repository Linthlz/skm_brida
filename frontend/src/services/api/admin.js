import { api } from './client.js';

/** Endpoint admin SKM di backend Go Singa Riset: /v1/admin/skm/* (role admin). */
const P = '/v1/admin/skm';

const data = (r) => r.data;
const daftar = (r) => r.data ?? [];

/**
 * Daftar berpaginasi. Backend memakai `page` + `limit` dan meta
 * { page, limit, total, total_page }; komponen (DataTable, Paginasi) memakai
 * { current_page, last_page, per_page, total }, jadi diterjemahkan di sini.
 */
async function halaman(path, { per_page, ...params } = {}) {
  const r = await api.get(path, { ...params, limit: per_page ?? params.limit });
  const m = r.meta || {};
  return {
    ...r,
    data: r.data ?? [],
    meta: { current_page: m.page ?? 1, last_page: Math.max(1, m.total_page ?? 1), per_page: m.limit, total: m.total ?? 0 },
  };
}

// Ringkasan & agregat
export const dashboard = (filter) => api.get(`${P}/dashboard`, filter).then(data);
export const hasil = (filter) => api.get(`${P}/hasil`, filter).then(data);
export const statistik = (filter) => api.get(`${P}/statistik`, filter).then(data);

// Responden (berpaginasi: kembalikan { data, meta })
export const daftarResponden = (params) => halaman(`${P}/responden`, params);
export const detailResponden = (id) => api.get(`${P}/responden/${id}`).then(data);

// Bank pertanyaan
export const daftarPertanyaan = () => api.get(`${P}/pertanyaan`).then(daftar);
export const tambahPertanyaan = (body) => api.post(`${P}/pertanyaan`, body).then(data);
export const ubahPertanyaan = (id, body) => api.patch(`${P}/pertanyaan/${id}`, body).then(data);
export const hapusPertanyaan = (id) => api.delete(`${P}/pertanyaan/${id}`);
export const urutkanPertanyaan = (ids) => api.put(`${P}/pertanyaan/urutan`, { ids }).then(daftar);

// Feedback & pesan kontak (berpaginasi). Hitungan per status / belum dibaca
// diambil dari endpoint terpisah lalu digabung ke meta seperti sebelumnya.
export async function daftarFeedback(params) {
  const [r, hitung] = await Promise.all([halaman(`${P}/feedback`, params), api.get(`${P}/feedback/hitung`).then(data)]);
  return { ...r, meta: { ...r.meta, hitung } };
}
export const ubahFeedback = (id, body) => api.patch(`${P}/feedback/${id}`, body).then(data);
export async function daftarPesan(params) {
  const [r, hitung] = await Promise.all([halaman(`${P}/pesan-kontak`, params), api.get(`${P}/pesan-kontak/hitung`).then(data)]);
  return { ...r, meta: { ...r.meta, belum_dibaca: hitung?.belum_dibaca ?? 0 } };
}
export const ubahPesan = (id, body) => api.patch(`${P}/pesan-kontak/${id}`, body).then(data);

// Laporan & pengaturan
export const unduhLaporan = (params) =>
  api.unduh(`${P}/laporan/export`, params, `laporan-skm${params?.periode ? '-' + params.periode : ''}.${params?.format || 'pdf'}`);
export const pengaturan = () => api.get(`${P}/pengaturan`).then(data);
export const simpanPengaturan = (body) => api.put(`${P}/pengaturan`, body).then(data);
