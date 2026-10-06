import { Card, GagalMuat, Memuat, SectionHead } from '../components/ui/index.js';
import { MeterIKM } from '../components/charts/index.js';
import Icon from '../components/Icon.jsx';
import { useApi } from '../hooks/useApi.js';
import { ringkasan } from '../services/api/publik.js';

const BLOK = [
  {i:'target', j:'Tujuan survei', d:'Mengukur tingkat kepuasan pengguna layanan, mengidentifikasi unsur pelayanan yang paling perlu diperbaiki, dan menyediakan data pembanding antarperiode.'},
  {i:'users', j:'Siapa yang dapat mengisi', d:'Setiap orang yang pernah menerima layanan BRIDA Kabupaten Buleleng dalam satu tahun terakhir. Satu pengguna cukup mengisi satu kali untuk setiap jenis layanan yang digunakan.'},
  {i:'clock', j:'Estimasi waktu pengisian', d:'Sekitar 3–5 menit untuk 9 pertanyaan penilaian, data responden singkat, dan satu kolom saran.'},
  {i:'shield', j:'Kerahasiaan data responden', d:'Nama dan kontak bersifat opsional. Data responden hanya diolah secara agregat untuk keperluan analisis dan tidak dipublikasikan per individu.'},
  {i:'trendingUp', j:'Pemanfaatan hasil survei', d:'Nilai per unsur dilaporkan dalam laporan SKM semesteran, dipublikasikan pada bagian Hasil Survei, dan menjadi bahan rapat evaluasi mutu pelayanan.'},
];

const MUTU = [
  ['A','88,31 – 100','Sangat Baik','var(--s5)','var(--s5i)'],
  ['B','76,61 – 88,30','Baik','var(--s4)','var(--s4i)'],
  ['C','65,00 – 76,60','Kurang Baik','var(--s2)','var(--s2i)'],
  ['D','25,00 – 64,99','Tidak Baik','var(--s1)','var(--s1i)'],
];

function Tentang() {
  const ring = useApi(ringkasan);
  return (
    <section id="tentang" className="section-anchor py-14 sm:py-16 border-t border-linesoft">
      <SectionHead eyebrow="Tentang survei" title="Apa itu Survei Kepuasan Masyarakat?"
        desc="Survei Kepuasan Masyarakat (SKM) adalah instrumen untuk memperoleh masukan masyarakat terhadap kualitas pelayanan publik. Hasilnya menjadi ukuran kinerja pelayanan sekaligus dasar perbaikan bagi unit penyelenggara."/>
      <div className="grid lg:grid-cols-[1.15fr_.85fr] gap-8 lg:gap-12 mt-10 items-start">
        <div className="space-y-7">
          {BLOK.map(b=>(
            <div key={b.j} className="flex gap-4">
              <span className="shrink-0 grid place-items-center w-10 h-10 rounded-[8px] bg-brandsoft text-[var(--brand)]"><Icon name={b.i} className="w-[19px] h-[19px]"/></span>
              <div>
                <h3 className="font-display font-bold text-[17px] mb-1.5">{b.j}</h3>
                <p className="text-[15px] text-ink2 leading-relaxed" style={{maxWidth:'60ch'}}>{b.d}</p>
              </div>
            </div>
          ))}
        </div>

        <div className="space-y-4 lg:sticky lg:top-[100px]">
          {/* Disembunyikan bila hasil belum dipublikasikan (API menjawab 403). */}
          {ring.error?.status !== 403 && <Card pad="p-6">
            {ring.error ? <GagalMuat error={ring.error} onRetry={ring.reload}/>
            : !ring.data ? <Memuat tinggi={110}/>
            : <>
              <p className="eyebrow text-[10.5px] font-bold text-muted mb-4">Ringkasan periode {ring.data.periode}</p>
              {ring.data.ikm !== null
                ? <MeterIKM nilai={ring.data.ikm}/>
                : <p className="text-[14px] text-muted">Belum ada penilaian yang masuk pada periode ini.</p>}
            </>}
          </Card>}
          <Card pad="p-6">
            <p className="eyebrow text-[10.5px] font-bold text-muted mb-4">Kategori mutu pelayanan</p>
            <ul className="space-y-2.5">
              {MUTU.map(([h,r,l,c,ci])=>(
                <li key={h} className="flex items-center gap-3">
                  <span className="w-7 h-7 grid place-items-center rounded-[6px] font-num text-[13px] font-semibold shrink-0" style={{background:c, color:ci}}>{h}</span>
                  <span className="font-num text-[13px] text-ink2 w-[104px]">{r}</span>
                  <span className="text-[13.5px]">{l}</span>
                </li>
              ))}
            </ul>
          </Card>
          <Card pad="p-6" className="!bg-goldsoft !border-[color:var(--gold)]/30">
            <div className="flex gap-3">
              <Icon name="info" className="w-5 h-5 text-[var(--gold-bright)] shrink-0 mt-0.5"/>
              <p className="text-[14px] text-ink2 leading-relaxed">
                Survei ini bukan kanal pengaduan resmi. Untuk pengaduan layanan yang memerlukan tindak lanjut cepat, gunakan formulir pada bagian <a href="#kontak" className="font-semibold text-[var(--brand)] underline underline-offset-2">Kontak</a>.
              </p>
            </div>
          </Card>
        </div>
      </div>
    </section>
  );
}

export default Tentang;
