/**
 * Klien HTTP tunggal untuk API Go Singa Riset (Fiber). Autentikasi admin
 * memakai JWT dari /v1/auth/login yang dikirim sebagai header
 * `Authorization: Bearer`. Token disimpan di sessionStorage (atau
 * localStorage bila "ingat saya" dicentang) karena backend tidak memakai
 * cookie sesi.
 *
 * Bentuk respons backend: { success, status, ResponseCode, message, data, error, meta? }.
 * Galat validasi dikirim sebagai 400 dengan `error` berbentuk "field: pesan";
 * klien menormalkannya menjadi ApiError berstatus 422 dengan `errors[field]`
 * agar halaman cukup memeriksa `ex.status === 422` dan `ex.field(nama)`.
 */
const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/+$/, '');
const KUNCI_TOKEN = 'skm_token';

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

function simpanan() {
  try { return [window.sessionStorage, window.localStorage]; } catch { return []; }
}

export const token = {
  ambil() {
    for (const s of simpanan()) {
      try { const t = s.getItem(KUNCI_TOKEN); if (t) return t; } catch { /* storage diblokir */ }
    }
    return null;
  },
  simpan(nilai, ingat) {
    token.hapus();
    const [sesi, lokal] = simpanan();
    try { (ingat ? lokal : sesi)?.setItem(KUNCI_TOKEN, nilai); } catch { /* storage diblokir */ }
  },
  hapus() {
    for (const s of simpanan()) {
      try { s.removeItem(KUNCI_TOKEN); } catch { /* storage diblokir */ }
    }
  },
};

/** Isi klaim JWT (tanpa verifikasi; verifikasi tetap di backend). */
export function klaimToken() {
  const t = token.ambil();
  if (!t) return null;
  try {
    const isi = t.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(decodeURIComponent(escape(atob(isi))));
  } catch {
    return null;
  }
}

function urlDengan(path, params) {
  const u = new URL(BASE + path, window.location.origin);
  Object.entries(params || {}).forEach(([k, v]) => {
    if (v !== undefined && v !== null && v !== '') u.searchParams.set(k, v);
  });
  return u.toString();
}

async function kirim(method, path, { params, body } = {}) {
  const headers = { Accept: 'application/json' };
  if (body !== undefined) headers['Content-Type'] = 'application/json';
  const t = token.ambil();
  if (t) headers.Authorization = `Bearer ${t}`;

  try {
    return await fetch(urlDengan(path, params), {
      method, headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  } catch {
    throw new ApiError(0, PESAN_JARINGAN);
  }
}

const POLA_FIELD = /^([a-z_][\w.]*): (.+)$/s;

async function json(res) {
  let isi = null;
  try { isi = await res.json(); } catch { /* respons tanpa body JSON */ }
  if (!res.ok) {
    if (res.status === 401 && token.ambil()) {
      token.hapus();
      window.dispatchEvent(new Event('skm:tidak-masuk'));
    }
    const detail = typeof isi?.error === 'string' ? isi.error : '';
    const cocok = detail.match(POLA_FIELD);
    if (res.status === 400 && cocok) {
      throw new ApiError(422, cocok[2], { [cocok[1]]: [cocok[2]] });
    }
    // Pesan dari kode SKM sendiri (403/409/400) ada di `error`; selain itu pakai `message`.
    const pesan = [400, 403, 409].includes(res.status) && detail && detail !== 'internal server error'
      ? detail : isi?.message;
    throw new ApiError(res.status, pesan || `Permintaan gagal (${res.status}).`);
  }
  // `data` dihilangkan backend bila kosong (omitempty).
  return { ...isi, data: isi?.data ?? null }; // { success, message, data, meta? }
}

export const api = {
  get: (path, params) => kirim('GET', path, { params }).then(json),
  post: (path, body) => kirim('POST', path, { body }).then(json),
  put: (path, body) => kirim('PUT', path, { body }).then(json),
  patch: (path, body) => kirim('PATCH', path, { body }).then(json),
  delete: (path) => kirim('DELETE', path).then(json),

  /**
   * Mengunduh berkas (ekspor laporan) lalu memicu dialog simpan di browser.
   * Content-Disposition tidak selalu terbaca lintas origin, jadi sediakan `namaCadangan`.
   */
  async unduh(path, params, namaCadangan = 'laporan') {
    const res = await kirim('GET', path, { params });
    if (!res.ok) await json(res);
    const blob = await res.blob();
    const cd = res.headers.get('Content-Disposition') || '';
    const nama = (cd.match(/filename\*?=(?:UTF-8'')?"?([^";]+)"?/i) || [])[1] || namaCadangan;
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = decodeURIComponent(nama);
    document.body.appendChild(a);
    a.click();
    a.remove();
    setTimeout(() => URL.revokeObjectURL(a.href), 1000);
  },
};
