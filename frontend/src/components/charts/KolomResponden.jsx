import { useState } from 'react';
import { useSize } from '../../hooks/useSize.js';
import Tooltip from '../ui/Tooltip.jsx';

function KolomResponden({labels, values, tahun}){
  const [ref,w] = useSize();
  const [hov,setHov] = useState(null);
  const h=210, padT=16, padB=28, padL=34, padR=6;
  const max = Math.max(20, Math.ceil(Math.max(0, ...values)/20)*20);
  const bw = (w-padL-padR)/labels.length;
  const y = v => padT + (1-v/max)*(h-padT-padB);
  return <div ref={ref} className="relative">
    <svg width={w} height={h} role="img" aria-label={`Jumlah responden per bulan tahun ${tahun}`}>
      {[0,max/2,max].map(g=> <g key={g}>
        <line x1={padL} x2={w-padR} y1={y(g)} y2={y(g)} stroke="var(--grid)" strokeWidth="1"/>
        <text x={padL-7} y={y(g)+4} textAnchor="end" fontSize="10.5" fill="var(--axis)" fontFamily="IBM Plex Mono">{g}</text>
      </g>)}
      {values.map((v,i)=>{
        const bx = padL + i*bw + 3, bwid = Math.max(4, bw-8);
        return <g key={i} onMouseEnter={()=>setHov(i)} onMouseLeave={()=>setHov(null)}>
          <rect x={padL+i*bw} y={padT-6} width={bw} height={h-padT-padB+6} fill="transparent"/>
          <rect x={bx} y={y(v)} width={bwid} height={h-padB-y(v)} rx="4" fill="var(--brand)" opacity={hov===null||hov===i?.92:.4}/>
          <text x={bx+bwid/2} y={h-9} textAnchor="middle" fontSize="11" fill="var(--axis)">{labels[i]}</text>
        </g>;
      })}
    </svg>
    {hov!==null && <Tooltip x={padL+hov*bw+bw/2} y={y(values[hov])}>{labels[hov]} {tahun}<br/><b className="font-num">{values[hov]}</b> responden</Tooltip>}
  </div>;
}

export default KolomResponden;
