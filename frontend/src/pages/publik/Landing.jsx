import { Fragment } from 'react';
import Shell from '../../layouts/Shell.jsx';
import Hero from '../../sections/Hero.jsx';
import Sekilas from '../../sections/Sekilas.jsx';
import Tentang from '../../sections/Tentang.jsx';
import Hasil from '../../sections/Hasil.jsx';
import Kontak from '../../sections/Kontak.jsx';
import Ajakan from '../../sections/Ajakan.jsx';

/**
 * Seluruh isi situs publik dijadikan satu halaman. Bagian Tentang, Hasil, dan
 * Kontak punya id sehingga bisa dituju lewat tautan jangkar (#tentang dst.);
 * jarak jangkarnya diatur kelas `.section-anchor` di src/index.css agar tidak
 * tertutup navbar yang sticky. Isi Survei tetap halaman tersendiri di /survei.
 */
function Landing() {
  return (
    <Fragment>
      <Hero />
      <Shell>
        <Sekilas />
        <Tentang />
        <Hasil />
        <Kontak />
        <Ajakan />
      </Shell>
    </Fragment>
  );
}

export default Landing;
