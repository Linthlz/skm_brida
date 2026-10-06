import { useState } from 'react';
import { useSize } from '../../hooks/useSize.js';
import Tooltip from '../ui/Tooltip.jsx';
import { fmt } from '../../utils/format.js';
import { mutu } from '../../utils/mutu.js';

function TrenIKM({data}){
  const [ref,w] = useSize();
  const [hov,setHov] = useState(null);
  const h=230, padL=42, padR=16, padT=18, padB=32;
  // Rentang sumbu Y melebar bila ada nilai di luar 78–92, dibulatkan ke kelipatan 4.
  const nilai = data.map(d=>d.ikm);
  const lo = Math.min(78, Math.floor((Math.min(...nilai)-2)/4)*4), hi = Math.max(92, Math.ceil((Math.max(...nilai)+2)/4)*4);
  const x = i => padL + i*((w-padL-padR)/Math.max(1,data.length-1));
  const y = v => padT + (1-(v-lo)/(hi-lo))*(h-padT-padB);
  const line = data.map((d,i)=>`${i?'L':'M'}${x(i)} ${y(d.ikm)}`).join(' ');
  const area = `${line} L${x(data.length-1)} ${h-padB} L${x(0)} ${h-padB} Z`;
  const gy = Array.from({length:Math.floor((hi-lo)/4)}, (_,i)=>lo+4*(i+1)).filter(g=>g<=hi);
  return <div ref={ref} className="relative">
    <svg width={w} height={h} role="img" aria-label="Tren Indeks Kepuasan Masyarakat per periode survei"
      onMouseLeave={()=>setHov(null)}
      onMouseMove={e=>{ const r=e.currentTarget.getBoundingClientRect(); const px=e.clientX-r.left;
        let best=0,bd=1e9; data.forEach((d,i)=>{const dd=Math.abs(x(i)-px); if(dd<bd){bd=dd;best=i;}}); setHov(best); }}>
      {gy.map(g=> <g key={g}>
        <line x1={padL} x2={w-padR} y1={y(g)} y2={y(g)} stroke="var(--grid)" strokeWidth="1"/>
        <text x={padL-8} y={y(g)+4} textAnchor="end" fontSize="11" fill="var(--axis)" fontFamily="IBM Plex Mono">{g}</text>
      </g>)}
      <path d={area} fill="var(--brand)" opacity=".10"/>
      <path d={line} fill="none" stroke="var(--brand)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
      {hov!==null && <line x1={x(hov)} x2={x(hov)} y1={padT-6} y2={h-padB} stroke="var(--brand)" strokeWidth="1" strokeDasharray="3 3" opacity=".6"/>}
      {data.map((d,i)=> <g key={d.periode}>
        <circle cx={x(i)} cy={y(d.ikm)} r={i===data.length-1?6:4.5} fill={i===data.length-1?'var(--gold-bright)':'var(--brand)'} stroke="var(--surface)" strokeWidth="2"/>
        <text x={x(i)} y={h-10} textAnchor="middle" fontSize="11.5" fill="var(--axis)" fontFamily="IBM Plex Mono">{d.periode}</text>
      </g>)}
      <text x={x(data.length-1)} y={y(data[data.length-1].ikm)-14} textAnchor="end" fontSize="12.5" fontWeight="700" fill="var(--ink)" fontFamily="IBM Plex Mono">{fmt(data[data.length-1].ikm)}</text>
    </svg>
    {hov!==null && <Tooltip x={x(hov)} y={y(data[hov].ikm)}>Periode {data[hov].periode}<br/><b className="font-num">{fmt(data[hov].ikm)}</b> · Mutu {mutu(data[hov].ikm).huruf}</Tooltip>}
  </div>;
}

export default TrenIKM;
