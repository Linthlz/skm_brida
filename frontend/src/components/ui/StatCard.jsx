import Icon from '../Icon.jsx';

function StatCard({label, value, sub, icon, tone='brand', big=false}){
  const c = tone==='gold'?'var(--gold)':tone==='ok'?'var(--ok)':'var(--brand)';
  return (
    <div className="bg-surface border border-line rounded-card p-5 flex flex-col gap-3">
      <div className="flex items-start justify-between gap-3">
        <span className="eyebrow text-[10.5px] font-bold text-muted">{label}</span>
        {icon && <span className="shrink-0 grid place-items-center w-8 h-8 rounded-[7px]" style={{background:'color-mix(in srgb,'+c+' 12%, transparent)', color:c}}><Icon name={icon} className="w-[17px] h-[17px]"/></span>}
      </div>
      <div className="flex items-baseline gap-2 flex-wrap">
        <span className={`font-num font-semibold leading-none ${big?'text-[38px] sm:text-[44px]':'text-[30px]'}`} style={{color:'var(--ink)'}}>{value}</span>
        {sub && <span className="text-[13px] text-muted">{sub}</span>}
      </div>
    </div>
  );
}

export default StatCard;
