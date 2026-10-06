import Badge from '../ui/Badge.jsx';
import { fmt } from '../../utils/format.js';
import { mutu } from '../../utils/mutu.js';

function MeterIKM({nilai}){
  const m = mutu(nilai);
  const pct = Math.max(0,Math.min(100,nilai));
  const bands=[{to:65,c:'var(--s1)',l:'D'},{to:76.61,c:'var(--s2)',l:'C'},{to:88.31,c:'var(--s4)',l:'B'},{to:100,c:'var(--s5)',l:'A'}];
  let prev=25;
  return <div>
    <div className="flex items-end justify-between gap-4 mb-3">
      <div>
        <p className="eyebrow text-[10.5px] font-bold text-muted mb-1.5">Indeks Kepuasan Masyarakat</p>
        <p className="font-num text-[44px] font-semibold leading-none">{fmt(nilai)}</p>
      </div>
      <Badge tone={m.tone} className="mb-1">Mutu {m.huruf} · {m.label}</Badge>
    </div>
    <div className="relative h-2.5 rounded-full overflow-hidden flex" style={{background:'var(--ground)'}}>
      {bands.map(b=>{ const wpc=((b.to-prev)/75)*100; prev=b.to;
        return <span key={b.l} style={{width:wpc+'%', background:b.c, opacity:.28}}/>; })}
      <span className="absolute top-0 bottom-0 w-[3px] rounded-full bg-[var(--ink)]" style={{left:`calc(${((pct-25)/75)*100}% - 1.5px)`}}/>
    </div>
    <div className="flex justify-between mt-2 text-[11px] text-muted font-num">
      <span>25,00</span><span>65,00</span><span>76,61</span><span>88,31</span><span>100</span>
    </div>
  </div>;
}

export default MeterIKM;
