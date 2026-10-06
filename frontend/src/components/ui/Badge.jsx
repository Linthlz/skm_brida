function Badge({tone='brand', children, dot=false, className=''}){
  const map={
    brand:'bg-brandsoft text-[var(--brand)] border-[color:var(--brand)]/30',
    ok:'bg-[color:var(--ok)]/12 text-[var(--ok)] border-[color:var(--ok)]/35',
    warn:'bg-[color:var(--warn)]/12 text-[var(--warn)] border-[color:var(--warn)]/35',
    bad:'bg-[color:var(--bad)]/12 text-[var(--bad)] border-[color:var(--bad)]/35',
    gold:'bg-goldsoft text-[var(--gold)] border-[color:var(--gold)]/35',
    neutral:'bg-[var(--ground)] text-muted border-line',
  };
  return <span className={`inline-flex items-center gap-1.5 text-[12px] font-semibold px-2.5 py-1 rounded-full border ${map[tone]} ${className}`}>
    {dot && <span className="w-1.5 h-1.5 rounded-full bg-current"/>}{children}
  </span>;
}

export default Badge;
