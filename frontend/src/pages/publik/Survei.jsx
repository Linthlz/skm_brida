import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Card, EmptyState, GagalMuat, Memuat, SectionHead, ProgressSteps, RatingScale, Field, Input, Textarea, Select } from '../../components/ui/index.js';
import { LegendDot } from '../../components/charts/index.js';
import Icon from '../../components/Icon.jsx';
import { useToast } from '../../hooks/useToast.jsx';
import { useApi } from '../../hooks/useApi.js';
import { formSurvei, kirimSurvei, referensi } from '../../services/api/publik.js';
import { skalaDari } from '../../data/skala.js';
import { tglID } from '../../utils/format.js';
import Shell from '../../layouts/Shell.jsx';

const STEPS = ['Identitas', 'Penilaian', 'Saran', 'Selesai'];
const DRAF = 'skm-draf-survei';
const AWAL = {
  nama:'', jk:'', usia:'', pendidikan:'', pekerjaan:'', kecamatan:'', layanan_id:'', frekuensi:'',
  nilai:{}, saran:'', apresiasi:'', setuju:false, website:'',
};
// Field langkah 1; dipakai untuk mengarahkan galat validasi server ke langkah yang tepat.
const FIELD_IDENTITAS = ['nama','jk','usia','pendidikan','pekerjaan','kecamatan','layanan_id','frekuensi'];

/** Draf jawaban disimpan di sessionStorage selama tab terbuka, agar tidak hilang saat halaman dimuat ulang. */
function bacaDraf(){
  try { const d = JSON.parse(sessionStorage.getItem(DRAF)); return d && d.f ? d : null; } catch { return null; }
}
function simpanDraf(f, step){
  try { sessionStorage.setItem(DRAF, JSON.stringify({f:{...f, website:''}, step})); } catch { /* penyimpanan diblokir: abaikan */ }
}
function hapusDraf(){ try { sessionStorage.removeItem(DRAF); } catch { /* abaikan */ } }

const muatData = () => Promise.all([formSurvei(), referensi()]).then(([form, ref]) => ({form, ref}));

function Survei() {
  const navigate = useNavigate();
  const toast = useToast();
  const { data, error, reload } = useApi(muatData);
  const draf = useRef(bacaDraf()).current;
  const [step,setStep] = useState(draf?.step ?? 0);
  const [err,setErr] = useState({});
  const [f,setF] = useState({...AWAL, ...(draf?.f||{})});
  const [mengirim,setMengirim] = useState(false);
  const set = (k,v)=> setF(p=>({...p,[k]:v}));
  const setNilai = (id,v)=> setF(p=>({...p, nilai:{...p.nilai,[id]:v}}));
  const topRef = useRef(null);
  const naik = ()=> topRef.current && topRef.current.scrollIntoView({behavior:'smooth', block:'start'});

  useEffect(()=>{ simpanDraf(f, step); },[f, step]);

  if (error) return <Shell><section className="py-10 sm:py-14"><Card className="mx-auto" pad="p-2"><GagalMuat error={error} onRetry={reload}/></Card></section></Shell>;
  if (!data) return <Shell><section className="py-10 sm:py-14 mx-auto space-y-4" style={{maxWidth:'840px'}}><Memuat tinggi={60}/><Memuat tinggi={420}/></section></Shell>;

  const { form, ref } = data;
  const pertanyaan = form.pertanyaan;
  const skala = skalaDari(form.skala);
  const layananNama = ref.layanan.find(l=>String(l.id)===String(f.layanan_id))?.nama;

  if (form.periode.ditutup) {
    return <Shell><section className="py-16 sm:py-24"><Card className="mx-auto" pad="p-2">
      <EmptyState icon="calendar" title={`Survei periode ${form.periode.tahun} sudah ditutup`}
        desc={`Pengisian ditutup pada ${tglID(form.periode.tanggal_tutup)}. Terima kasih atas minat Anda untuk memberikan penilaian.`}
        action={<Button variant="outline" icon="home" onClick={()=>navigate('/')}>Kembali ke Beranda</Button>}/>
    </Card></section></Shell>;
  }

  function validasi(s){
    const e={};
    if(s===0){
      if(!form.anonim && !f.nama.trim()) e.nama='Nama wajib diisi.';
      if(!f.jk) e.jk='Silakan pilih salah satu jawaban sebelum melanjutkan.';
      if(!f.usia) e.usia='Usia wajib diisi.';
      else if(isNaN(+f.usia) || +f.usia<10 || +f.usia>99) e.usia='Masukkan usia antara 10 dan 99 tahun.';
      if(!f.pendidikan) e.pendidikan='Silakan pilih pendidikan terakhir Anda.';
      if(!f.pekerjaan) e.pekerjaan='Silakan pilih pekerjaan Anda.';
      if(!f.layanan_id) e.layanan_id='Silakan pilih layanan yang Anda gunakan.';
      if(!f.frekuensi) e.frekuensi='Silakan pilih salah satu jawaban sebelum melanjutkan.';
    }
    if(s===1){
      const belum = pertanyaan.filter(p=>p.wajib && !f.nilai[p.id]);
      if(belum.length) e.nilai = `Silakan pilih salah satu jawaban sebelum melanjutkan. ${belum.length} pertanyaan belum dinilai.`;
    }
    if(s===2){ if(!f.setuju) e.setuju='Centang persetujuan untuk mengirim jawaban Anda.'; }
    setErr(e); return Object.keys(e).length===0;
  }
  function lanjut(){ if(validasi(step)){ setStep(s=>s+1); naik(); } else naik(); }
  function kembali(){ setErr({}); setStep(s=>s-1); naik(); }

  async function kirim(){
    if(!validasi(2)){ naik(); return; }
    setMengirim(true);
    try{
      const hasil = await kirimSurvei({
        nama: f.nama.trim() || null, jk:f.jk, usia:Number(f.usia), pendidikan:f.pendidikan, pekerjaan:f.pekerjaan,
        kecamatan: f.kecamatan || null, layanan_id:Number(f.layanan_id), frekuensi:f.frekuensi,
        jawaban: pertanyaan.filter(p=>f.nilai[p.id]).map(p=>({pertanyaan_id:p.id, nilai:f.nilai[p.id]})),
        saran: f.saran.trim() || null, apresiasi: f.apresiasi.trim() || null, setuju:f.setuju, website:f.website,
      });
      hapusDraf();
      toast('Jawaban survei berhasil dikirim');
      navigate('/survei/selesai', { state: hasil });
    }catch(ex){
      if(ex.status===422){
        const e={}; let ke=2;
        Object.keys(ex.errors).forEach(k=>{
          if(FIELD_IDENTITAS.includes(k)){ e[k]=ex.field(k); ke=0; }
          else if(k.startsWith('jawaban')){ e.nilai = e.nilai || ex.field(k); ke=Math.min(ke,1); }
          else e[k]=ex.field(k);
        });
        setErr(e); setStep(ke); naik();
        toast('Ada isian yang perlu diperbaiki.', 'galat');
      } else {
        toast(ex.message, 'galat');
        if(ex.status===403) reload(); // periode baru saja ditutup
      }
    }finally{ setMengirim(false); }
  }

  const terisi = pertanyaan.filter(p=>f.nilai[p.id]).length;

  return <Shell>
    <div ref={topRef} className="scroll-mt-24"/>
    <section className="py-10 sm:py-14">
      <div className="mx-auto" style={{maxWidth:'840px'}}>
        <div className="mb-8"><ProgressSteps steps={STEPS} current={step}/></div>

        {step===0 && (
          <Card pad="p-6 sm:p-8" className="anim-rise">
            <SectionHead eyebrow="Langkah 01" title="Data responden"
              desc={form.anonim ? 'Data ini dipakai untuk analisis kelompok pengguna layanan. Nama boleh dikosongkan.' : 'Data ini dipakai untuk analisis kelompok pengguna layanan.'}/>
            <div className="grid sm:grid-cols-2 gap-5 mt-7">
              <Field label="Nama" required={!form.anonim} error={err.nama} className="sm:col-span-2">
                <Input id="s-nama" maxLength={100} value={f.nama} onChange={e=>set('nama',e.target.value)} placeholder={form.anonim?'Boleh dikosongkan':'Nama lengkap Anda'}/>
              </Field>
              <Field label="Jenis kelamin" required error={err.jk}>
                <div className="grid grid-cols-2 gap-2">
                  {ref.jk.map(o=>(
                    <button key={o} type="button" onClick={()=>set('jk',o)} aria-pressed={f.jk===o}
                      className={`py-2.5 rounded-[8px] border text-[14.5px] font-semibold transition-colors ${f.jk===o?'bg-brandsoft border-[var(--brand)] text-[var(--brand)]':'bg-surface border-line text-ink2 hover:border-[var(--brand)]'}`}>{o}</button>
                  ))}
                </div>
              </Field>
              <Field label="Usia" required error={err.usia} hint="Dalam tahun, contoh: 26">
                <Input id="s-usia" inputMode="numeric" value={f.usia} onChange={e=>set('usia',e.target.value.replace(/\D/g,'').slice(0,2))} placeholder="26"/>
              </Field>
              <Field label="Pendidikan terakhir" required error={err.pendidikan}>
                <Select id="s-pend" options={ref.pendidikan} value={f.pendidikan} onChange={e=>set('pendidikan',e.target.value)}/>
              </Field>
              <Field label="Pekerjaan" required error={err.pekerjaan}>
                <Select id="s-kerja" options={ref.pekerjaan} value={f.pekerjaan} onChange={e=>set('pekerjaan',e.target.value)}/>
              </Field>
              <Field label="Jenis layanan yang digunakan" required error={err.layanan_id}>
                <Select id="s-lay" options={ref.layanan.map(l=>({value:String(l.id), label:l.nama}))} value={f.layanan_id} onChange={e=>set('layanan_id',e.target.value)} placeholder="Pilih layanan BRIDA yang Anda gunakan"/>
              </Field>
              <Field label="Kecamatan domisili" error={err.kecamatan}>
                <Select id="s-kec" options={ref.kecamatan} value={f.kecamatan} onChange={e=>set('kecamatan',e.target.value)} placeholder="Pilih kecamatan"/>
              </Field>
              <Field label="Frekuensi menggunakan layanan" required error={err.frekuensi} className="sm:col-span-2">
                <div className="grid sm:grid-cols-4 gap-2">
                  {ref.frekuensi.map(o=>(
                    <button key={o} type="button" onClick={()=>set('frekuensi',o)} aria-pressed={f.frekuensi===o}
                      className={`py-2.5 px-2 rounded-[8px] border text-[13.5px] font-semibold transition-colors ${f.frekuensi===o?'bg-brandsoft border-[var(--brand)] text-[var(--brand)]':'bg-surface border-line text-ink2 hover:border-[var(--brand)]'}`}>{o}</button>
                  ))}
                </div>
              </Field>
            </div>
            <p className="mt-6 flex gap-2.5 text-[13px] text-muted bg-[var(--ground)] border border-linesoft rounded-[8px] p-3.5">
              <Icon name="shield" className="w-4 h-4 shrink-0 mt-0.5 text-[var(--brand)]"/>
              Data responden digunakan untuk keperluan analisis survei dan akan dikelola sesuai ketentuan yang berlaku.
            </p>
          </Card>
        )}

        {step===1 && (
          <div className="anim-rise space-y-4">
            <Card pad="p-6 sm:p-8">
              <SectionHead eyebrow="Langkah 02" title="Penilaian kepuasan"
                desc="Nilai setiap unsur pelayanan sesuai pengalaman Anda. Pertanyaan bertanda wajib harus dijawab."/>
              <div className="mt-6 flex flex-wrap items-center gap-x-5 gap-y-2 text-[12.5px]">
                {skala.map(s=> <LegendDot key={s.v} color={s.warna}>{s.v} — {s.label}</LegendDot>)}
              </div>
              {err.nilai && <p className="mt-5 flex items-center gap-2 text-[13.5px] font-semibold text-[var(--bad)] bg-[color:var(--bad)]/10 border border-[color:var(--bad)]/25 rounded-[8px] px-3.5 py-3">
                <Icon name="alert" className="w-4 h-4 shrink-0"/>{err.nilai}</p>}
              <div className="mt-5 flex items-center gap-3">
                <div className="h-1.5 flex-1 rounded-full bg-[var(--ground)] overflow-hidden">
                  <div className="h-full bg-[var(--gold-bright)] transition-[width] duration-300" style={{width:`${(terisi/Math.max(1,pertanyaan.length))*100}%`}}/>
                </div>
                <span className="font-num text-[12.5px] text-muted">{terisi}/{pertanyaan.length} dinilai</span>
              </div>
            </Card>
            {pertanyaan.map(p=>(
              <Card key={p.id} pad="p-5 sm:p-6" className={f.nilai[p.id]?'!border-[color:var(--brand)]/35':''}>
                <div className="flex items-start gap-3 mb-4">
                  <span className="font-num text-[11.5px] font-semibold text-[var(--brand)] bg-brandsoft rounded-[5px] px-1.5 py-1 shrink-0 mt-0.5">{p.kode}</span>
                  <div>
                    <p className="text-[11px] eyebrow font-bold text-muted mb-1">{p.unsur}{!p.wajib && <span className="normal-case font-normal"> (opsional)</span>}</p>
                    <p className="font-display font-bold text-[15.5px] sm:text-[16.5px] leading-snug">{p.teks}</p>
                  </div>
                </div>
                <RatingScale name={p.unsur} skala={skala} value={f.nilai[p.id]} onChange={v=>setNilai(p.id,v)} invalid={!!err.nilai && p.wajib && !f.nilai[p.id]}/>
              </Card>
            ))}
          </div>
        )}

        {step===2 && (
          <div className="anim-rise space-y-4">
            <Card pad="p-6 sm:p-8">
              <SectionHead eyebrow="Langkah 03" title="Apa yang dapat kami tingkatkan?"
                desc="Masukan tertulis Anda dibaca langsung oleh tim pelayanan dan dikelompokkan berdasarkan unsur pelayanan."/>
              <div className="mt-7 space-y-5">
                <Field label="Saran, kritik, atau masukan Anda" error={err.saran}>
                  <Textarea id="s-saran" value={f.saran} onChange={e=>set('saran',e.target.value.slice(0,600))}
                    placeholder="Contoh: alur pengajuan rekomendasi penelitian sebaiknya dijelaskan dalam satu halaman infografis."/>
                  <span className="mt-1.5 block text-right font-num text-[12px] text-muted">{f.saran.length}/600</span>
                </Field>
                <Field label="Hal yang paling Anda apresiasi dari pelayanan BRIDA Kabupaten Buleleng" error={err.apresiasi}>
                  <Textarea id="s-apre" className="!min-h-[90px]" value={f.apresiasi} onChange={e=>set('apresiasi',e.target.value.slice(0,300))}
                    placeholder="Contoh: petugas menjelaskan persyaratan dengan rinci sejak awal."/>
                </Field>
              </div>
            </Card>

            <Card pad="p-6 sm:p-8">
              <h3 className="font-display text-[19px] font-bold mb-1">Ringkasan jawaban Anda</h3>
              <p className="text-[13.5px] text-muted mb-5">Periksa kembali sebelum mengirim. Anda masih bisa kembali untuk mengubah jawaban.</p>
              <div className="rounded-[9px] border border-line overflow-hidden">
                <div className="bg-[var(--ground)] px-4 py-3 flex items-center justify-between">
                  <span className="eyebrow text-[10.5px] font-bold text-muted">Data responden</span>
                  <button onClick={()=>{setStep(0);naik();}} className="text-[12.5px] font-semibold text-[var(--brand)]">Ubah</button>
                </div>
                <dl className="divide-y divide-[color:var(--line-soft)]">
                  {[['Nama', f.nama||'— tidak diisi —'],['Jenis kelamin',f.jk],['Usia', f.usia+' tahun'],['Pendidikan',f.pendidikan],['Pekerjaan',f.pekerjaan],['Layanan',layananNama],['Kecamatan',f.kecamatan||'— tidak diisi —'],['Frekuensi',f.frekuensi]].map(([k,v])=>(
                    <div key={k} className="flex gap-4 px-4 py-2.5 text-[14px]">
                      <dt className="text-muted w-[130px] shrink-0">{k}</dt><dd className="text-ink2 min-w-0">{v}</dd>
                    </div>
                  ))}
                </dl>
              </div>
              <div className="rounded-[9px] border border-line overflow-hidden mt-4">
                <div className="bg-[var(--ground)] px-4 py-3 flex items-center justify-between">
                  <span className="eyebrow text-[10.5px] font-bold text-muted">Penilaian {pertanyaan.length} unsur</span>
                  <button onClick={()=>{setStep(1);naik();}} className="text-[12.5px] font-semibold text-[var(--brand)]">Ubah</button>
                </div>
                <ul className="divide-y divide-[color:var(--line-soft)]">
                  {pertanyaan.map(p=>{ const v=f.nilai[p.id]; const s=skala.find(x=>x.v===v);
                    return <li key={p.id} className="flex items-center gap-3 px-4 py-2.5">
                      <span className="font-num text-[11.5px] text-muted w-7 shrink-0">{p.kode}</span>
                      <span className="text-[14px] flex-1 min-w-0 truncate">{p.unsur}</span>
                      {s ? <span className="text-[12.5px] font-semibold px-2 py-1 rounded-[5px] shrink-0" style={{background:s.warna, color:s.ink}}>{s.v} · {s.pendek}</span>
                         : <span className={`text-[12.5px] font-semibold shrink-0 ${p.wajib?'text-[var(--bad)]':'text-muted'}`}>{p.wajib?'Belum dinilai':'Dilewati'}</span>}
                    </li>;
                  })}
                </ul>
              </div>

              <label className="mt-6 flex gap-3 items-start cursor-pointer">
                <input id="s-setuju" type="checkbox" checked={f.setuju} onChange={e=>set('setuju',e.target.checked)}
                  className="mt-0.5 w-[18px] h-[18px] accent-[var(--brand)] shrink-0"/>
                <span className="text-[14.5px] text-ink2 leading-relaxed">Saya bersedia data jawaban saya digunakan untuk keperluan evaluasi pelayanan.</span>
              </label>
              {err.setuju && <p className="mt-2 flex items-center gap-2 text-[13px] font-semibold text-[var(--bad)]"><Icon name="alert" className="w-4 h-4"/>{err.setuju}</p>}
              {/* Honeypot: tersembunyi dari pengguna, hanya diisi bot. */}
              <input type="text" name="website" tabIndex={-1} autoComplete="off" aria-hidden="true" value={f.website} onChange={e=>set('website',e.target.value)}
                style={{position:'absolute', left:'-10000px', width:1, height:1, overflow:'hidden'}}/>
            </Card>
          </div>
        )}

        {/* Navigasi langkah */}
        <div className="mt-6 flex flex-col-reverse sm:flex-row gap-3 sm:justify-between">
          {step>0
            ? <Button variant="outline" size="lg" icon="arrowLeft" onClick={kembali} disabled={mengirim}>Kembali</Button>
            : <Button variant="quiet" size="lg" icon="arrowLeft" onClick={()=>navigate('/')}>Batalkan</Button>}
          {step<2
            ? <Button size="lg" iconRight="arrowRight" onClick={lanjut}>Lanjutkan</Button>
            : <Button size="lg" icon="send" onClick={kirim} disabled={mengirim}>{mengirim?'Mengirim…':'Kirim Survei'}</Button>}
        </div>
      </div>
    </section>
  </Shell>;
}

export default Survei;
