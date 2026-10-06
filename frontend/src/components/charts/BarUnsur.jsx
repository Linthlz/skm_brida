import { useState } from 'react';
import { useSize } from '../../hooks/useSize.js';
import Tooltip from '../ui/Tooltip.jsx';
import { fmt } from '../../utils/format.js';

/** `max` = nilai tertinggi skala, `pengali` = konversi NRR ke indeks (100 / max). */
function BarUnsur({data, max=5, pengali=100/max}){
  const [ref,w] = useSize();
  const [hov,setHov] = useState(null);
  const labelW = w<520 ? 0 : 210;
  const padR = 46, rowH = w<520 ? 50 : 36, top = 22;
  const h = top + data.length*rowH + 10;
  const x0 = labelW, x1 = Math.max(x0+80, w - padR);
  const sc = v => x0 + ((v-1)/(max-1))*(x1-x0);
  const ticks = Array.from({length:max}, (_,i)=>i+1);
  return <div ref={ref} className="relative">
    <svg width={w} height={h} role="img" aria-label="Nilai rata-rata tiap unsur pelayanan">
      {ticks.map(t=> <g key={t}>
        <line x1={sc(t)} x2={sc(t)} y1={top-8} y2={h-16} stroke="var(--grid)" strokeWidth="1"/>
        <text x={sc(t)} y={h-4} textAnchor="middle" fontSize="11" fill="var(--axis)" fontFamily="IBM Plex Mono">{t}</text>
      </g>)}
      {data.map((d,i)=>{
        const y = top + i*rowH;
        const barY = w<520 ? y+20 : y+9;
        const on = hov===i;
        return <g key={d.kode} onMouseEnter={()=>setHov(i)} onMouseLeave={()=>setHov(null)}>
          <rect x="0" y={y-4} width={w} height={rowH} fill={on?'var(--brand-soft)':'transparent'}/>
          {w<520
            ? <text x="2" y={y+12} fontSize="12.5" fill="var(--ink)" fontWeight="600">{d.kode} · {d.nama}</text>
            : <text x={labelW-12} y={y+18} textAnchor="end" fontSize="12.5" fill="var(--ink)">{d.nama}</text>}
          <rect x={x0} y={barY} width={Math.max(2, sc(d.nilai)-x0)} height="12" rx="4" fill="var(--brand)" opacity={on?1:.9}/>
          <text x={sc(d.nilai)+8} y={barY+11} fontSize="12.5" fill="var(--ink)" fontWeight="600" fontFamily="IBM Plex Mono">{fmt(d.nilai)}</text>
        </g>;
      })}
    </svg>
    {hov!==null && <Tooltip x={Math.min(w-20, sc(data[hov].nilai))} y={top + hov*rowH + (w<520?24:14)}>
      <b>{data[hov].kode} · {data[hov].nama}</b><br/>Nilai rata-rata {fmt(data[hov].nilai)} dari {fmt(max)}<br/>Konversi indeks {fmt(data[hov].nilai*pengali)}
    </Tooltip>}
  </div>;
}

export default BarUnsur;
