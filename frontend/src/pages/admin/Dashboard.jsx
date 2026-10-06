import { StatCard, Card, EmptyState, GagalMuat, Memuat } from '../../components/ui/index.js';
import { ChartCard, BarUnsur, DonutKepuasan, KolomResponden, MeterIKM } from '../../components/charts/index.js';
import { useApi } from '../../hooks/useApi.js';
import { dashboard } from '../../services/api/admin.js';
import { distribusiBerwarna } from '../../data/skala.js';
import { fmt, fmtInt, unsurGrafik } from '../../utils/format.js';

function Dashboard() {
  const { data:d, error, reload } = useApi(dashboard);

  if (error) return <Card><GagalMuat error={error} onRetry={reload}/></Card>;
  if (!d) return <div className="space-y-4">
    <div className="grid sm:grid-cols-2 xl:grid-cols-4 gap-4">{[0,1,2,3].map(i=><Memuat key={i} tinggi={112}/>)}</div>
    <div className="grid xl:grid-cols-[1.25fr_.75fr] gap-4"><Memuat tinggi={300}/><Memuat tinggi={300}/></div>
  </div>;

  const kosong = d.total_responden === 0;
  return <div className="space-y-4">
    <div className="grid sm:grid-cols-2 xl:grid-cols-4 gap-4">
      <StatCard label="Total responden" value={fmtInt(d.total_responden)} sub={`+${fmtInt(d.bulan_ini)} bulan ini`} icon="users"/>
      <StatCard label="Survei terkumpul" value={fmtInt(d.total_responden)} sub={d.persen_lengkap!==null ? `${d.persen_lengkap}% lengkap` : undefined} icon="clipboard"/>
      <StatCard label="Indeks kepuasan" value={fmt(d.ikm)} sub={d.mutu && `Mutu ${d.mutu.huruf}`} icon="gauge" tone="gold"/>
      <StatCard label="Feedback masuk" value={fmtInt(d.feedback_total)} sub={`${fmtInt(d.feedback_baru)} belum ditinjau`} icon="message"/>
    </div>

    {kosong
      ? <Card><EmptyState icon="clipboard" title={`Belum ada responden pada periode ${d.periode}`} desc="Grafik akan muncul setelah survei pertama masuk."/></Card>
      : <>
      <div className="grid xl:grid-cols-[1.25fr_.75fr] gap-4 items-start">
        <ChartCard title="Responden per bulan" desc={`Jumlah survei yang masuk sepanjang periode ${d.periode}.`}>
          <KolomResponden labels={d.per_bulan.map(b=>b.label)} values={d.per_bulan.map(b=>b.jumlah)} tahun={d.periode}/>
        </ChartCard>
        <ChartCard title="Distribusi kepuasan" desc="Sebaran kategori penilaian responden.">
          <DonutKepuasan data={distribusiBerwarna(d.distribusi)}/>
        </ChartCard>
      </div>
      <div className="grid xl:grid-cols-[1.25fr_.75fr] gap-4 items-start">
        <ChartCard title="Nilai kepuasan per unsur" desc={`Skala 1–${d.skala_maks}. Tiga unsur terendah otomatis masuk daftar prioritas.`}>
          <BarUnsur data={unsurGrafik(d.unsur)} max={d.skala_maks} pengali={d.pengali}/>
        </ChartCard>
        <div className="space-y-4 min-w-0">
          <Card>
            <p className="eyebrow text-[10.5px] font-bold text-muted mb-4">Prioritas perbaikan</p>
            <ul className="space-y-3">
              {d.prioritas.map((u,i)=>(
                <li key={u.kode} className="flex items-center gap-3">
                  <span className="font-num text-[12px] font-semibold w-6 h-6 grid place-items-center rounded-[6px] bg-[color:var(--warn)]/14 text-[var(--warn)] shrink-0">{i+1}</span>
                  <span className="text-[14px] flex-1 min-w-0">{u.nama}</span>
                  <span className="font-num text-[13.5px] font-semibold">{fmt(u.nrr)}</span>
                </li>
              ))}
            </ul>
          </Card>
          {d.ikm !== null && <Card><MeterIKM nilai={d.ikm}/></Card>}
        </div>
      </div>
    </>}
  </div>;
}

export default Dashboard;
