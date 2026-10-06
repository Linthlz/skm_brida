import { useEffect, useMemo, useRef, useState } from 'react';
import Icon from '../Icon.jsx';
import Button from './Button.jsx';
import EmptyState from './EmptyState.jsx';
import { inputCls } from './Form.jsx';
import { Memuat, GagalMuat } from './StatusData.jsx';
import { fmtInt } from '../../utils/format.js';

/**
 * Tabel dengan pencarian, pengurutan, paginasi, dan tampilan kartu di layar kecil.
 *
 * Dua mode:
 * - klien: kirim `rows`, semua diolah di browser.
 * - server: kirim `fetchPage(params)` yang mengembalikan { data, meta } dari API
 *   berpaginasi Laravel; `params` (filter tambahan) ikut dikirim. Kolom dengan
 *   `sortable:false` tidak bisa diurutkan, `key` kolom dipakai sebagai nama sort.
 */
function DataTable({columns, rows, fetchPage, params, pageSize=8, empty, onRow, defaultSort}){
  const server = !!fetchPage;
  const [q,setQ] = useState('');
  const [cari,setCari] = useState('');
  const [sort,setSort] = useState(defaultSort || {key:columns[0].key, dir:'asc'});

  // Pencarian ke server ditunda sebentar agar tidak memanggil API setiap ketukan.
  useEffect(()=>{ const t=setTimeout(()=>setCari(q.trim()), server?350:0); return ()=>clearTimeout(t); },[q,server]);

  // Halaman kembali ke 1 setiap kali kata kunci, urutan, atau filter berubah.
  const kunci = JSON.stringify([cari, sort, params||{}]);
  const [hal,setHal] = useState({k:kunci, n:1});
  const page = hal.k===kunci ? hal.n : 1;
  const keHalaman = n => setHal({k:kunci, n});

  // ── mode server ──
  const fetchRef = useRef(fetchPage); fetchRef.current = fetchPage;
  const [srv,setSrv] = useState({rows:[], total:0, pages:1, loading:true, error:null});
  const [ulang,setUlang] = useState(0);
  useEffect(()=>{
    if(!server) return;
    let batal=false;
    setSrv(s=>({...s, loading:true, error:null}));
    fetchRef.current({...(params||{}), search:cari||undefined, sort:sort.key, dir:sort.dir, page, per_page:pageSize}).then(
      r=>{ if(!batal) setSrv({rows:r.data, total:r.meta.total, pages:Math.max(1,r.meta.last_page), loading:false, error:null}); },
      e=>{ if(!batal) setSrv(s=>({...s, loading:false, error:e})); });
    return ()=>{ batal=true; };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  },[server, kunci, page, pageSize, ulang]);

  // ── mode klien ──
  const filtered = useMemo(()=>{
    if(server) return [];
    const s=cari.toLowerCase();
    let r = !s ? rows : rows.filter(x=> columns.some(c=> String(c.raw?c.raw(x):x[c.key]).toLowerCase().includes(s)));
    const col = columns.find(c=>c.key===sort.key);
    return [...r].sort((a,b)=>{
      const va = col?.sortVal? col.sortVal(a) : a[sort.key], vb = col?.sortVal? col.sortVal(b) : b[sort.key];
      const cmp = typeof va==='number' ? va-vb : String(va).localeCompare(String(vb),'id');
      return sort.dir==='asc'?cmp:-cmp;
    });
  },[server,cari,rows,sort,columns]);

  const total = server ? srv.total : filtered.length;
  const pages = server ? srv.pages : Math.max(1, Math.ceil(filtered.length/pageSize));
  const view = server ? srv.rows : filtered.slice((page-1)*pageSize, page*pageSize);
  const memuatAwal = server && srv.loading && srv.rows.length===0;

  function urutkan(c){
    if(c.sortable===false) return;
    setSort(s=>({key:c.key, dir:s.key===c.key&&s.dir==='asc'?'desc':'asc'}));
  }

  let isi;
  if(server && srv.error) isi = <div className="border border-line rounded-card bg-surface"><GagalMuat error={srv.error} onRetry={()=>setUlang(u=>u+1)}/></div>;
  else if(memuatAwal) isi = <div className="border border-line rounded-card bg-surface p-5"><Memuat baris={pageSize}/></div>;
  else if(view.length===0) isi = <div className="border border-line rounded-card bg-surface">{empty || <EmptyState title="Tidak ada data yang cocok" desc="Ubah kata kunci pencarian atau atur ulang filter."/>}</div>;
  else isi = <div className={`transition-opacity ${server&&srv.loading?'opacity-60':''}`} aria-busy={server&&srv.loading}>
    {/* desktop */}
    <div className="hidden md:block border border-line rounded-card bg-surface overflow-hidden">
      <div className="scroll-x">
        <table className="w-full text-[14px]" style={{minWidth:'700px'}}>
          <thead><tr className="bg-[var(--ground)] text-left">
            {columns.map(c=>(
              <th key={c.key} className={`px-4 py-3 ${c.align==='right'?'text-right':''}`} aria-sort={sort.key===c.key?(sort.dir==='asc'?'ascending':'descending'):undefined}>
                {c.sortable===false
                  ? <span className="eyebrow text-[11px] font-bold text-muted">{c.label}</span>
                  : <button onClick={()=>urutkan(c)}
                      className={`eyebrow text-[11px] font-bold text-muted hover:text-ink inline-flex items-center gap-1 ${c.align==='right'?'flex-row-reverse':''}`}>
                      {c.label}{sort.key===c.key && <Icon name={sort.dir==='asc'?'chevronUp':'chevronDown'} className="w-3.5 h-3.5"/>}
                    </button>}
              </th>
            ))}
            {onRow && <th className="px-4 py-3 w-16"></th>}
          </tr></thead>
          <tbody className="divide-y divide-[color:var(--line-soft)]">
            {view.map((r,i)=>(
              <tr key={r.id||i} className="hover:bg-[var(--ground)]">
                {columns.map(c=> <td key={c.key} className={`px-4 py-3 ${c.align==='right'?'text-right':''} ${c.mono?'font-num':''}`}>{c.cell?c.cell(r):r[c.key]}</td>)}
                {onRow && <td className="px-4 py-3 text-right">
                  <button onClick={()=>onRow(r)} aria-label="Lihat detail" className="p-1.5 rounded-[6px] text-muted hover:text-[var(--brand)] hover:bg-brandsoft"><Icon name="eye" className="w-[17px] h-[17px]"/></button>
                </td>}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
    {/* mobile: kartu */}
    <ul className="md:hidden space-y-2.5">
      {view.map((r,i)=>(
        <li key={r.id||i} className={`bg-surface border border-line rounded-card p-4 ${onRow?'cursor-pointer':''}`} onClick={()=>onRow&&onRow(r)}>
          {columns.map((c,ci)=>(
            <div key={c.key} className={`flex gap-3 ${ci?'mt-2':''}`}>
              <span className="text-[12px] text-muted w-[92px] shrink-0">{c.label}</span>
              <span className={`text-[13.5px] flex-1 min-w-0 ${c.mono?'font-num':''}`}>{c.cell?c.cell(r):r[c.key]}</span>
            </div>
          ))}
        </li>
      ))}
    </ul>
  </div>;

  return <div>
    <div className="flex flex-col sm:flex-row gap-3 mb-4">
      <div className="relative flex-1">
        <Icon name="search" className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-muted"/>
        <input value={q} onChange={e=>setQ(e.target.value)} placeholder="Cari…" aria-label="Cari data" maxLength={100}
          className={`${inputCls} pl-9`}/>
      </div>
      <span className="text-[13px] text-muted self-center font-num whitespace-nowrap">{memuatAwal?'Memuat…':`${fmtInt(total)} baris`}</span>
    </div>
    {isi}
    <div className="flex items-center justify-between gap-3 mt-4">
      <p className="text-[13px] text-muted">Halaman <span className="font-num text-ink">{page}</span> dari <span className="font-num">{pages}</span></p>
      <div className="flex gap-2">
        <Button size="sm" variant="outline" icon="chevronLeft" disabled={page<=1 || (server&&srv.loading)} onClick={()=>keHalaman(page-1)}>Sebelumnya</Button>
        <Button size="sm" variant="outline" iconRight="chevronRight" disabled={page>=pages || (server&&srv.loading)} onClick={()=>keHalaman(page+1)}>Berikutnya</Button>
      </div>
    </div>
  </div>;
}

export default DataTable;
