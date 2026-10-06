import { useEffect, useState } from 'react';
import { NavLink, Outlet, useLocation, useNavigate } from 'react-router-dom';
import Icon from '../components/Icon.jsx';
import Logo from '../components/Logo.jsx';
import { Button } from '../components/ui/index.js';
import { useAuth } from '../hooks/useAuth.jsx';
import { useToast } from '../hooks/useToast.jsx';

export const NAV_ADMIN = [
  { to: '/admin', end: true, label: 'Dashboard', icon: 'grid', judul: 'Dashboard', desc: 'Ringkasan hasil Survei Kepuasan Masyarakat periode aktif' },
  { to: '/admin/responden', label: 'Data Responden', icon: 'users', judul: 'Data Responden', desc: 'Daftar pengisi survei beserta karakteristiknya' },
  { to: '/admin/survei', label: 'Data Survei', icon: 'clipboard', judul: 'Data Survei', desc: 'Seluruh pengisian survei yang masuk' },
  { to: '/admin/pertanyaan', label: 'Pertanyaan', icon: 'listChecks', judul: 'Kelola Pertanyaan', desc: 'Susun, aktifkan, dan urutkan pertanyaan survei' },
  { to: '/admin/hasil', label: 'Hasil Survei', icon: 'barChart', judul: 'Hasil Survei', desc: 'Rekapitulasi nilai per unsur pelayanan' },
  { to: '/admin/statistik', label: 'Statistik', icon: 'pieChart', judul: 'Statistik', desc: 'Karakteristik responden dan tren antarperiode' },
  { to: '/admin/feedback', label: 'Feedback', icon: 'message', judul: 'Feedback Masyarakat', desc: 'Kritik dan saran beserta status tindak lanjutnya' },
  { to: '/admin/laporan', label: 'Laporan', icon: 'fileText', judul: 'Laporan', desc: 'Rekapitulasi dan ekspor hasil survei' },
  { to: '/admin/pengaturan', label: 'Pengaturan', icon: 'sliders', judul: 'Pengaturan', desc: 'Konfigurasi survei dan pengelola' },
];

function Sidebar({ mobile, onClose, onKeluar }) {
  const { user } = useAuth();
  return (
    <div className="h-full flex flex-col" style={{ background: 'var(--brand-deep)' }}>
      <div className="flex items-center gap-3 px-5 h-[70px] shrink-0 border-b border-white/10">
        <Logo tone="nav" className="w-8 h-8" />
        <span className="font-display font-extrabold text-[13.5px] text-white leading-tight">
          BRIDA Buleleng
          <br />
          <span className="text-white/55 text-[11px] font-semibold">Panel SKM</span>
        </span>
        {mobile && (
          <button onClick={onClose} className="ml-auto p-1.5 text-white/70" aria-label="Tutup menu">
            <Icon name="x" className="w-5 h-5" />
          </button>
        )}
      </div>

      <nav className="flex-1 overflow-y-auto p-3 space-y-0.5">
        {NAV_ADMIN.map((n) => (
          <NavLink
            key={n.to}
            to={n.to}
            end={n.end}
            className={({ isActive }) =>
              `w-full flex items-center gap-3 px-3 py-2.5 rounded-[8px] text-[14px] font-semibold transition-colors ${
                isActive ? 'bg-white/12 text-white' : 'text-white/60 hover:text-white hover:bg-white/10'
              }`
            }
          >
            {({ isActive }) => (
              <>
                <Icon name={n.icon} className="w-[18px] h-[18px] shrink-0" />
                {n.label}
                {isActive && <span className="ml-auto w-1.5 h-1.5 rounded-full bg-[var(--gold-bright)]" />}
              </>
            )}
          </NavLink>
        ))}
      </nav>

      <div className="p-3 border-t border-white/10">
        <div className="flex items-center gap-3 px-2 py-2 mb-1">
          <span className="w-8 h-8 rounded-full grid place-items-center bg-[var(--gold-bright)] text-[var(--gold-ink)] font-display font-extrabold text-[13px]">
            {user?.inisial}
          </span>
          <span className="min-w-0">
            <span className="block text-[13px] font-semibold text-white truncate">{user?.nama}</span>
            <span className="block text-[11.5px] text-white/50 truncate">{user?.jabatan || 'Admin'}</span>
          </span>
        </div>
        <button
          onClick={onKeluar}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-[8px] text-[13.5px] font-semibold text-white/60 hover:text-white hover:bg-white/10"
        >
          <Icon name="logOut" className="w-[17px] h-[17px]" /> Keluar
        </button>
      </div>
    </div>
  );
}

function AdminLayout() {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const { logout } = useAuth();
  const toast = useToast();
  useEffect(() => {
    setOpen(false);
    window.scrollTo({ top: 0 });
  }, [pathname]);

  async function keluar() {
    try {
      await logout();
      toast('Anda telah keluar');
    } catch {
      // Sesi lokal tetap dihapus walau server tidak terjangkau.
    }
    navigate('/admin/login', { replace: true });
  }

  const aktif =
    [...NAV_ADMIN].reverse().find((n) => (n.end ? pathname === n.to : pathname.startsWith(n.to))) || NAV_ADMIN[0];

  return (
    <div className="min-h-[100svh] flex bg-ground">
      <aside className="hidden lg:block w-[248px] shrink-0 sticky top-0 h-[100svh]">
        <Sidebar onKeluar={keluar} />
      </aside>

      {open && (
        <div className="lg:hidden fixed inset-0 z-[70] flex">
          <div className="absolute inset-0 bg-[#2A0B0B]/60" onClick={() => setOpen(false)} />
          <div className="relative w-[262px] h-full anim-rise">
            <Sidebar mobile onClose={() => setOpen(false)} onKeluar={keluar} />
          </div>
        </div>
      )}

      <div className="flex-1 min-w-0">
        <header
          className="sticky z-40 bg-surface/95 backdrop-blur border-b border-line"
          style={{ top: 'env(safe-area-inset-top, 0px)' }}
        >
          <div className="flex items-center gap-3 px-4 sm:px-6 h-[70px]">
            <button
              className="lg:hidden p-2 rounded-[8px] border border-line"
              aria-label="Buka menu"
              onClick={() => setOpen(true)}
            >
              <Icon name="menu" className="w-5 h-5" />
            </button>
            <div className="min-w-0">
              <h1 className="font-display font-extrabold text-[18px] sm:text-[20px] truncate">{aktif.judul}</h1>
              <p className="text-[12.5px] text-muted truncate hidden sm:block">{aktif.desc}</p>
            </div>
            <div className="ml-auto flex items-center gap-2">
              <Button
                size="sm"
                variant="outline"
                icon="arrowUpRight"
                onClick={() => navigate('/')}
                className="hidden sm:inline-flex"
              >
                Situs publik
              </Button>
            </div>
          </div>
        </header>

        <main className="p-4 sm:p-6 pb-16">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

export default AdminLayout;
