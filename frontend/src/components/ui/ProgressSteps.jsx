import Icon from '../Icon.jsx';

function ProgressSteps({steps, current}){
  const pct = (current/(steps.length-1))*100;
  return <div>
    <div className="flex items-center justify-between mb-3">
      <p className="text-[13px] font-semibold text-muted">Langkah <span className="font-num text-ink">{current+1}</span> dari <span className="font-num">{steps.length}</span></p>
      <p className="text-[13px] font-semibold text-[var(--brand)]">{steps[current]}</p>
    </div>
    <div className="h-1.5 rounded-full bg-[var(--ground)] overflow-hidden">
      <div className="h-full rounded-full bg-[var(--brand)] transition-[width] duration-300" style={{width:`${Math.max(6,pct)}%`}}/>
    </div>
    <ol className="hidden sm:flex items-center gap-2 mt-4">
      {steps.map((s,i)=>(
        <li key={s} className="flex items-center gap-2 flex-1 last:flex-none">
          <span className={`font-num text-[11px] font-semibold w-6 h-6 grid place-items-center rounded-full shrink-0 border
            ${i<current?'bg-[var(--brand)] text-[var(--on-brand)] border-[var(--brand)]': i===current?'border-[var(--brand)] text-[var(--brand)] bg-brandsoft':'border-line text-muted'}`}>
            {i<current ? <Icon name="check" className="w-3.5 h-3.5" sw={2.5}/> : String(i+1).padStart(2,'0')}
          </span>
          <span className={`text-[12.5px] font-semibold ${i<=current?'text-ink':'text-muted'}`}>{s}</span>
          {i<steps.length-1 && <span className="flex-1 h-px bg-line ml-1"/>}
        </li>
      ))}
    </ol>
  </div>;
}

export default ProgressSteps;
