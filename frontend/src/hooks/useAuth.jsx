import { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import * as authApi from '../services/api/auth.js';

/**
 * Status sesi admin. Sumber kebenaran tetap backend: frontend hanya
 * menanyakan /auth/me dan mengikuti hasilnya. Setiap respons 401 dari API
 * memicu event `skm:tidak-masuk` sehingga pengguna dikembalikan ke login.
 */
const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [status, setStatus] = useState('memuat'); // memuat | masuk | tamu | galat

  const periksa = useCallback(() => {
    setStatus('memuat');
    authApi.saya().then(
      (u) => { setUser(u); setStatus('masuk'); },
      (e) => { setUser(null); setStatus(e.status === 401 ? 'tamu' : 'galat'); },
    );
  }, []);

  useEffect(() => { periksa(); }, [periksa]);

  useEffect(() => {
    const keluarPaksa = () => { setUser(null); setStatus('tamu'); };
    window.addEventListener('skm:tidak-masuk', keluarPaksa);
    return () => window.removeEventListener('skm:tidak-masuk', keluarPaksa);
  }, []);

  const login = useCallback(async (kredensial) => {
    const u = await authApi.login(kredensial);
    setUser(u);
    setStatus('masuk');
    return u;
  }, []);

  const logout = useCallback(async () => {
    try { await authApi.logout(); } finally { setUser(null); setStatus('tamu'); }
  }, []);

  return (
    <AuthContext.Provider value={{ user, status, login, logout, periksa }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);

/** Penjaga rute admin: tamu dialihkan ke /admin/login (alamat asal diingat). */
export function RequireAuth({ children, fallback }) {
  const { status } = useAuth();
  const loc = useLocation();
  if (status === 'masuk') return children;
  if (status === 'tamu') return <Navigate to="/admin/login" replace state={{ dari: loc.pathname + loc.search }} />;
  return fallback;
}
