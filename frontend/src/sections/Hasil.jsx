import { useState } from 'react';
import { Badge, Card, EmptyState, GagalMuat, Memuat, SectionHead, StatCard, Select } from '../components/ui/index.js';
import { ChartCard, BarUnsur, DonutKepuasan, TrenIKM } from '../components/charts/index.js';
import Icon from '../components/Icon.jsx';
import { useApi } from '../hooks/useApi.js';
import { daftarPeriode, hasil, referensi } from '../services/api/publik.js';
import { distribusiBerwarna } from '../data/skala.js';
import { fmt, fmtInt, unsurGrafik } from '../utils/format.js';
import { mutu } from '../utils/mutu.js';

function Hasil() {
  const [periode,setPeriode] = useState('');
  const [layanan,setLayanan] = useState('');
  const per = useApi(daftarPeriode);
  const ref = useApi(referensi);
  const h = useApi(()=>hasil({periode, layanan_id:layanan}), [periode, layanan]);

  // Periode kosong = periode aktif di server; tampilkan tahunnya di dropdown.
  const aktif = per.data?.find(p=>p.aktif)?.tahun;
  const opsiPeriode = (per.data||[]).map(p=>String(p.tahun));
  const opsiLayanan = (ref.data?.layanan||[]).map(l=>({value:String(l.id), label:l.nama}));

  return (
    <section id="hasil" className="section-anchor py-14 sm:py-16 border-t border-linesoft">
      <div className="flex flex-col lg:flex-row lg:items-end justify-between gap-5">
        <SectionHead eyebrow="Transparansi hasil" title="Hasil Survei Kepuasan Masyarakat"
          desc="Rekapitulasi penilaian pengguna layanan BRIDA Kabupaten Buleleng."/>
        {h.error?.status !== 403 && <div className="flex flex-wrap gap-2.5">
          <div className="min-w-[140px]"><Select aria-label="Periode" options={opsiPeriode} placeholder={null} value={periode || String(aktif ?? '')} onChange={e=>setPeriode(e.target.value)}/></div>
          <div className="min-w-[200px]"><Select aria-label="Jenis layanan" options={opsiLayanan} placeholder="Semua jenis layanan" value={layanan} onChange={e=>setLayanan(e.target.value)}/></div>
        </div>}
      </div>

      <IsiHasil h={h}/>
    </section>
  );
}

function IsiHasil({h}) {
  if (h.error) {
    return <Card className="mt-8">{h.error.status === 403
      ? <EmptyState icon="lock" title="Hasil survei belum dipublikasikan" desc="Rekapitulasi akan tampil di sini setelah dipublikasikan oleh BRIDA Kabupaten Buleleng."/>
      : <GagalMuat error={h.error} onRetry={h.reload}/>}</Card>;
  }
  if (!h.data) {
    return <div className="mt-8 space-y-4">
      <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">{[0,1,2,3].map(i=><Memuat key={i} tinggi={112}/>)}</div>
      <div className="grid lg:grid-cols-2 gap-4"><Memuat tinggi={380}/><Memuat tinggi={380}/></div>
    </div>;
  }

  const d = h.data;
  const m = mutu(d.ikm ?? 0);
  const unsur = unsurGrafik(d.unsur);
  const urut = [...unsur].sort((a,b)=>b.nilai-a.nilai);
  const tertinggi = urut[0], terendah = urut[urut.length-1];

  return <div className={`transition-opacity ${h.loading?'opacity-60':''}`} aria-busy={h.loading}>
    <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4 mt-8">
      <StatCard label="Indeks Kepuasan Masyarakat" value={fmt(d.ikm)} sub={d.mutu && `Mutu ${d.mutu.huruf}`} icon="gauge"/>
      <StatCard label="Jumlah responden" value={fmtInt(d.total_responden)} sub="orang" icon="users"/>
      <StatCard label="Periode survei" value={String(d.periode)} sub={d.rentang_bulan || undefined} icon="calendar"/>
      <StatCard label="Tingkat kepuasan" value={d.mutu?.label || '–'} icon="award" tone="gold"/>
    </div>

    {unsur.length === 0
      ? <Card className="mt-4"><EmptyState icon="barChart" title="Belum ada data pada pilihan ini" desc="Belum ada responden untuk periode atau jenis layanan yang dipilih."/></Card>
      : <>
      <div className="grid lg:grid-cols-2 gap-4 mt-4 items-start">
        <ChartCard title="Tingkat kepuasan per unsur pelayanan"
          desc={`Nilai rata-rata tiap unsur pada skala 1–${d.skala_maks}. Unsur dengan nilai terendah menjadi prioritas perbaikan tahun berjalan.`}>
          <BarUnsur data={unsur} max={d.skala_maks} pengali={d.pengali}/>
          <div className="mt-5 pt-5 border-t border-linesoft grid sm:grid-cols-2 gap-4">
            <div className="flex gap-3">
              <span className="shrink-0 w-8 h-8 grid place-items-center rounded-[7px] bg-[color:var(--ok)]/12 text-[var(--ok)]"><Icon name="star" className="w-4 h-4"/></span>
              <div><p className="text-[12px] text-muted">Nilai tertinggi</p><p className="text-[14px] font-semibold">{tertinggi.nama} <span className="font-num text-muted">({fmt(tertinggi.nilai)})</span></p></div>
            </div>
            <div className="flex gap-3">
              <span className="shrink-0 w-8 h-8 grid place-items-center rounded-[7px] bg-[color:var(--warn)]/12 text-[var(--warn)]"><Icon name="alert" className="w-4 h-4"/></span>
              <div><p className="text-[12px] text-muted">Prioritas perbaikan</p><p className="text-[14px] font-semibold">{terendah.nama} <span className="font-num text-muted">({fmt(terendah.nilai)})</span></p></div>
            </div>
          </div>
        </ChartCard>

        <div className="space-y-4 min-w-0">
          <ChartCard title="Distribusi tingkat kepuasan responden"
            desc="Sebaran kategori kepuasan berdasarkan nilai rata-rata tiap responden.">
            <DonutKepuasan data={distribusiBerwarna(d.distribusi)}/>
          </ChartCard>
          {d.tren.length > 0 && <ChartCard title="Tren Indeks Kepuasan Masyarakat"
            desc="Perubahan nilai IKM antarperiode survei. Titik emas menandai periode terakhir.">
            <TrenIKM data={d.tren}/>
          </ChartCard>}
        </div>
      </div>

      <Card className="mt-4" pad="p-0">
        <div className="px-5 sm:px-6 py-5 border-b border-linesoft">
          <h3 className="font-display text-[17px] font-bold">Rekapitulasi nilai per unsur</h3>
          <p className="text-[13.5px] text-muted mt-1">Nilai rata-rata (NRR) dikonversi ke indeks dengan pengali {fmt(d.pengali,0)} karena survei ini memakai skala 1–{d.skala_maks}.</p>
        </div>
        <div className="scroll-x">
          <table className="w-full text-[14px]" style={{minWidth:'620px'}}>
            <thead>
              <tr className="text-left text-[12px] eyebrow font-bold text-muted bg-[var(--ground)]">
                <th className="px-5 py-3 font-bold">Kode</th><th className="px-5 py-3 font-bold">Unsur pelayanan</th>
                <th className="px-5 py-3 font-bold text-right">NRR</th><th className="px-5 py-3 font-bold text-right">Indeks</th><th className="px-5 py-3 font-bold">Mutu</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[color:var(--line-soft)]">
              {d.unsur.map(u=>{ const mm=mutu(u.indeks);
                return <tr key={u.kode} className="hover:bg-[var(--ground)]">
                  <td className="px-5 py-3 font-num text-muted">{u.kode}</td>
                  <td className="px-5 py-3">{u.nama}</td>
                  <td className="px-5 py-3 font-num text-right">{fmt(u.nrr)}</td>
                  <td className="px-5 py-3 font-num text-right font-semibold">{fmt(u.indeks)}</td>
                  <td className="px-5 py-3"><Badge tone={mm.tone}>{mm.huruf} · {mm.label}</Badge></td>
                </tr>;
              })}
              <tr className="bg-brandsoft font-semibold">
                <td className="px-5 py-3.5" colSpan="2">Rata-rata seluruh unsur</td>
                <td className="px-5 py-3.5 font-num text-right">{fmt(d.nrr)}</td>
                <td className="px-5 py-3.5 font-num text-right">{fmt(d.ikm)}</td>
                <td className="px-5 py-3.5"><Badge tone={m.tone}>{m.huruf} · {m.label}</Badge></td>
              </tr>
            </tbody>
          </table>
        </div>
      </Card>
    </>}
  </div>;
}

export default Hasil;
