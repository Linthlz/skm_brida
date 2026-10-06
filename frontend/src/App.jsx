import { useEffect } from 'react';
import { Navigate, Outlet, Route, Routes, useLocation } from 'react-router-dom';
import { ToastProvider } from './hooks/useToast.jsx';
import { AuthProvider, RequireAuth, useAuth } from './hooks/useAuth.jsx';
import { Card, GagalMuat, Memuat } from './components/ui/index.js';

import PublicLayout from './layouts/PublicLayout.jsx';
import AdminLayout from './layouts/AdminLayout.jsx';

import Landing from './pages/publik/Landing.jsx';
import Survei from './pages/publik/Survei.jsx';
import Selesai from './pages/publik/Selesai.jsx';

import Login from './pages/admin/Login.jsx';
import Dashboard from './pages/admin/Dashboard.jsx';
import Responden from './pages/admin/Responden.jsx';
import DataSurvei from './pages/admin/DataSurvei.jsx';
import Pertanyaan from './pages/admin/Pertanyaan.jsx';
import HasilSurvei from './pages/admin/HasilSurvei.jsx';
import Statistik from './pages/admin/Statistik.jsx';
import Feedback from './pages/admin/Feedback.jsx';
import Laporan from './pages/admin/Laporan.jsx';
import Pengaturan from './pages/admin/Pengaturan.jsx';

/**
 * Mulai dari atas setiap kali pindah halaman. Kalau alamatnya membawa jangkar
 * (#tentang, #hasil, #kontak), gulir ke bagian itu — dibutuhkan karena navigasi
 * router tidak memicu lompatan jangkar bawaan browser.
 */
function ScrollToTop() {
  const { pathname, hash } = useLocation();
  useEffect(() => {
    const sasaran = hash ? document.getElementById(hash.slice(1)) : null;
    if (sasaran) sasaran.scrollIntoView();
    else window.scrollTo({ top: 0 });
  }, [pathname, hash]);
  return null;
}

/** Tampilan selama sesi admin diperiksa, atau bila server tidak dapat dihubungi. */
function LayarSesi() {
  const { status, periksa } = useAuth();
  return (
    <div className="min-h-[100svh] grid place-items-center p-4 bg-ground">
      <Card className="w-full" pad="p-6">
        <div style={{ maxWidth: '420px' }} className="mx-auto">
          {status === 'galat'
            ? <GagalMuat error={{ status: 0, message: 'Sesi admin tidak dapat diperiksa. Pastikan server API berjalan.' }} onRetry={periksa} />
            : <Memuat baris={4} />}
        </div>
      </Card>
    </div>
  );
}

export default function App() {
  return (
    <ToastProvider>
      <ScrollToTop />
      <Routes>
        <Route element={<PublicLayout />}>
          <Route index element={<Landing />} />
          <Route path="survei" element={<Survei />} />
          <Route path="survei/selesai" element={<Selesai />} />
        </Route>

        <Route path="admin" element={<AuthProvider><Outlet /></AuthProvider>}>
          <Route path="login" element={<Login />} />
          <Route element={<RequireAuth fallback={<LayarSesi />}><AdminLayout /></RequireAuth>}>
            <Route index element={<Dashboard />} />
            <Route path="responden" element={<Responden />} />
            <Route path="survei" element={<DataSurvei />} />
            <Route path="pertanyaan" element={<Pertanyaan />} />
            <Route path="hasil" element={<HasilSurvei />} />
            <Route path="statistik" element={<Statistik />} />
            <Route path="feedback" element={<Feedback />} />
            <Route path="laporan" element={<Laporan />} />
            <Route path="pengaturan" element={<Pengaturan />} />
          </Route>
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </ToastProvider>
  );
}
