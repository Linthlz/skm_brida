import { useState } from 'react';
import { Button, Card, SectionHead, Field, Input, Textarea, Select } from '../components/ui/index.js';
import Icon from '../components/Icon.jsx';
import { useToast } from '../hooks/useToast.jsx';
import { useApi } from '../hooks/useApi.js';
import { kirimKontak, referensi } from '../services/api/publik.js';

const INFO = [
  {i:'mapPin', j:'Kantor', d:'Jl. Ngurah Rai No. 2, Singaraja\nKabupaten Buleleng, Bali 81116'},
  {i:'clock',  j:'Jam layanan', d:'Senin – Kamis  08.00 – 15.30 WITA\nJumat  08.00 – 14.00 WITA'},
  {i:'phone',  j:'Telepon', d:'(0362) 000 000'},
  {i:'mail',   j:'Surel', d:'brida@bulelengkab.go.id'},
];

const KOSONG = {nama:'',email:'',topik:'',pesan:'',website:''};

function Kontak() {
  const toast = useToast();
  const ref = useApi(referensi);
  const [f,setF] = useState(KOSONG);
  const [err,setErr] = useState({});
  const [kirimMasuk,setKirimMasuk] = useState(false);
  const set=(k,v)=>setF(p=>({...p,[k]:v}));
  async function kirim(e){
    e.preventDefault();
    // Pemeriksaan cepat di browser hanya untuk kenyamanan; backend tetap memvalidasi ulang.
    const er={};
    if(!f.nama.trim()) er.nama='Nama wajib diisi.';
    if(!/^\S+@\S+\.\S+$/.test(f.email)) er.email='Masukkan alamat email yang valid, contoh: nama@email.com';
    if(!f.topik) er.topik='Silakan pilih salah satu topik.';
    if(f.pesan.trim().length<10) er.pesan='Tuliskan pesan minimal 10 karakter agar kami dapat menindaklanjuti.';
    setErr(er);
    if(Object.keys(er).length) return;
    setKirimMasuk(true);
    try{
      const r = await kirimKontak(f);
      setF(KOSONG); toast(r.message);
    }catch(ex){
      if(ex.status===422) setErr({nama:ex.field('nama'), email:ex.field('email'), topik:ex.field('topik'), pesan:ex.field('pesan')});
      toast(ex.status===422 ? 'Periksa kembali isian formulir.' : ex.message, 'galat');
    }finally{ setKirimMasuk(false); }
  }
  return (
    <section id="kontak" className="section-anchor py-14 sm:py-16 border-t border-linesoft">
      <SectionHead eyebrow="Kontak" title="Hubungi BRIDA Kabupaten Buleleng"
        desc="Untuk pertanyaan seputar survei, permintaan data hasil SKM, atau pengaduan layanan yang memerlukan tindak lanjut."/>
      <div className="grid lg:grid-cols-[1fr_1.1fr] gap-8 mt-9 items-start">
        <div className="space-y-4">
          {INFO.map(c=>(
            <Card key={c.j} pad="p-5" className="flex gap-4">
              <span className="shrink-0 grid place-items-center w-10 h-10 rounded-[8px] bg-brandsoft text-[var(--brand)]"><Icon name={c.i} className="w-[19px] h-[19px]"/></span>
              <div>
                <p className="font-display font-bold text-[15.5px] mb-1">{c.j}</p>
                <p className="text-[14px] text-ink2 leading-relaxed whitespace-pre-line select-text">{c.d}</p>
              </div>
            </Card>
          ))}
          <Card pad="p-5" className="!bg-goldsoft !border-[color:var(--gold)]/30 flex gap-3">
            <Icon name="info" className="w-5 h-5 text-[var(--gold-bright)] shrink-0 mt-0.5"/>
            <p className="text-[13.5px] text-ink2 leading-relaxed">Alamat dan nomor di atas adalah contoh untuk prototipe. Ganti dengan data resmi sebelum publikasi.</p>
          </Card>
        </div>

        <Card pad="p-6 sm:p-8">
          <h3 className="font-display text-[19px] font-bold mb-1">Kirim pesan</h3>
          <p className="text-[13.5px] text-muted mb-6">Kami membalas paling lambat dua hari kerja.</p>
          <form onSubmit={kirim} className="space-y-5">
            <Field label="Nama" required error={err.nama}><Input id="k-nama" value={f.nama} onChange={e=>set('nama',e.target.value)} placeholder="Nama lengkap Anda"/></Field>
            <Field label="Email" required error={err.email}><Input id="k-mail" type="email" value={f.email} onChange={e=>set('email',e.target.value)} placeholder="nama@email.com"/></Field>
            <Field label="Topik" required error={err.topik}>
              <Select id="k-topik" options={ref.data?.topik_kontak || []}
                value={f.topik} onChange={e=>set('topik',e.target.value)}/>
            </Field>
            <Field label="Pesan" required error={err.pesan}><Textarea id="k-pesan" maxLength={2000} value={f.pesan} onChange={e=>set('pesan',e.target.value)} placeholder="Tuliskan pertanyaan atau masukan Anda"/></Field>
            {/* Honeypot: tersembunyi dari pengguna, hanya diisi bot. */}
            <input type="text" name="website" tabIndex={-1} autoComplete="off" aria-hidden="true" value={f.website} onChange={e=>set('website',e.target.value)}
              style={{position:'absolute', left:'-10000px', width:1, height:1, overflow:'hidden'}}/>
            <Button type="submit" size="lg" icon="send" className="w-full sm:w-auto" disabled={kirimMasuk}>{kirimMasuk?'Mengirim…':'Kirim Pesan'}</Button>
          </form>
        </Card>
      </div>
    </section>
  );
}

export default Kontak;
