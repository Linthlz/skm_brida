import { Card, GagalMuat, Memuat, SectionHead } from '../components/ui/index.js';
import Icon from '../components/Icon.jsx';
import { useApi } from '../hooks/useApi.js';
import { formSurvei } from '../services/api/publik.js';

const RINGKAS = [
  {i:'target', j:'Tujuan Survei', d:'Mengetahui tingkat kepuasan masyarakat terhadap mutu pelayanan BRIDA Kabupaten Buleleng.'},
  {i:'clock',  j:'Waktu Pengisian', d:'Sekitar 3–5 menit. Jawaban tersimpan di perangkat Anda selama pengisian berlangsung.'},
  {i:'users',  j:'Responden', d:'Masyarakat dan pengguna layanan BRIDA — peneliti, mahasiswa, OPD, pelaku usaha, dan umum.'},
  {i:'lightbulb', j:'Manfaat', d:'Menjadi bahan evaluasi tahunan dan dasar perbaikan standar pelayanan publik.'},
];

function Sekilas() {
  const form = useApi(formSurvei);
  const unsur = form.data?.pertanyaan || [];
  const maks = form.data?.periode.skala_maks ?? 5;
  return (
    <>
      <section className="py-14 sm:py-16">
        <SectionHead eyebrow="Sekilas survei" title="Yang perlu Anda ketahui sebelum mengisi"
          desc="Survei ini terbuka untuk siapa saja yang pernah menggunakan layanan BRIDA Kabupaten Buleleng, baik secara langsung di kantor maupun secara daring."/>
        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4 mt-8">
          {RINGKAS.map(c=>(
            <Card key={c.j} className="hover:border-[color:var(--brand)] transition-colors">
              <span className="grid place-items-center w-10 h-10 rounded-[8px] bg-brandsoft text-[var(--brand)] mb-4"><Icon name={c.i} className="w-[19px] h-[19px]"/></span>
              <h3 className="font-display font-bold text-[16px] mb-2">{c.j}</h3>
              <p className="text-[14px] text-ink2 leading-relaxed">{c.d}</p>
            </Card>
          ))}
        </div>
      </section>

      <section className="pb-14 sm:pb-16">
        <div className="bg-surface border border-line rounded-card overflow-hidden">
          <div className="grid lg:grid-cols-[.9fr_1.1fr]">
            <div className="p-7 sm:p-9 border-b lg:border-b-0 lg:border-r border-line">
              <SectionHead eyebrow={`${unsur.length || 9} unsur pelayanan`} title="Apa saja yang Anda nilai?"
                desc={`Unsur penilaian mengikuti Permen PANRB 14/2017. Setiap unsur dinilai pada skala 1 sampai ${maks}, dan hasilnya dikonversi menjadi Indeks Kepuasan Masyarakat.`}/>
            </div>
            {form.error ? <div className="p-4"><GagalMuat error={form.error} onRetry={form.reload}/></div>
            : !form.data ? <div className="p-6"><Memuat baris={9}/></div>
            : <ul className="p-4 sm:p-6 grid sm:grid-cols-2 gap-x-6">
              {unsur.map(u=>(
                <li key={u.id} className="flex items-center gap-3 py-2.5 border-b border-linesoft last:border-0">
                  <span className="font-num text-[11.5px] font-semibold text-[var(--brand)] bg-brandsoft rounded-[5px] px-1.5 py-1 shrink-0">{u.kode}</span>
                  <span className="text-[14px] text-ink2">{u.unsur}</span>
                </li>
              ))}
            </ul>}
          </div>
        </div>
      </section>
    </>
  );
}

export default Sekilas;
