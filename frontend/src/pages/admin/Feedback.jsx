import { Fragment, useState } from 'react';
import { Badge, Button, Card, EmptyState, Field, GagalMuat, Memuat, Modal, Paginasi, Select, Textarea } from '../../components/ui/index.js';
import Icon from '../../components/Icon.jsx';
import { useToast } from '../../hooks/useToast.jsx';
import { useApi } from '../../hooks/useApi.js';
import * as adminApi from '../../services/api/admin.js';
import { STATUS_FEEDBACK, STATUS_TONE } from '../../data/feedback.js';
import { tglID } from '../../utils/format.js';

const TAB = [['feedback','Feedback survei','message'], ['pesan','Pesan kontak','mail']];

function Feedback() {
  const [tab,setTab] = useState('feedback');
  return <div className="space-y-4">
    <div role="tablist" className="inline-flex p-1 bg-surface border border-line rounded-[10px] gap-1">
      {TAB.map(([k,l,i])=>(
        <button key={k} role="tab" aria-selected={tab===k} onClick={()=>setTab(k)}
          className={`inline-flex items-center gap-2 px-3.5 py-2 rounded-[7px] text-[13.5px] font-semibold transition-colors ${tab===k?'bg-brandsoft text-[var(--brand)]':'text-muted hover:text-ink'}`}>
          <Icon name={i} className="w-4 h-4"/>{l}
        </button>
      ))}
    </div>
    {tab==='feedback' ? <DaftarFeedback/> : <DaftarPesan/>}
  </div>;
}

function DaftarFeedback() {
  const toast = useToast();
  const [filter,setFilter] = useState('');
  const [page,setPage] = useState(1);
  const [buka,setBuka] = useState(null);
  const res = useApi(()=>adminApi.daftarFeedback({status:filter, page, per_page:10}), [filter, page]);
  const kategori = useApi(adminApi.daftarPertanyaan);
  const [sibuk,setSibuk] = useState(null);

  const pilihFilter = s => { setFilter(filter===s?'':s); setPage(1); };
  const hitung = res.data?.meta.hitung;

  async function majukan(f){
    setSibuk(f.id);
    try{
      const baru = await adminApi.ubahFeedback(f.id, {status:f.status_berikut});
      toast(`${baru.nomor} → ${baru.status}`);
      res.reload();
      if(buka?.id===baru.id) setBuka(baru);
    }catch(ex){ toast(ex.message,'galat'); }
    finally{ setSibuk(null); }
  }

  return <div className="space-y-4">
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
      {STATUS_FEEDBACK.map(s=>(
        <button key={s} onClick={()=>pilihFilter(s)} aria-pressed={filter===s}
          className={`text-left bg-surface border rounded-card p-4 transition-colors ${filter===s?'border-[var(--brand)] ring-1 ring-[color:var(--brand)]/30':'border-line hover:border-[var(--brand)]'}`}>
          <div className="flex items-center gap-2 mb-2"><span className="w-2 h-2 rounded-full" style={{background:`var(--${STATUS_TONE[s]})`}}/><span className="text-[12.5px] text-muted">{s}</span></div>
          <p className="font-num text-[24px] font-semibold leading-none">{hitung ? hitung[s] : '–'}</p>
        </button>
      ))}
    </div>
    {filter && <div><Button size="sm" variant="quiet" icon="x" onClick={()=>pilihFilter(filter)}>Tampilkan semua status</Button></div>}

    {res.error ? <Card><GagalMuat error={res.error} onRetry={res.reload}/></Card>
    : !res.data ? <div className="space-y-3">{[0,1,2].map(i=><Memuat key={i} tinggi={150}/>)}</div>
    : res.data.data.length===0 ? <Card><EmptyState icon="message" title="Belum ada feedback pada status ini" desc="Pilih status lain atau tampilkan seluruh masukan yang masuk."/></Card>
    : <ul className={`space-y-3 transition-opacity ${res.loading?'opacity-60':''}`}>
      {res.data.data.map(f=>(
        <li key={f.id} className="bg-surface border border-line rounded-card p-5">
          <div className="flex flex-wrap items-center gap-2.5 mb-3">
            <span className="font-num text-[12.5px] text-muted">{f.nomor}</span>
            <span className="text-[12.5px] text-muted">·</span>
            <span className="text-[12.5px] text-muted">{tglID(f.tanggal)}</span>
            {f.kategori && <Badge tone="neutral">{f.kategori.nama}</Badge>}
            <Badge tone={STATUS_TONE[f.status]} dot className="ml-auto">{f.status}</Badge>
          </div>
          <p className="text-[15px] text-ink2 leading-relaxed" style={{maxWidth:'72ch'}}>{f.isi}</p>
          <div className="mt-4 pt-4 border-t border-linesoft flex flex-wrap items-center gap-2">
            <span className="text-[12.5px] text-muted mr-1">{f.layanan}</span>
            <div className="ml-auto flex flex-wrap gap-2">
              <Button size="sm" variant="outline" icon="eye" onClick={()=>setBuka(f)}>Detail</Button>
              {f.status_berikut && <Button size="sm" disabled={sibuk===f.id} onClick={()=>majukan(f)}>Tandai {f.status_berikut}</Button>}
            </div>
          </div>
        </li>
      ))}
    </ul>}
    <Paginasi meta={res.data?.meta} onPage={setPage} disabled={res.loading}/>

    <DetailFeedback f={buka} kategori={kategori.data||[]} onClose={()=>setBuka(null)} onMaju={majukan} sibuk={sibuk===buka?.id}
      onSimpan={baru=>{ setBuka(baru); res.reload(); }}/>
  </div>;
}

function DetailFeedback({f, kategori, onClose, onMaju, onSimpan, sibuk}) {
  const toast = useToast();
  const [draf,setDraf] = useState({id:null, catatan:'', kategori_id:''});
  const [menyimpan,setMenyimpan] = useState(false);
  // Isi ulang draf setiap kali feedback lain dibuka.
  if(f && draf.id!==f.id) setDraf({id:f.id, catatan:f.catatan||'', kategori_id:f.kategori?String(f.kategori.id):''});

  async function simpan(){
    setMenyimpan(true);
    try{
      const baru = await adminApi.ubahFeedback(f.id, {catatan:draf.catatan.trim()||null, kategori_id:draf.kategori_id?Number(draf.kategori_id):null});
      toast('Catatan tindak lanjut disimpan'); onSimpan(baru);
    }catch(ex){ toast(ex.status===422 ? Object.values(ex.errors).flat()[0] : ex.message,'galat'); }
    finally{ setMenyimpan(false); }
  }

  return <Modal open={!!f} onClose={onClose} title={f?.nomor} desc="Detail masukan masyarakat" width="max-w-xl"
    footer={<Fragment>
      <Button variant="outline" onClick={onClose}>Tutup</Button>
      {f?.status_berikut && <Button variant="outline" disabled={sibuk} onClick={()=>onMaju(f)}>Tandai {f.status_berikut}</Button>}
      <Button icon="check" onClick={simpan} disabled={menyimpan}>{menyimpan?'Menyimpan…':'Simpan'}</Button>
    </Fragment>}>
    {f && <div className="space-y-4">
      <div className="flex flex-wrap gap-2"><Badge tone={STATUS_TONE[f.status]} dot>{f.status}</Badge><Badge tone="neutral">{tglID(f.tanggal)}</Badge><Badge tone="neutral">{f.nomor_responden}</Badge></div>
      <p className="text-[15px] leading-relaxed">{f.isi}</p>
      <p className="text-[13px] text-muted">{f.layanan}{f.ditangani_oleh && ` · terakhir ditangani ${f.ditangani_oleh}`}</p>
      <Field label="Kategori unsur" hint="Diisi otomatis dengan unsur bernilai terendah dari responden ini.">
        <Select id="fb-kat" options={kategori.map(k=>({value:String(k.id), label:`${k.kode} · ${k.unsur}`}))} placeholder="Tanpa kategori"
          value={draf.kategori_id} onChange={e=>setDraf({...draf, kategori_id:e.target.value})}/>
      </Field>
      <Field label="Catatan tindak lanjut"><Textarea id="fb-note" maxLength={2000} className="!min-h-[90px]" value={draf.catatan} onChange={e=>setDraf({...draf, catatan:e.target.value})} placeholder="Tuliskan langkah tindak lanjut yang sudah atau akan dilakukan."/></Field>
    </div>}
  </Modal>;
}

function DaftarPesan() {
  const toast = useToast();
  const [page,setPage] = useState(1);
  const [buka,setBuka] = useState(null);
  const res = useApi(()=>adminApi.daftarPesan({page, per_page:10}), [page]);

  async function tandai(p, dibaca){
    try{
      const baru = await adminApi.ubahPesan(p.id, {dibaca});
      res.setData(r=>({...r, data:r.data.map(x=>x.id===p.id?baru:x), meta:{...r.meta, belum_dibaca:r.meta.belum_dibaca+(dibaca?-1:1)}}));
      if(buka?.id===p.id) setBuka(baru);
    }catch(ex){ toast(ex.message,'galat'); }
  }
  function lihat(p){ setBuka(p); if(!p.dibaca) tandai(p, true); }

  if(res.error) return <Card><GagalMuat error={res.error} onRetry={res.reload}/></Card>;
  if(!res.data) return <div className="space-y-3">{[0,1,2].map(i=><Memuat key={i} tinggi={110}/>)}</div>;

  const {data, meta} = res.data;
  return <div className="space-y-4">
    <p className="text-[14px] text-muted"><span className="font-num font-semibold text-ink">{meta.belum_dibaca}</span> pesan belum dibaca. Pesan dikirim melalui formulir Kontak di situs publik.</p>
    {data.length===0 ? <Card><EmptyState icon="mail" title="Belum ada pesan" desc="Pesan dari formulir kontak akan muncul di sini."/></Card>
    : <ul className={`space-y-3 transition-opacity ${res.loading?'opacity-60':''}`}>
      {data.map(p=>(
        <li key={p.id} className={`bg-surface border rounded-card p-5 ${p.dibaca?'border-line':'border-[color:var(--brand)]/40'}`}>
          <div className="flex flex-wrap items-center gap-2.5 mb-2">
            {!p.dibaca && <span className="w-2 h-2 rounded-full bg-[var(--brand)]" aria-label="Belum dibaca"/>}
            <span className="text-[14px] font-semibold">{p.nama}</span>
            <span className="text-[12.5px] text-muted">{p.email}</span>
            <Badge tone="neutral" className="ml-auto">{p.topik}</Badge>
          </div>
          <p className="text-[14.5px] text-ink2 leading-relaxed line-clamp-2" style={{maxWidth:'72ch'}}>{p.pesan}</p>
          <div className="mt-3 flex flex-wrap items-center gap-2">
            <span className="text-[12.5px] text-muted">{tglID(p.tanggal)}</span>
            <div className="ml-auto flex gap-2">
              <Button size="sm" variant="quiet" onClick={()=>tandai(p, !p.dibaca)}>{p.dibaca?'Tandai belum dibaca':'Tandai dibaca'}</Button>
              <Button size="sm" variant="outline" icon="eye" onClick={()=>lihat(p)}>Baca</Button>
            </div>
          </div>
        </li>
      ))}
    </ul>}
    <Paginasi meta={meta} onPage={setPage} disabled={res.loading}/>

    <Modal open={!!buka} onClose={()=>setBuka(null)} title={buka?.topik} desc={buka && `${buka.nama} · ${tglID(buka.tanggal)}`} width="max-w-xl"
      footer={<Fragment>
        <Button variant="outline" onClick={()=>setBuka(null)}>Tutup</Button>
        {buka && <Button as="a" icon="mail" href={`mailto:${encodeURIComponent(buka.email)}?subject=${encodeURIComponent('Balasan: '+buka.topik)}`}>Balas lewat email</Button>}
      </Fragment>}>
      {buka && <div className="space-y-3">
        <p className="text-[13.5px] text-muted">{buka.email}</p>
        <p className="text-[15px] leading-relaxed whitespace-pre-line">{buka.pesan}</p>
      </div>}
    </Modal>
  </div>;
}

export default Feedback;
