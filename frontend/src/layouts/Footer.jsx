import { Link } from 'react-router-dom';
import Icon from '../components/Icon.jsx';
import Logo from '../components/Logo.jsx';
import { BAGIAN_LANDING } from './Navbar.jsx';

const KONTAK = [
  { icon: 'mapPin', isi: 'Jl. Ngurah Rai No. 2, Singaraja,\nKabupaten Buleleng, Bali 81116' },
  { icon: 'phone', isi: '(0362) 000 000', num: true },
  { icon: 'mail', isi: 'brida@bulelengkab.go.id' },
  { icon: 'clock', isi: 'Senin-Jumat, 08.00-15.30 WITA' },
];

function Footer() {
  return (
    <footer className="mt-20 text-[var(--nav-ink)]" style={{ background: 'var(--nav)' }}>
      <div className="mx-auto max-w-[1180px] px-4 sm:px-6 py-12 grid gap-10 md:grid-cols-[1.4fr_1fr_1fr]">
        <div>
          <div className="flex items-center gap-3 mb-4">
            <Logo tone="nav" className="w-10 h-10" />
            <span className="font-display font-extrabold text-[16px] leading-tight">
              Badan Riset dan Inovasi Daerah
              <br />
              <span className="text-[var(--nav-ink-2)] font-semibold text-[13px]">Kabupaten Buleleng</span>
            </span>
          </div>
          <p className="text-[14px] text-[var(--nav-ink-2)] leading-relaxed" style={{ maxWidth: '44ch' }}>
            Survei Kepuasan Masyarakat diselenggarakan mengacu pada Peraturan Menteri PANRB Nomor 14 Tahun 2017
            tentang Pedoman Penyusunan Survei Kepuasan Masyarakat Unit Penyelenggara Pelayanan Publik.
          </p>
        </div>

        <div>
          <p className="eyebrow text-[10.5px] font-bold text-[var(--nav-ink-2)] mb-3">Navigasi</p>
          <ul className="space-y-2">
            <li>
              <Link to="/survei" className="text-[14px] text-[var(--nav-ink-2)] hover:text-[var(--nav-ink)] hover:underline">
                Isi Survei
              </Link>
            </li>
            {BAGIAN_LANDING.map((n) => (
              <li key={n.to}>
                <Link to={n.to} className="text-[14px] text-[var(--nav-ink-2)] hover:text-[var(--nav-ink)] hover:underline">
                  {n.label}
                </Link>
              </li>
            ))}
          </ul>
        </div>

        <div>
          <p className="eyebrow text-[10.5px] font-bold text-[var(--nav-ink-2)] mb-3">Kontak</p>
          <ul className="space-y-3 text-[14px] text-[var(--nav-ink-2)]">
            {KONTAK.map((k) => (
              <li key={k.icon} className="flex gap-2.5">
                <Icon name={k.icon} className="w-4 h-4 mt-0.5 shrink-0 text-[var(--nav-ink)]" />
                <span className={`whitespace-pre-line select-text ${k.num ? 'font-num' : ''}`}>{k.isi}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>

      <div className="border-t border-[var(--nav-line)]">
        <div className="mx-auto max-w-[1180px] px-4 sm:px-6 py-5 text-[12.5px] text-[var(--nav-ink-2)]">
          <span>
            &copy; {new Date().getFullYear()} BRIDA Kabupaten Buleleng. Survei Kepuasan Masyarakat.
          </span>
        </div>
      </div>
    </footer>
  );
}

export default Footer;
