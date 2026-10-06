import { useState } from 'react';
import { useSize } from '../../hooks/useSize.js';
import { fmt, fmtInt } from '../../utils/format.js';

function DonutKepuasan({data}){
  const [ref,w] = useSize();
  const [hov,setHov] = useState(null);
  const stack = w < 440;
  const size = stack ? Math.min(w, 268) : Math.min(w*0.44, 236);
  const cx=size/2, cy=size/2, rOut=size/2-6, rIn=size*0.30;
  const total = data.reduce((a,d)=>a+d.jumlah,0);
  let acc = -Math.PI/2;
  const arcs = data.map((d,i)=>{
    const frac = d.jumlah/total, gap = 0.012;
    const a0 = acc+gap/2, a1 = acc+frac*2*Math.PI-gap/2; acc += frac*2*Math.PI;
    const pt=(r,a)=>[cx+r*Math.cos(a), cy+r*Math.sin(a)];
    const big = (a1-a0)>Math.PI?1:0;
    const [x1,y1]=pt(rOut,a0),[x2,y2]=pt(rOut,a1),[x3,y3]=pt(rIn,a1),[x4,y4]=pt(rIn,a0);
    return {...d, i, frac, d:`M${x1} ${y1}A${rOut} ${rOut} 0 ${big} 1 ${x2} ${y2}L${x3} ${y3}A${rIn} ${rIn} 0 ${big} 0 ${x4} ${y4}Z`};
  });
  const cur = hov!==null ? arcs[hov] : null;
  return <div ref={ref} className={`flex ${stack?'flex-col':'flex-row'} items-center gap-5 w-full min-w-0`}>
    <div className="relative shrink-0" style={{width:size, height:size}}>
      <svg width={size} height={size} role="img" aria-label="Distribusi tingkat kepuasan responden">
        {arcs.map(a=> <path key={a.v} d={a.d} fill={a.warna} stroke="var(--surface)" strokeWidth="2"
          opacity={hov===null||hov===a.i?1:.35} onMouseEnter={()=>setHov(a.i)} onMouseLeave={()=>setHov(null)}/>)}
      </svg>
      <div className="absolute inset-0 grid place-content-center text-center pointer-events-none">
        <div className="font-num text-[24px] font-semibold leading-none">{cur? Math.round(cur.frac*1000)/10+'%' : fmtInt(total)}</div>
        <div className="text-[11.5px] text-muted mt-1.5 px-6">{cur? cur.label : 'Responden'}</div>
      </div>
    </div>
    <ul className="w-full min-w-0 flex-1 space-y-1.5">
      {arcs.map(a=> <li key={a.v} onMouseEnter={()=>setHov(a.i)} onMouseLeave={()=>setHov(null)}
        className={`flex items-center gap-3 px-2.5 py-2 rounded-[7px] ${hov===a.i?'bg-brandsoft':''}`}>
        <span className="w-3 h-3 rounded-[3px] shrink-0" style={{background:a.warna}}/>
        <span className="text-[13.5px] flex-1 min-w-0">{a.label}</span>
        <span className="font-num text-[13px] text-ink2">{fmtInt(a.jumlah)}</span>
        <span className="font-num text-[13px] font-semibold w-[52px] text-right">{fmt(a.frac*100,1)}%</span>
      </li>)}
    </ul>
  </div>;
}

export default DonutKepuasan;
