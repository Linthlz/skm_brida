import { useState } from 'react';
import { Button, Card, DataTable, GagalMuat, Memuat, Modal, Select } from '../../components/ui/index.js';
import Icon from '../../components/Icon.jsx';
import { useApi } from '../../hooks/useApi.js';
import { daftarResponden, detailResponden } from '../../services/api/admin.js';
import { referensi } from '../../services/api/publik.js';
import { fmt, tglID } from '../../utils/format.js';

function Responden() {
  const [detailId,setDetailId] = useState(null);
  const [layanan,setLayanan] = useState('');
  const ref = useApi(referensi);
  const cols = [
    {key:'nomor', label:'ID', mono:true},
    {key:'usia', label:'Usia', mono:true, align:'right'},
    {key:'jk', label:'Jenis Kelamin'},
    {key:'pekerjaan', label:'Pekerjaan'},
    {key:'layanan', label:'Layanan'},
    {key:'nilai', label:'Nilai', mono:true, align:'right', cell:r=><span className="font-semibold">{fmt(r.nilai)}</span>},
    {key:'tanggal', label:'Tanggal', cell:r=><span className="whitespace-nowrap">{tglID(r.tanggal)}</span>},
  ];
  return <div className="space-y-4">
    <Card pad="p-4 sm:p-5">
      <div className="flex flex-col sm:flex-row gap-3 sm:items-center">
        <span className="flex items-center gap-2 text-[13px] font-semibold text-muted"><Icon name="filter" className="w-4 h-4"/>Filter</span>
        <div className="sm:w-[260px]"><Select aria-label="Jenis layanan" options={(ref.data?.layanan||[]).map(l=>({value:String(l.id), label:l.nama}))} placeholder="Semua jenis layanan" value={layanan} onChange={e=>setLayanan(e.target.value)}/></div>
        {layanan && <Button size="sm" variant="quiet" icon="x" onClick={()=>setLayanan('')}>Atur ulang</Button>}
        <p className="sm:ml-auto text-[12.5px] text-muted">Nama dan kontak responden tidak ditampilkan pada tabel ini.</p>
      </div>
    </Card>
    <DataTable columns={cols} fetchPage={daftarResponden} params={{layanan_id:layanan||undefined}}
      defaultSort={{key:'tanggal', dir:'desc'}} onRow={r=>setDetailId(r.id)} pageSize={10}/>
    <DetailResponden id={detailId} onClose={()=>setDetailId(null)}/>
  </div>;
}

function DetailResponden({id, onClose}) {
  const { data:d, error, reload } = useApi(()=> id ? detailResponden(id) : Promise.resolve(null), [id]);
  const muat = id && !d && !error;
  return <Modal open={!!id} onClose={onClose} title={d?.nomor || 'Detail responden'} desc="Detail responden dan penilaian" width="max-w-xl"
    footer={<Button variant="outline" onClick={onClose}>Tutup</Button>}>
    {error ? <GagalMuat error={error} onRetry={reload}/>
    : muat || d?.id !== id ? <Memuat baris={8}/>
    : <div className="space-y-5">
      <dl className="divide-y divide-[color:var(--line-soft)]">
        {[['Usia', d.usia+' tahun'],['Jenis kelamin',d.jk],['Pendidikan',d.pendidikan],['Pekerjaan',d.pekerjaan],
          ['Kecamatan',d.kecamatan||'—'],['Jenis layanan',d.layanan],['Frekuensi',d.frekuensi],
          ['Tanggal pengisian',tglID(d.tanggal)],['Nilai rata-rata',`${fmt(d.nilai)} / ${fmt(d.skala_maks)}`],['Konversi indeks',fmt(d.indeks)],['Kategori',d.kategori]].map(([k,v])=>(
          <div key={k} className="flex gap-4 py-2.5 text-[14px]"><dt className="text-muted w-[150px] shrink-0">{k}</dt><dd className="font-medium min-w-0">{v}</dd></div>
        ))}
      </dl>
      <div>
        <p className="eyebrow text-[10.5px] font-bold text-muted mb-2">Nilai per unsur</p>
        <ul className="grid sm:grid-cols-2 gap-x-6">
          {d.jawaban.map(j=>(
            <li key={j.kode} className="flex items-center gap-3 py-1.5 text-[13.5px] border-b border-linesoft">
              <span className="font-num text-[11.5px] text-muted w-7 shrink-0">{j.kode}</span>
              <span className="flex-1 min-w-0 truncate">{j.unsur}</span>
              <span className="font-num font-semibold">{j.nilai}</span>
            </li>
          ))}
        </ul>
      </div>
      {(d.saran || d.apresiasi) && <div className="space-y-3">
        {d.saran && <div><p className="eyebrow text-[10.5px] font-bold text-muted mb-1">Saran</p><p className="text-[14px] text-ink2 leading-relaxed">{d.saran}</p></div>}
        {d.apresiasi && <div><p className="eyebrow text-[10.5px] font-bold text-muted mb-1">Apresiasi</p><p className="text-[14px] text-ink2 leading-relaxed">{d.apresiasi}</p></div>}
      </div>}
    </div>}
  </Modal>;
}

export default Responden;
