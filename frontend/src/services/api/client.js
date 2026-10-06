/**
 * Klien HTTP tunggal untuk API Laravel. Autentikasi memakai Sanctum mode SPA:
 * sesi disimpan di cookie httpOnly, jadi tidak ada token yang disimpan di
 * browser. Untuk permintaan yang mengubah data, header X-XSRF-TOKEN diambil
 * dari cookie XSRF-TOKEN yang dipasang oleh /sanctum/csrf-cookie.
 */
const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/+$/, '');
const CSRF_URL = `${new URL(BASE || window.location.origin, window.location.origin).origin}/sanctum/csrf-cookie`;

/** Galat API yang sudah dinormalisasi. status 0 = gagal terhubung. */
export class ApiError extends Error {
  constructor(status, message, errors = {}) {
    super(message);
    this.status = status;
    this.errors = errors || {};
  }
  /** Pesan pertama untuk sebuah field, atau undefined. */
  field(nama) {
    const v = this.errors[nama];
    return Array.isArray(v) ? v[0] : v;
  }
}

const PESAN_JARINGAN = 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda lalu coba lagi.';

function bacaCookie(nama) {
  const m = document.cookie.match(new RegExp('(?:^|; )' + nama + '=([^;]*)'));
  return m ? decodeURIComponent(m[1]) : null;
}

export async function siapkanCsrf() {
  try {
    await fetch(CSRF_URL, { credentials: 'include', headers: { Accept: 'application/json' } });
  } catch {
    throw new ApiError(0, PESAN_JARINGAN);
  }
}

function urlDengan(path, params) {
  const u = new URL(BASE + path, window.location.origin);
  Object.entries(params || {}).forEach(([k, v]) => {
    if (v !== undefined && v !== null && v !== '') u.searchParams.set(k, v);
  });
  return u.toString();
}

async function kirim(method, path, { params, body } = {}, sudahUlang = false) {
  const ubah = method !== 'GET';
  if (ubah && !bacaCookie('XSRF-TOKEN')) await siapkanCsrf();

  const headers = { Accept: 'application/json', 'X-Requested-With': 'XMLHttpRequest' };
  if (body !== undefined) headers['Content-Type'] = 'application/json';
  if (ubah) headers['X-XSRF-TOKEN'] = bacaCookie('XSRF-TOKEN') || '';

  let res;
  try {
    res = await fetch(urlDengan(path, params), {
      method, headers, credentials: 'include',
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  } catch {
    throw new ApiError(0, PESAN_JARINGAN);
  }

  // Token CSRF kedaluwarsa: ambil ulang sekali lalu ulangi permintaan.
  if (res.status === 419 && !sudahUlang) {
    await siapkanCsrf();
    return kirim(method, path, { params, body }, true);
  }

  return res;
}

async function json(res) {
  let isi = null;
  try { isi = await res.json(); } catch { /* respons tanpa body JSON */ }
  if (!res.ok) {
    if (res.status === 401) window.dispatchEvent(new Event('skm:tidak-masuk'));
    throw new ApiError(res.status, isi?.message || `Permintaan gagal (${res.status}).`, isi?.errors);
  }
  return isi; // { success, message, data, meta? }
}

export const api = {
  get: (path, params) => kirim('GET', path, { params }).then(json),
  post: (path, body) => kirim('POST', path, { body }).then(json),
  put: (path, body) => kirim('PUT', path, { body }).then(json),
  patch: (path, body) => kirim('PATCH', path, { body }).then(json),
  delete: (path) => kirim('DELETE', path).then(json),

  /** Mengunduh berkas (ekspor laporan) lalu memicu dialog simpan di browser. */
  async unduh(path, params) {
    const res = await kirim('GET', path, { params });
    if (!res.ok) await json(res);
    const blob = await res.blob();
    const cd = res.headers.get('Content-Disposition') || '';
    const nama = (cd.match(/filename\*?=(?:UTF-8'')?"?([^";]+)"?/i) || [])[1] || 'laporan';
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = decodeURIComponent(nama);
    document.body.appendChild(a);
    a.click();
    a.remove();
    setTimeout(() => URL.revokeObjectURL(a.href), 1000);
  },
};
