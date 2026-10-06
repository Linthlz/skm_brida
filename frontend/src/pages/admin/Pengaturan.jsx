import { useState } from 'react';
import { Button, Card, Field, GagalMuat, Input, Memuat, Select } from '../../components/ui/index.js';
import { useToast } from '../../hooks/useToast.jsx';
import { useApi } from '../../hooks/useApi.js';
import { pengaturan, simpanPengaturan } from '../../services/api/admin.js';
import { hapusCache } from '../../services/api/publik.js';

const OPSI_SKALA = [{value:'4', label:'1 – 4 (Permen PANRB 14/2017)'}, {value:'5', label:'1 – 5'}];

function Pengaturan() {
  const { data, error, reload } = useApi(pengaturan);
  if (error) return <Card><GagalMuat error={error} onRetry={reload}/></Card>;
  if (!data) return <div className="grid lg:grid-cols-2 gap-4"><Memuat tinggi={360}/><Memuat tinggi={360}/></div>;
  return <FormPengaturan awal={data}/>;
}

function FormPengaturan({awal}) {
  const toast = useToast();
  const [asal,setAsal] = useState(awal);
  const [s,setS] = useState(()=>keForm(awal));
  const [galat,setGalat] = useState({});
  const [menyimpan,setMenyimpan] = useState(false);

  // Periode lama + tahun berikutnya, agar periode baru bisa diaktifkan dari sini.
  const tahunDepan = Math.max(...asal.opsi_periode, new Date().getFullYear()) + 1;
  const opsiPeriode = [...new Set([tahunDepan, ...asal.opsi_periode])].sort((a,b)=>b-a).map(String);
  const skalaTerkunci = asal.skala_terkunci && String(asal.periode)===s.periode;

  async function simpan(){
    setMenyimpan(true); setGalat({});
    try{
      const baru = await simpanPengaturan({...s, periode:Number(s.periode), skala_maks:Number(s.skala_maks), tanggal_tutup:s.tanggal_tutup||null});
      setAsal(baru); setS(keForm(baru)); hapusCache();
      toast('Pengaturan disimpan');
    }catch(ex){
      if(ex.status===422) setGalat(Object.fromEntries(Object.keys(ex.errors).map(k=>[k, ex.field(k)])));
      toast(ex.status===422 ? 'Periksa kembali isian pengaturan.' : ex.message, 'galat');
    }finally{ setMenyimpan(false); }
  }

  return <div className="grid lg:grid-cols-2 gap-4">
    <Card pad="p-6">
      <h3 className="font-display text-[17px] font-bold mb-1">Identitas survei</h3>
      <p className="text-[13.5px] text-muted mb-6">Tampil pada halaman publik dan pada kop laporan.</p>
      <div className="space-y-5">
        <Field label="Judul survei" required error={galat.judul}><Input id="g-judul" maxLength={150} value={s.judul} onChange={e=>setS({...s,judul:e.target.value})}/></Field>
        <div className="grid sm:grid-cols-2 gap-4">
          <Field label="Periode aktif" required error={galat.periode}><Select id="g-per" placeholder={null} options={opsiPeriode} value={s.periode} onChange={e=>setS({...s,periode:e.target.value})}/></Field>
          <Field label="Skala penilaian" required error={galat.skala_maks}
            hint={skalaTerkunci ? 'Terkunci karena periode ini sudah memiliki responden.' : undefined}>
            <Select id="g-skala" placeholder={null} options={OPSI_SKALA} value={s.skala_maks} disabled={skalaTerkunci} onChange={e=>setS({...s,skala_maks:e.target.value})}/>
          </Field>
        </div>
        <Field label="Tanggal penutupan" error={galat.tanggal_tutup} hint="Survei tidak menerima jawaban setelah tanggal ini. Kosongkan bila tidak dibatasi.">
          <Input id="g-tutup" type="date" value={s.tanggal_tutup} onChange={e=>setS({...s,tanggal_tutup:e.target.value})}/>
        </Field>
      </div>
    </Card>
    <div className="space-y-4">
      <Card pad="p-6">
        <h3 className="font-display text-[17px] font-bold mb-1">Preferensi</h3>
        <p className="text-[13.5px] text-muted mb-5">Berlaku untuk pengisian survei berikutnya.</p>
        <ul className="divide-y divide-[color:var(--line-soft)]">
          {[['anonim','Izinkan pengisian tanpa nama','Responden dapat mengosongkan kolom nama dan kontak.'],
            ['publik','Publikasikan hasil ke halaman publik','Statistik agregat tampil pada menu Hasil Survei.']].map(([k,j,d])=>(
            <li key={k} className="flex items-start gap-4 py-4 first:pt-0">
              <div className="flex-1"><p className="text-[14.5px] font-semibold">{j}</p><p className="text-[13px] text-muted mt-0.5">{d}</p></div>
              <button onClick={()=>setS(p=>({...p,[k]:!p[k]}))} role="switch" aria-checked={s[k]} aria-label={j}
                className={`relative w-[42px] h-[24px] rounded-full shrink-0 transition-colors ${s[k]?'bg-[var(--brand)]':'bg-[var(--line)]'}`}>
                <span className={`absolute top-[3px] w-[18px] h-[18px] rounded-full bg-white transition-all ${s[k]?'left-[21px]':'left-[3px]'}`}/>
              </button>
            </li>
          ))}
        </ul>
      </Card>
      <Card pad="p-6">
        <h3 className="font-display text-[17px] font-bold mb-4">Pengelola</h3>
        <ul className="space-y-3">
          {asal.pengelola.map(p=>(
            <li key={p.id} className="flex items-center gap-3">
              <span className="w-9 h-9 rounded-full grid place-items-center bg-brandsoft text-[var(--brand)] font-display font-extrabold text-[13px]">{p.inisial}</span>
              <span className="flex-1 min-w-0"><span className="block text-[14px] font-semibold truncate">{p.nama}</span><span className="block text-[12.5px] text-muted">{p.jabatan || 'Admin'}</span></span>
            </li>
          ))}
        </ul>
      </Card>
      <Button size="lg" icon="check" className="w-full" onClick={simpan} disabled={menyimpan}>{menyimpan?'Menyimpan…':'Simpan Pengaturan'}</Button>
    </div>
  </div>;
}

const keForm = (d) => ({
  judul:d.judul, periode:String(d.periode), skala_maks:String(d.skala_maks),
  tanggal_tutup:d.tanggal_tutup||'', anonim:d.anonim, publik:d.publik,
});

export default Pengaturan;
