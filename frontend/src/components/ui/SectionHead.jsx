function SectionHead({eyebrow, title, desc, align='left', className=''}){
  return <div className={`${align==='center'?'text-center mx-auto':''} ${className}`} style={align==='center'?{maxWidth:'62ch'}:null}>
    {eyebrow && <p className="eyebrow text-[11px] font-bold text-[var(--gold)] mb-2">{eyebrow}</p>}
    <h2 className="font-display text-[26px] sm:text-[32px] font-extrabold leading-[1.15] tracking-[-0.01em]">{title}</h2>
    {desc && <p className="mt-3 text-[16px] text-ink2 leading-relaxed" style={{maxWidth:'62ch'}}>{desc}</p>}
  </div>;
}

export default SectionHead;
