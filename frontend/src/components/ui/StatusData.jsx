import Button from './Button.jsx';
import EmptyState from './EmptyState.jsx';

/** Blok kerangka abu-abu selama data dimuat. */
function Memuat({ tinggi = 220, baris = 0, className = '' }) {
  if (baris) {
    return <div className={`space-y-2.5 ${className}`} aria-busy="true" aria-label="Memuat data">
      {Array.from({ length: baris }, (_, i) => <div key={i} className="h-4 rounded-[6px] bg-[var(--ground)] animate-pulse" style={{ width: `${92 - (i % 3) * 14}%` }}/>)}
    </div>;
  }
  return <div className={`rounded-[9px] bg-[var(--ground)] animate-pulse ${className}`} style={{ height: tinggi }} aria-busy="true" aria-label="Memuat data"/>;
}

/** Pesan galat dengan tombol muat ulang. Menyesuaikan teks dengan kode status API. */
function GagalMuat({ error, onRetry, judul }) {
  const s = error?.status;
  const icon = s === 0 ? 'refresh' : s === 403 ? 'lock' : 'alert';
  const t = judul || (s === 0 ? 'Tidak dapat terhubung ke server' : s === 403 ? 'Akses tidak tersedia' : s === 404 ? 'Data tidak ditemukan' : 'Data gagal dimuat');
  return <EmptyState icon={icon} title={t} desc={error?.message || 'Terjadi kesalahan. Silakan coba lagi.'}
    action={onRetry && s !== 403 && <Button size="sm" variant="outline" icon="refresh" onClick={onRetry}>Coba lagi</Button>}/>;
}

export { Memuat, GagalMuat };
