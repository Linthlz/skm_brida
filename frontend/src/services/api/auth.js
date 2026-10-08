import { api, ApiError, klaimToken, token } from './client.js';

/**
 * Login memakai akun Singa Riset (POST /v1/auth/login). Panel SKM hanya
 * untuk role "admin"; aturan yang sama ditegakkan backend di setiap
 * endpoint /v1/admin/skm.
 */

/** Bentuk user backend (UserResponse) -> bentuk yang dipakai komponen. */
function keUser(u) {
  const nama = (u?.name || '').trim();
  const inisial = nama.split(/\s+/).filter(Boolean).slice(0, 2).map((k) => k[0]).join('').toUpperCase();
  return { id: u?.public_id, nama, email: u?.email, jabatan: u?.position, role: u?.role, inisial };
}

const tolakBukanAdmin = () => {
  token.hapus();
  return new ApiError(422, 'Akun ini tidak memiliki akses ke panel SKM.', { email: ['Akun ini tidak memiliki akses ke panel SKM.'] });
};

export async function login({ email, password, ingat }) {
  let r;
  try {
    r = await api.post('/v1/auth/login', { email, password });
  } catch (ex) {
    if (ex.status === 401) {
      const pesan = 'Email atau kata sandi tidak sesuai.';
      throw new ApiError(422, pesan, { email: [pesan] });
    }
    throw ex;
  }
  const user = keUser(r.data?.Data);
  if (user.role !== 'admin') throw tolakBukanAdmin();
  token.simpan(r.data.Token, ingat);
  return user;
}

/** JWT tidak punya sesi di server; keluar cukup menghapus token. */
export async function logout() {
  token.hapus();
}

/** Profil admin yang sedang masuk (GET /v1/users/:public_id). */
export async function saya() {
  const klaim = klaimToken();
  if (!klaim?.pub_id || (klaim.exp && klaim.exp * 1000 < Date.now())) {
    token.hapus();
    throw new ApiError(401, 'Silakan masuk terlebih dahulu.');
  }
  const user = keUser((await api.get(`/v1/users/${klaim.pub_id}`)).data);
  if (user.role !== 'admin') {
    token.hapus();
    throw new ApiError(401, 'Akun ini tidak memiliki akses ke panel SKM.');
  }
  return user;
}
