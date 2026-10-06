import { useNavigate } from 'react-router-dom';
import { Button } from '../components/ui/index.js';

function Ajakan() {
  const navigate = useNavigate();
  return (
    <section className="pb-4">
      <div className="rounded-card p-8 sm:p-10 text-center" style={{background:'var(--brand-deep)'}}>
        <h2 className="font-display text-[24px] sm:text-[30px] font-extrabold text-white leading-tight" style={{maxWidth:'26ch', margin:'0 auto'}}>
          Lima menit dari Anda, satu tahun perbaikan bagi kami
        </h2>
        <p className="mt-3 text-[15.5px] text-white/75 mx-auto" style={{maxWidth:'56ch'}}>
          Setiap jawaban direkap menjadi nilai per unsur dan dibahas dalam evaluasi mutu pelayanan BRIDA Kabupaten Buleleng.
        </p>
        <div className="mt-7 flex flex-col sm:flex-row gap-3 justify-center">
          <Button size="lg" variant="gold" icon="clipboard" onClick={()=>navigate('/survei')}>Isi Survei</Button>
          <Button as="a" href="#hasil" size="lg" variant="outline" className="!border-white/30 !text-white hover:!bg-white/10">Lihat Hasil Survei</Button>
        </div>
      </div>
    </section>
  );
}

export default Ajakan;
