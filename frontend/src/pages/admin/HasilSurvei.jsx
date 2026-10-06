import { Card, EmptyState, GagalMuat, Memuat, StatCard } from '../../components/ui/index.js';
import { ChartCard, BarUnsur, DonutKepuasan, TrenIKM } from '../../components/charts/index.js';
import { useApi } from '../../hooks/useApi.js';
import { hasil } from '../../services/api/admin.js';
import { distribusiBerwarna } from '../../data/skala.js';
import { fmt, fmtInt, unsurGrafik } from '../../utils/format.js';

function HasilSurvei() {
  const { data:d, error, reload } = useApi(hasil);
  if (error) return <Card><GagalMuat error={error} onRetry={reload}/></Card>;
  if (!d) return <div className="space-y-4">
    <div className="grid sm:grid-cols-2 xl:grid-cols-4 gap-4">{[0,1,2,3].map(i=><Memuat key={i} tinggi={112}/>)}</div>
    <div className="grid xl:grid-cols-2 gap-4"><Memuat tinggi={380}/><Memuat tinggi={380}/></div>
  </div>;

  return <div className="space-y-4">
    <div className="grid sm:grid-cols-2 xl:grid-cols-4 gap-4">
      <StatCard label="Indeks Kepuasan Masyarakat" value={fmt(d.ikm)} sub={d.mutu && `Mutu ${d.mutu.huruf}`} icon="gauge"/>
      <StatCard label="Nilai rata-rata unsur" value={fmt(d.nrr)} sub={`dari ${fmt(d.skala_maks)}`} icon="star"/>
      <StatCard label="Responden" value={fmtInt(d.total_responden)} sub={`periode ${d.periode}`} icon="users"/>
      <StatCard label="Kategori mutu" value={d.mutu?.label || '–'} icon="award" tone="gold"/>
    </div>
    {d.unsur.length===0
      ? <Card><EmptyState icon="barChart" title={`Belum ada data pada periode ${d.periode}`} desc="Rekap akan muncul setelah survei pertama masuk."/></Card>
      : <div className="grid xl:grid-cols-2 gap-4 items-start">
        <ChartCard title="Nilai per unsur pelayanan" desc={`Skala 1–${d.skala_maks} dengan konversi indeks pengali ${fmt(d.pengali,0)}.`}>
          <BarUnsur data={unsurGrafik(d.unsur)} max={d.skala_maks} pengali={d.pengali}/>
        </ChartCard>
        <div className="space-y-4 min-w-0">
          <ChartCard title="Distribusi kepuasan"><DonutKepuasan data={distribusiBerwarna(d.distribusi)}/></ChartCard>
          {d.tren.length>0 && <ChartCard title="Tren antarperiode"><TrenIKM data={d.tren}/></ChartCard>}
        </div>
      </div>}
  </div>;
}

export default HasilSurvei;
