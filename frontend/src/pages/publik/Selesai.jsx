import { Navigate, useLocation, useNavigate } from 'react-router-dom';
import { Badge, Button } from '../../components/ui/index.js';
import Icon from '../../components/Icon.jsx';
import { fmt } from '../../utils/format.js';
import { mutu } from '../../utils/mutu.js';
import Shell from '../../layouts/Shell.jsx';

function Selesai() {
  const navigate = useNavigate();
  const { state } = useLocation();
  if (!state) return <Navigate to="/survei" replace />;
  const hasil = state;
  const m = mutu(hasil.indeks);
  return <Shell>
    <section className="py-16 sm:py-24">
      <div className="mx-auto text-center anim-rise" style={{maxWidth:'640px'}}>
        <div className="mx-auto w-16 h-16 grid place-items-center rounded-full mb-6" style={{background:'color-mix(in srgb, var(--ok) 14%, transparent)', color:'var(--ok)'}}>
          <Icon name="checkCircle" className="w-8 h-8" sw={2}/>
        </div>
        <h1 className="font-display text-[32px] sm:text-[40px] font-extrabold tracking-[-0.02em]">Terima kasih!</h1>
        <p className="mt-3 text-[17px] text-ink2">Pendapat Anda telah berhasil dikirim.</p>
        <div className="mt-8 bg-surface border border-line rounded-card p-6 text-left">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div>
              <p className="eyebrow text-[10.5px] font-bold text-muted mb-1.5">Nomor responden</p>
              <p className="font-num text-[19px] font-semibold select-all">{hasil.nomor}</p>
            </div>
            <div className="text-right">
              <p className="eyebrow text-[10.5px] font-bold text-muted mb-1.5">Nilai rata-rata Anda</p>
              <p className="font-num text-[19px] font-semibold">{fmt(hasil.rata)} <span className="text-muted text-[14px]">/ {fmt(hasil.skala_maks)}</span></p>
            </div>
            <Badge tone={m.tone}>Setara indeks {fmt(hasil.indeks)} · {m.label}</Badge>
          </div>
          <p className="mt-5 pt-5 border-t border-linesoft text-[14px] text-ink2 leading-relaxed">
            Simpan nomor responden bila Anda ingin menanyakan tindak lanjut masukan. Masukan Anda membantu BRIDA Kabupaten Buleleng dalam meningkatkan kualitas pelayanan publik.
          </p>
        </div>
        <div className="mt-7 flex flex-col sm:flex-row gap-3 justify-center">
          <Button size="lg" variant="outline" icon="home" onClick={()=>navigate('/')}>Kembali ke Beranda</Button>
          <Button size="lg" icon="barChart" onClick={()=>navigate('/#hasil')}>Lihat Hasil Survei</Button>
        </div>
      </div>
    </section>
  </Shell>;
}

export default Selesai;
