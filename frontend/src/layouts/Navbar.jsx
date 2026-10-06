import { useNavigate } from 'react-router-dom';
import Logo from '../components/Logo.jsx';

/** Bagian-bagian landing page, dipakai sebagai tautan jangkar di footer. */
export const BAGIAN_LANDING = [
  { to: '/#tentang', label: 'Tentang Survei' },
  { to: '/#hasil', label: 'Hasil Survei' },
  { to: '/#kontak', label: 'Kontak' },
];

function Navbar() {
  const navigate = useNavigate();

  return (
    <header
      className="sticky z-50 text-[var(--nav-ink)]"
      style={{ top: 'env(safe-area-inset-top, 0px)', background: 'var(--nav)' }}
    >
      <div className="h-1 bg-[var(--gold-bright)]" />
      <div className="mx-auto max-w-[1180px] px-4 sm:px-6">
        <div className="flex items-center gap-4 h-[70px]">
          <button onClick={() => navigate('/')} className="flex items-center gap-3 text-left shrink-0">
            <Logo tone="nav" className="w-9 h-9 sm:w-10 sm:h-10" />
            <span className="leading-tight">
              <span className="block font-display font-extrabold text-[15px] sm:text-[16px] tracking-[-0.01em]">
                BRIDA Kabupaten Buleleng
              </span>
              <span className="block text-[11px] sm:text-[12px] text-[var(--nav-ink-2)]">
                Survei Kepuasan Masyarakat
              </span>
            </span>
          </button>
        </div>
      </div>
    </header>
  );
}

export default Navbar;
