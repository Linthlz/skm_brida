import { SKALA } from '../../data/skala.js';

/** `skala` berasal dari skalaDari(api) — 4 atau 5 tingkat; bawaannya SKALA lokal 1–5. */
function RatingScale({value, onChange, name, invalid, skala=SKALA}){
  return <div role="radiogroup" aria-label={name} className={`grid gap-1.5 sm:gap-2 ${invalid?'rounded-[9px] ring-2 ring-[color:var(--bad)]/40 p-1 -m-1':''}`}
    style={{gridTemplateColumns:`repeat(${skala.length}, minmax(0, 1fr))`}}>
    {skala.map(s=>{
      const on = value===s.v;
      return (
        <button key={s.v} type="button" role="radio" aria-checked={on} onClick={()=>onChange(s.v)}
          className={`min-h-[76px] sm:min-h-[84px] rounded-[9px] border px-1 py-2.5 flex flex-col items-center justify-center gap-1.5 transition-all
            ${on ? 'shadow-card' : 'bg-surface border-line text-ink2 hover:border-[var(--brand)] hover:bg-brandsoft'}`}
          style={on?{background:s.warna, borderColor:s.warna, color:s.ink}:null}>
          <span className="font-num text-[17px] font-semibold">{s.v}</span>
          <span className={`text-[10.5px] sm:text-[11.5px] font-semibold leading-tight text-center ${on?'opacity-90':'text-muted'}`}>{s.pendek}</span>
        </button>
      );
    })}
  </div>;
}

export default RatingScale;
