import { useState } from 'react';
import { Badge, Button, Card, EmptyState, Field, GagalMuat, Memuat, Select } from '../../components/ui/index.js';
import Icon from '../../components/Icon.jsx';
import { useToast } from '../../hooks/useToast.jsx';
import { useApi } from '../../hooks/useApi.js';
import { hasil, unduhLaporan } from '../../services/api/admin.js';
import { daftarPeriode, referensi } from '../../services/api/publik.js';
import { fmt } from '../../utils/format.js';
import { mutu } from '../../utils/mutu.js';

function Laporan() {
  const toast = useToast();
  const [f,setF] = useState({periode:'', layanan_id:'', pendidikan:''});
  const [ekspor,setEkspor] = useState(null);
  const per = useApi(daftarPeriode);
  const ref = useApi(referensi);
  const h = useApi(()=>hasil(f), [f.periode, f.layanan_id, f.pendidikan]);

  const aktif = per.data?.find(p=>p.aktif)?.tahun;
  const layananNama = ref.data?.layanan.find(l=>String(l.id)===f.layanan_id)?.nama;

  async function unduh(format){
    setEkspor(format);
    try{ await unduhLaporan({...f, format}); toast(`Laporan ${format==='pdf'?'PDF':'Excel'} diunduh`); }
    catch(ex){ toast(ex.message,'galat'); }
    finally{ setEkspor(null); }
  }

  const d = h.data;
  const m = d && mutu(d.ikm ?? 0);
  return <div className="space-y-4">
    <Card pad="p-5 sm:p-6">
      <div className="flex items-center gap-2 mb-4"><Icon name="filter" className="w-4 h-4 text-[var(--brand)]"/><h3 className="font-display font-bold text-[16px]">Filter laporan</h3></div>
      <div className="grid sm:grid-cols-3 gap-4">
        <Field label="Periode survei" required><Select id="l-per" placeholder={null} options={(per.data||[]).map(p=>String(p.tahun))} value={f.periode || String(aktif ?? '')} onChange={e=>setF({...f,periode:e.target.value})}/></Field>
        <Field label="Jenis layanan"><Select id="l-lay" options={(ref.data?.layanan||[]).map(l=>({value:String(l.id), label:l.nama}))} placeholder="Semua layanan" value={f.layanan_id} onChange={e=>setF({...f,layanan_id:e.target.value})}/></Field>
        <Field label="Pendidikan responden"><Select id="l-pend" options={ref.data?.pendidikan||[]} placeholder="Semua jenjang" value={f.pendidikan} onChange={e=>setF({...f,pendidikan:e.target.value})}/></Field>
      </div>
      <div className="mt-5 pt-5 border-t border-linesoft flex flex-col sm:flex-row gap-3 sm:items-center">
        <p className="text-[13px] text-muted flex-1">Berkas ekspor memuat rekapitulasi sesuai filter di atas.</p>
        <div className="flex gap-2">
          <Button variant="outline" icon="fileText" disabled={!!ekspor} onClick={()=>unduh('pdf')}>{ekspor==='pdf'?'Menyiapkan…':'Export PDF'}</Button>
          <Button variant="outline" icon="download" disabled={!!ekspor} onClick={()=>unduh('xlsx')}>{ekspor==='xlsx'?'Menyiapkan…':'Export Excel'}</Button>
        </div>
      </div>
    </Card>

    {h.error ? <Card><GagalMuat error={h.error} onRetry={h.reload}/></Card>
    : !d ? <Memuat tinggi={420}/>
    : <Card pad="p-0" className={`transition-opacity ${h.loading?'opacity-60':''}`}>
      <div className="px-5 sm:px-6 py-5 border-b border-linesoft flex flex-wrap items-center gap-3">
        <div>
          <h3 className="font-display text-[17px] font-bold">Rekapitulasi hasil — periode {d.periode}</h3>
          <p className="text-[13px] text-muted mt-0.5">{layananNama || 'Semua layanan'} · {f.pendidikan || 'Semua jenjang pendidikan'} · {d.total_responden} responden</p>
        </div>
        {d.ikm!==null && <Badge tone={m.tone} className="ml-auto">IKM {fmt(d.ikm)} · Mutu {m.huruf}</Badge>}
      </div>
      {d.unsur.length===0 ? <EmptyState icon="fileText" title="Belum ada data pada filter ini" desc="Ubah periode, jenis layanan, atau jenjang pendidikan."/>
      : <div className="scroll-x">
        <table className="w-full text-[14px]" style={{minWidth:'640px'}}>
          <thead><tr className="bg-[var(--ground)] text-left eyebrow text-[11px] font-bold text-muted">
            <th className="px-5 py-3 font-bold">Kode</th><th className="px-5 py-3 font-bold">Unsur pelayanan</th>
            <th className="px-5 py-3 font-bold text-right">NRR</th><th className="px-5 py-3 font-bold text-right">NRR tertimbang</th>
            <th className="px-5 py-3 font-bold text-right">Indeks</th><th className="px-5 py-3 font-bold">Mutu</th>
          </tr></thead>
          <tbody className="divide-y divide-[color:var(--line-soft)]">
            {d.unsur.map(u=>{ const mm=mutu(u.indeks);
              return <tr key={u.kode} className="hover:bg-[var(--ground)]">
                <td className="px-5 py-3 font-num text-muted">{u.kode}</td>
                <td className="px-5 py-3">{u.nama}</td>
                <td className="px-5 py-3 font-num text-right">{fmt(u.nrr)}</td>
                <td className="px-5 py-3 font-num text-right text-muted">{fmt(u.nrr_tertimbang,3)}</td>
                <td className="px-5 py-3 font-num text-right font-semibold">{fmt(u.indeks)}</td>
                <td className="px-5 py-3"><Badge tone={mm.tone}>{mm.huruf}</Badge></td>
              </tr>;
            })}
            <tr className="bg-brandsoft font-semibold">
              <td className="px-5 py-3.5" colSpan="2">Total</td>
              <td className="px-5 py-3.5 font-num text-right">{fmt(d.nrr)}</td>
              <td className="px-5 py-3.5 font-num text-right">{fmt(d.nrr/d.unsur.length,3)}</td>
              <td className="px-5 py-3.5 font-num text-right">{fmt(d.ikm)}</td>
              <td className="px-5 py-3.5"><Badge tone={m.tone}>{m.huruf}</Badge></td>
            </tr>
          </tbody>
        </table>
      </div>}
    </Card>}
  </div>;
}

export default Laporan;
