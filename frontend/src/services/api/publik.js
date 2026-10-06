import { api } from './client.js';

/**
 * Endpoint publik. Data yang jarang berubah (referensi, form, ringkasan,
 * periode) di-cache per muat halaman agar beberapa bagian landing yang
 * memakainya tidak memanggil API berulang kali.
 */
const cache = new Map();
function sekali(kunci, fn) {
  if (!cache.has(kunci)) {
    cache.set(kunci, fn().catch((e) => { cache.delete(kunci); throw e; }));
  }
  return cache.get(kunci);
}

/** Dipanggil setelah admin mengubah pengaturan atau pertanyaan agar data publik dimuat ulang. */
export const hapusCache = () => cache.clear();

export const referensi =() => sekali('referensi', () => api.get('/referensi')).then((r) => r.data);
export const formSurvei = () => sekali('form', () => api.get('/survei/form')).then((r) => r.data);
export const ringkasan = () => sekali('ringkasan', () => api.get('/ringkasan')).then((r) => r.data);
export const daftarPeriode = () => sekali('periode', () => api.get('/periode')).then((r) => r.data);

export const hasil = (filter) => api.get('/hasil', filter).then((r) => r.data);

/** @returns {Promise<{nomor:string, rata:number, skala_maks:number, indeks:number, mutu:{huruf,label}}>} */
export const kirimSurvei = (payload) => api.post('/survei', payload).then((r) => r.data);

export const kirimKontak = (payload) => api.post('/kontak', payload);
