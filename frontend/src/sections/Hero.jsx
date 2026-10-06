import { useNavigate } from 'react-router-dom';
import { Badge, Button, GagalMuat, Memuat } from '../components/ui/index.js';
import { MeterIKM } from '../components/charts/index.js';
import Icon from '../components/Icon.jsx';
import { useApi } from '../hooks/useApi.js';
import { formSurvei, ringkasan } from '../services/api/publik.js';
import { fmtInt } from '../utils/format.js';
import Shell from '../layouts/Shell.jsx';

function Hero() {
  const navigate = useNavigate();
  const form = useApi(formSurvei);
  const ring = useApi(ringkasan);
  const per = form.data?.periode;
  const r = ring.data;

  return (
    <section className="bg-surface border-b border-line">
      <Shell>
        <div className="grid lg:grid-cols-[1.05fr_.95fr] gap-10 lg:gap-14 items-center py-12 sm:py-16">
          <div className="anim-rise">
            {per && <Badge tone={per.ditutup ? 'neutral' : 'gold'} className="mb-5">
              Periode survei {per.tahun} · {per.ditutup ? 'sudah ditutup' : 'sedang berlangsung'}
            </Badge>}
            <h1 className="font-display text-[34px] sm:text-[46px] font-extrabold leading-[1.08] tracking-[-0.02em]">
              Survei Kepuasan Masyarakat
            </h1>
            <p className="mt-4 text-[17px] text-ink2 leading-relaxed" style={{maxWidth:'52ch'}}>
              Bantu kami meningkatkan kualitas pelayanan BRIDA Kabupaten Buleleng melalui pendapat dan pengalaman Anda.
            </p>
            <div className="mt-7 flex flex-col sm:flex-row gap-3">
              <Button size="lg" icon="clipboard" onClick={()=>navigate('/survei')}>Mulai Survei</Button>
              <Button as="a" href="#tentang" size="lg" variant="outline" iconRight="arrowRight">Pelajari Survei</Button>
            </div>
            <div className="mt-7 flex flex-wrap items-center gap-x-6 gap-y-2 text-[13.5px] text-muted">
              <span className="flex items-center gap-2"><Icon name="clock" className="w-4 h-4 text-[var(--brand)]"/>± 3–5 menit</span>
              {form.data?.anonim !== false && <span className="flex items-center gap-2"><Icon name="shield" className="w-4 h-4 text-[var(--brand)]"/>Identitas boleh dikosongkan</span>}
              {form.data && <span className="flex items-center gap-2"><Icon name="listChecks" className="w-4 h-4 text-[var(--brand)]"/>{form.data.pertanyaan.length} unsur pelayanan</span>}
            </div>
          </div>

          <div className="anim-rise" style={{animationDelay:'.08s'}}>
            <div className="bg-[var(--ground)] border border-line rounded-card p-6">
              {ring.error ? <GagalMuat error={ring.error} onRetry={ring.reload}/>
              : !r ? <Memuat tinggi={210} className="!bg-surface"/>
              : <>
                <div className="flex items-center justify-between gap-3 mb-5">
                  <p className="eyebrow text-[10.5px] font-bold text-muted">Hasil survei periode berjalan</p>
                  <Badge tone="neutral">{r.periode}</Badge>
                </div>
                {r.ikm !== null
                  ? <MeterIKM nilai={r.ikm}/>
                  : <p className="text-[14px] text-muted py-6">Belum ada penilaian yang masuk pada periode ini.</p>}
                <div className="grid grid-cols-3 gap-3 mt-6 pt-5 border-t border-line">
                  {[['Responden', fmtInt(r.total_responden)], ['Unsur dinilai', fmtInt(r.jumlah_unsur)], ['Unit layanan', fmtInt(r.jumlah_layanan)]].map(([k,v])=>(
                    <div key={k}>
                      <p className="font-num text-[19px] font-semibold leading-none">{v}</p>
                      <p className="text-[12px] text-muted mt-1.5">{k}</p>
                    </div>
                  ))}
                </div>
                <a href="#hasil" className="mt-5 inline-flex items-center gap-1.5 text-[14px] font-semibold text-[var(--brand)] hover:gap-2.5 transition-all">
                  Lihat rincian hasil survei <Icon name="arrowRight" className="w-4 h-4"/>
                </a>
              </>}
            </div>
          </div>
        </div>
      </Shell>
    </section>
  );
}

export default Hero;
