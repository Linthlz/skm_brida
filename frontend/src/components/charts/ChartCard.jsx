function ChartCard({title, desc, legend, children, right}){
  return <div className="bg-surface border border-line rounded-card p-5 sm:p-6 min-w-0">
    <div className="flex items-start justify-between gap-4 mb-1">
      <h3 className="font-display text-[17px] font-bold leading-snug">{title}</h3>
      {right}
    </div>
    {desc && <p className="text-[13.5px] text-muted mb-4" style={{maxWidth:'58ch'}}>{desc}</p>}
    {legend && <div className="flex flex-wrap gap-x-4 gap-y-2 mb-4">{legend}</div>}
    {children}
  </div>;
}

const LegendDot = ({ color, children }) => (
  <span className="inline-flex items-center gap-2 text-[12.5px] text-ink2">
    <span className="w-2.5 h-2.5 rounded-[3px] shrink-0" style={{ background: color }} />
    {children}
  </span>
);

export { ChartCard, LegendDot };
export default ChartCard;
