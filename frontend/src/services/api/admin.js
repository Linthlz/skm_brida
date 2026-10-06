import { api } from './client.js';

const data = (r) => r.data;

// Ringkasan & agregat
export const dashboard = (filter) => api.get('/admin/dashboard', filter).then(data);
export const hasil = (filter) => api.get('/admin/hasil', filter).then(data);
export const statistik = (filter) => api.get('/admin/statistik', filter).then(data);

// Responden (berpaginasi: kembalikan { data, meta })
export const daftarResponden = (params) => api.get('/admin/responden', params);
export const detailResponden = (id) => api.get(`/admin/responden/${id}`).then(data);

// Bank pertanyaan
export const daftarPertanyaan = () => api.get('/admin/pertanyaan').then(data);
export const tambahPertanyaan = (body) => api.post('/admin/pertanyaan', body).then(data);
export const ubahPertanyaan = (id, body) => api.patch(`/admin/pertanyaan/${id}`, body).then(data);
export const hapusPertanyaan = (id) => api.delete(`/admin/pertanyaan/${id}`);
export const urutkanPertanyaan = (ids) => api.put('/admin/pertanyaan/urutan', { ids }).then(data);

// Feedback & pesan kontak (berpaginasi)
export const daftarFeedback = (params) => api.get('/admin/feedback', params);
export const ubahFeedback = (id, body) => api.patch(`/admin/feedback/${id}`, body).then(data);
export const daftarPesan = (params) => api.get('/admin/pesan-kontak', params);
export const ubahPesan = (id, body) => api.patch(`/admin/pesan-kontak/${id}`, body).then(data);

// Laporan & pengaturan
export const unduhLaporan = (params) => api.unduh('/admin/laporan/export', params);
export const pengaturan = () => api.get('/admin/pengaturan').then(data);
export const simpanPengaturan = (body) => api.put('/admin/pengaturan', body).then(data);
