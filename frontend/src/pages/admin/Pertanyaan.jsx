import { useState } from 'react';
import { Fragment } from 'react';
import { Badge, Button, Card, EmptyState, Field, GagalMuat, Input, Memuat, Modal, Textarea } from '../../components/ui/index.js';
import Icon from '../../components/Icon.jsx';
import { useToast } from '../../hooks/useToast.jsx';
import { useApi } from '../../hooks/useApi.js';
import * as adminApi from '../../services/api/admin.js';
import { formSurvei, hapusCache } from '../../services/api/publik.js';

function Pertanyaan() {
  const toast = useToast();
  const { data:list, error, reload, setData:setList } = useApi(adminApi.daftarPertanyaan);
  const form = useApi(formSurvei);
  const maks = form.data?.periode.skala_maks ?? 5;
  const [edit,setEdit] = useState(null);
  const [hapus,setHapus] = useState(null);
  const [draft,setDraft] = useState({unsur:'',teks:'',wajib:true});
  const [galat,setGalat] = useState({});
  const [sibuk,setSibuk] = useState(false);

  function buka(p){ setGalat({}); setEdit(p||{id:null}); setDraft(p?{unsur:p.unsur,teks:p.teks,wajib:p.wajib}:{unsur:'',teks:'',wajib:true}); }

  async function simpan(){
    const g={};
    if(!draft.unsur.trim()) g.unsur='Unsur pelayanan wajib diisi.';
    if(!draft.teks.trim()) g.teks='Teks pertanyaan wajib diisi.';
    setGalat(g); if(Object.keys(g).length) return;
    setSibuk(true);
    try{
      if(edit.id){
        const p = await adminApi.ubahPertanyaan(edit.id, draft);
        setList(l=>l.map(x=>x.id===p.id?p:x)); hapusCache(); toast('Pertanyaan diperbarui');
      } else {
        const p = await adminApi.tambahPertanyaan(draft);
        setList(l=>[...l,p]); hapusCache(); toast('Pertanyaan ditambahkan');
      }
      setEdit(null);
    }catch(ex){
      if(ex.status===422) setGalat({unsur:ex.field('unsur'), teks:ex.field('teks')});
      else toast(ex.message,'galat');
    }finally{ setSibuk(false); }
  }

  async function toggleAktif(p){
    setList(l=>l.map(x=>x.id===p.id?{...x,aktif:!p.aktif}:x));
    try{ await adminApi.ubahPertanyaan(p.id, {aktif:!p.aktif}); hapusCache(); toast(p.aktif?'Pertanyaan dinonaktifkan':'Pertanyaan diaktifkan'); }
    catch(ex){ setList(l=>l.map(x=>x.id===p.id?{...x,aktif:p.aktif}:x)); toast(ex.message,'galat'); }
  }

  async function pindah(i,arah){
    const j=i+arah; if(j<0||j>=list.length) return;
    const lama=list;
    const baru=[...list]; [baru[i],baru[j]]=[baru[j],baru[i]];
    setList(baru.map((x,k)=>({...x,urutan:k+1})));
    try{ setList(await adminApi.urutkanPertanyaan(baru.map(x=>x.id))); hapusCache(); }
    catch(ex){ setList(lama); toast(ex.message,'galat'); }
  }

  async function konfirmasiHapus(){
    setSibuk(true);
    try{ await adminApi.hapusPertanyaan(hapus.id); setList(l=>l.filter(x=>x.id!==hapus.id)); hapusCache(); setHapus(null); toast('Pertanyaan dihapus'); }
    catch(ex){ toast(ex.message,'galat'); }
    finally{ setSibuk(false); }
  }

  if(error) return <Card><GagalMuat error={error} onRetry={reload}/></Card>;
  if(!list) return <div className="space-y-2.5">{[0,1,2,3,4].map(i=><Memuat key={i} tinggi={92}/>)}</div>;

  return <div className="space-y-4">
    <div className="flex flex-col sm:flex-row sm:items-center gap-3">
      <p className="text-[14px] text-muted flex-1">Urutan pertanyaan menentukan tampilan pada halaman survei. Pertanyaan nonaktif tidak ditampilkan kepada responden.</p>
      <Button icon="plus" onClick={()=>buka(null)}>Tambah Pertanyaan</Button>
    </div>
    {list.length===0 && <Card><EmptyState icon="listChecks" title="Belum ada pertanyaan" desc="Tambahkan pertanyaan agar survei dapat diisi responden."/></Card>}
    <ul className="space-y-2.5">
      {list.map((p,i)=>(
        <li key={p.id} className={`bg-surface border rounded-card p-4 sm:p-5 flex flex-col sm:flex-row sm:items-center gap-4 ${p.aktif?'border-line':'border-line opacity-60'}`}>
          <div className="flex sm:flex-col gap-1">
            <button onClick={()=>pindah(i,-1)} disabled={i===0} aria-label="Naikkan urutan" className="p-1 rounded-[5px] text-muted hover:text-ink hover:bg-[var(--ground)] disabled:opacity-30"><Icon name="chevronUp" className="w-4 h-4"/></button>
            <button onClick={()=>pindah(i,1)} disabled={i===list.length-1} aria-label="Turunkan urutan" className="p-1 rounded-[5px] text-muted hover:text-ink hover:bg-[var(--ground)] disabled:opacity-30"><Icon name="chevronDown" className="w-4 h-4"/></button>
          </div>
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2 mb-1.5">
              <span className="font-num text-[11.5px] font-semibold text-[var(--brand)] bg-brandsoft rounded-[5px] px-1.5 py-0.5">{p.kode}</span>
              <span className="text-[12px] text-muted">{p.unsur}</span>
              <Badge tone="neutral">Skala 1–{maks}</Badge>
              {p.wajib && <Badge tone="gold">Wajib</Badge>}
            </div>
            <p className="text-[14.5px] leading-snug">{p.teks}</p>
          </div>
          <div className="flex items-center gap-1.5 shrink-0">
            <button onClick={()=>toggleAktif(p)}
              className={`relative w-[42px] h-[24px] rounded-full transition-colors shrink-0 ${p.aktif?'bg-[var(--brand)]':'bg-[var(--line)]'}`}
              role="switch" aria-checked={p.aktif} aria-label="Aktifkan pertanyaan">
              <span className={`absolute top-[3px] w-[18px] h-[18px] rounded-full bg-white transition-all ${p.aktif?'left-[21px]':'left-[3px]'}`}/>
            </button>
            <button onClick={()=>buka(p)} aria-label="Edit" className="p-2 rounded-[7px] text-muted hover:text-[var(--brand)] hover:bg-brandsoft"><Icon name="pencil" className="w-[17px] h-[17px]"/></button>
            <button onClick={()=>setHapus(p)} aria-label="Hapus" className="p-2 rounded-[7px] text-muted hover:text-[var(--bad)] hover:bg-[color:var(--bad)]/10"><Icon name="trash" className="w-[17px] h-[17px]"/></button>
          </div>
        </li>
      ))}
    </ul>

    <Modal open={!!edit} onClose={()=>setEdit(null)} title={edit?.id?'Edit pertanyaan':'Tambah pertanyaan'}
      desc={`Pertanyaan memakai skala penilaian 1–${maks}.`}
      footer={<Fragment><Button variant="outline" onClick={()=>setEdit(null)}>Batal</Button><Button icon="check" onClick={simpan} disabled={sibuk}>{sibuk?'Menyimpan…':'Simpan'}</Button></Fragment>}>
      <div className="space-y-5">
        <Field label="Unsur pelayanan" required error={galat.unsur}><Input id="p-unsur" maxLength={100} value={draft.unsur} onChange={e=>setDraft({...draft,unsur:e.target.value})} placeholder="Contoh: Waktu Penyelesaian"/></Field>
        <Field label="Teks pertanyaan" required error={galat.teks}><Textarea id="p-teks" maxLength={500} className="!min-h-[96px]" value={draft.teks} onChange={e=>setDraft({...draft,teks:e.target.value})} placeholder="Bagaimana pendapat Anda tentang …"/></Field>
        <label className="flex items-center gap-3 cursor-pointer text-[14.5px]">
          <input type="checkbox" checked={draft.wajib} onChange={e=>setDraft({...draft,wajib:e.target.checked})} className="w-[17px] h-[17px] accent-[var(--brand)]"/> Wajib dijawab responden
        </label>
      </div>
    </Modal>

    <Modal open={!!hapus} onClose={()=>setHapus(null)} title="Hapus pertanyaan?" desc="Pertanyaan tidak lagi tampil di survei. Jawaban yang sudah masuk tetap tersimpan untuk rekap periode sebelumnya."
      footer={<Fragment><Button variant="outline" onClick={()=>setHapus(null)}>Batal</Button>
        <Button variant="danger" icon="trash" onClick={konfirmasiHapus} disabled={sibuk}>Hapus</Button></Fragment>}>
      <p className="text-[14.5px] text-ink2 leading-relaxed">{hapus?.teks}</p>
    </Modal>
  </div>;
}

export default Pertanyaan;
