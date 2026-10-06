import { useState } from 'react';
import { Navigate, useLocation, useNavigate } from 'react-router-dom';
import { Button, Field, Input } from '../../components/ui/index.js';
import Icon from '../../components/Icon.jsx';
import Logo from '../../components/Logo.jsx';
import { useToast } from '../../hooks/useToast.jsx';
import { useAuth } from '../../hooks/useAuth.jsx';
import { useApi } from '../../hooks/useApi.js';
import { ringkasan } from '../../services/api/publik.js';
import { fmt, fmtInt } from '../../utils/format.js';

function Login() {
  const navigate = useNavigate();
  const { state } = useLocation();
  const toast = useToast();
  const { status, login } = useAuth();
  const ring = useApi(ringkasan);
  const [f,setF] = useState({email:'', pass:'', ingat:false});
  const [err,setErr] = useState({});
  const [masuk,setMasuk] = useState(false);

  // Sudah punya sesi aktif: langsung ke panel.
  if (status === 'masuk') return <Navigate to={state?.dari || '/admin'} replace />;

  async function kirim(e){
    e.preventDefault();
    const er={};
    if(!f.email.trim()) er.email='Email wajib diisi.';
    if(!f.pass) er.pass='Kata sandi wajib diisi.';
    setErr(er);
    if(Object.keys(er).length) return;
    setMasuk(true);
    try{
      const u = await login({email:f.email.trim(), password:f.pass, ingat:f.ingat});
      toast(`Berhasil masuk sebagai ${u.nama}`);
      navigate(state?.dari || '/admin', {replace:true});
    }catch(ex){
      // 422 = kredensial salah, 429 = terlalu banyak percobaan; pesan dari server sudah ramah pengguna.
      if(ex.status===422 || ex.status===429) setErr({email: ex.field('email') || ex.message, pass: ex.field('password')});
      else toast(ex.message, 'galat');
      setF(p=>({...p, pass:''}));
    }finally{ setMasuk(false); }
  }

  const r = ring.data;
  return <div className="min-h-[100svh] grid lg:grid-cols-2">
    <div className="hidden lg:flex flex-col justify-between p-12" style={{background:'var(--brand-deep)'}}>
      <div className="flex items-center gap-3"><Logo className="w-11 h-11"/>
        <span className="font-display font-extrabold text-[16px] text-white leading-tight">BRIDA Kabupaten Buleleng<br/><span className="text-white/60 text-[13px] font-semibold">Panel Administrator SKM</span></span>
      </div>
      <div>
        <p className="font-display text-[30px] font-extrabold text-white leading-[1.15]" style={{maxWidth:'18ch'}}>Kelola dan analisis hasil survei dalam satu tempat.</p>
        {r && <div className="mt-8 grid grid-cols-3 gap-6">
          {[['Responden',fmtInt(r.total_responden)],['Indeks',fmt(r.ikm)],['Feedback', fmtInt(r.total_feedback)]].map(([k,v])=>(
            <div key={k}><p className="font-num text-[24px] font-semibold text-[var(--gold-bright)] leading-none">{v}</p><p className="text-[12px] text-white/55 mt-2">{k}</p></div>
          ))}
        </div>}
      </div>
      <p className="text-[12.5px] text-white/45">{r ? `Periode survei ${r.periode}` : 'Survei Kepuasan Masyarakat'}</p>
    </div>
    <div className="flex items-center justify-center p-6 sm:p-10 bg-surface">
      <div className="w-full" style={{maxWidth:'400px'}}>
        <div className="lg:hidden flex items-center gap-3 mb-8"><Logo className="w-10 h-10"/>
          <span className="font-display font-extrabold text-[15px] leading-tight">BRIDA Kabupaten Buleleng<br/><span className="text-muted text-[12px]">Panel Administrator</span></span></div>
        <h1 className="font-display text-[26px] font-extrabold">Masuk ke panel admin</h1>
        <p className="text-[14.5px] text-muted mt-2">Gunakan akun resmi yang terdaftar pada BRIDA Kabupaten Buleleng.</p>
        <form onSubmit={kirim} className="mt-7 space-y-5" noValidate>
          <Field label="Email dinas" required error={err.email}><Input id="a-user" type="email" autoComplete="username" value={f.email} onChange={e=>setF({...f,email:e.target.value})} placeholder="nama@bulelengkab.go.id"/></Field>
          <Field label="Kata sandi" required error={err.pass}>
            <Input id="a-pass" type="password" autoComplete="current-password" value={f.pass} onChange={e=>setF({...f,pass:e.target.value})} placeholder="Masukkan kata sandi"/>
          </Field>
          <div className="flex items-center justify-between">
            <label className="flex items-center gap-2.5 cursor-pointer text-[14px] text-ink2">
              <input type="checkbox" checked={f.ingat} onChange={e=>setF({...f,ingat:e.target.checked})} className="w-[17px] h-[17px] accent-[var(--brand)]"/> Ingat perangkat ini
            </label>
            <button type="button" onClick={()=>toast('Hubungi administrator sistem BRIDA untuk mengatur ulang kata sandi.')} className="text-[13.5px] font-semibold text-[var(--brand)]">Lupa sandi?</button>
          </div>
          <Button type="submit" size="lg" className="w-full" icon="logOut" disabled={masuk}>{masuk?'Memeriksa…':'Masuk'}</Button>
        </form>
        <button onClick={()=>navigate('/')} className="mt-6 inline-flex items-center gap-2 text-[14px] font-semibold text-muted hover:text-ink">
          <Icon name="arrowLeft" className="w-4 h-4"/> Kembali ke situs publik
        </button>
      </div>
    </div>
  </div>;
}

export default Login;
