import Button from './Button.jsx';
import { fmtInt } from '../../utils/format.js';

/** Navigasi halaman untuk daftar berpaginasi dari API (meta Laravel). */
function Paginasi({meta, onPage, disabled}){
  if(!meta || meta.last_page<=1) return null;
  return <div className="flex items-center justify-between gap-3">
    <p className="text-[13px] text-muted">Halaman <span className="font-num text-ink">{meta.current_page}</span> dari <span className="font-num">{meta.last_page}</span> · <span className="font-num">{fmtInt(meta.total)}</span> data</p>
    <div className="flex gap-2">
      <Button size="sm" variant="outline" icon="chevronLeft" disabled={disabled||meta.current_page<=1} onClick={()=>onPage(meta.current_page-1)}>Sebelumnya</Button>
      <Button size="sm" variant="outline" iconRight="chevronRight" disabled={disabled||meta.current_page>=meta.last_page} onClick={()=>onPage(meta.current_page+1)}>Berikutnya</Button>
    </div>
  </div>;
}

export default Paginasi;
