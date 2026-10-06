import { Card, GagalMuat, Memuat } from '../../components/ui/index.js';
import { ChartCard, TrenIKM } from '../../components/charts/index.js';
import { useApi } from '../../hooks/useApi.js';
import { statistik } from '../../services/api/admin.js';
import { fmtInt } from '../../utils/format.js';

const Baris = ({label,n,max})=>(
  <li className="flex items-center gap-3 py-2">
    <span className="text-[13.5px] w-[150px] sm:w-[190px] shrink-0">{label}</span>
    <span className="flex-1 h-2.5 rounded-full bg-[var(--ground)] overflow-hidden"><span className="block h-full rounded-full bg-[var(--brand)]" style={{width:`${max?(n/max)*100:0}%`}}/></span>
    <span className="font-num text-[13px] w-12 text-right">{fmtInt(n)}</span>
  </li>
);

function Statistik() {
  const { data:d, error, reload } = useApi(statistik);
  if (error) return <Card><GagalMuat error={error} onRetry={reload}/></Card>;
  if (!d) return <div className="space-y-4"><div className="grid lg:grid-cols-2 gap-4"><Memuat tinggi={300}/><Memuat tinggi={300}/></div><Memuat tinggi={280}/></div>;

  const maxA = Math.max(0, ...d.pendidikan.map(x=>x.n)), maxB = Math.max(0, ...d.pekerjaan.map(x=>x.n));
  const ket = `Jumlah responden periode ${d.periode} (${fmtInt(d.total_responden)} orang).`;
  return <div className="space-y-4">
    <div className="grid lg:grid-cols-2 gap-4">
      <ChartCard title="Responden menurut pendidikan terakhir" desc={ket}>
        <ul>{d.pendidikan.map(x=><Baris key={x.label} {...x} max={maxA}/>)}</ul>
      </ChartCard>
      <ChartCard title="Responden menurut pekerjaan" desc="Kelompok pengguna layanan terbanyak menjadi acuan penyusunan standar layanan.">
        <ul>{d.pekerjaan.map(x=><Baris key={x.label} {...x} max={maxB}/>)}</ul>
      </ChartCard>
    </div>
    {d.tren.length>0 && <ChartCard title="Tren Indeks Kepuasan Masyarakat" desc={`${d.tren.length} periode survei yang memiliki data.`}>
      <TrenIKM data={d.tren}/>
    </ChartCard>}
  </div>;
}

export default Statistik;
