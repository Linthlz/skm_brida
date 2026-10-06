import { DataTable } from '../../components/ui/index.js';
import { daftarResponden } from '../../services/api/admin.js';
import { warnaKategori } from '../../data/skala.js';
import { fmt, tglID } from '../../utils/format.js';

function DataSurvei() {
  const cols=[
    {key:'nomor', label:'ID Survei', mono:true},
    {key:'layanan', label:'Layanan'},
    {key:'nilai', label:'Nilai', mono:true, align:'right', cell:r=><span className="font-semibold">{fmt(r.nilai)}</span>},
    {key:'kategori', label:'Kategori', sortable:false, cell:r=>(
      <span className="inline-flex items-center gap-2 text-[13px]"><span className="w-2.5 h-2.5 rounded-[3px]" style={{background:warnaKategori(r.kategori)}}/>{r.kategori}</span>
    )},
    {key:'tanggal', label:'Tanggal', cell:r=><span className="whitespace-nowrap">{tglID(r.tanggal)}</span>},
  ];
  return <DataTable columns={cols} fetchPage={daftarResponden} defaultSort={{key:'tanggal', dir:'desc'}} pageSize={10}/>;
}

export default DataSurvei;
