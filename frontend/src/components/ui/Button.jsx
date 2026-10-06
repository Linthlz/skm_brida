import Icon from '../Icon.jsx';

function Button({as='button', variant='primary', size='md', icon, iconRight, children, className='', ...rest}){
  const base='inline-flex items-center justify-center gap-2 font-semibold rounded-[8px] transition-colors duration-150 select-none disabled:opacity-50 disabled:cursor-not-allowed';
  const sizes={sm:'text-[13px] px-3 py-1.5', md:'text-[15px] px-4 py-2.5', lg:'text-[16px] px-5 py-3'};
  const vars={
    primary:'bg-[var(--brand)] text-[var(--on-brand)] hover:brightness-110 shadow-card',
    gold:'bg-[var(--gold-bright)] text-[var(--gold-ink)] hover:brightness-105 shadow-card',
    outline:'border border-line text-ink hover:bg-brandsoft hover:border-[var(--brand)]',
    ghost:'text-ink hover:bg-brandsoft',
    quiet:'text-muted hover:text-ink hover:bg-[var(--ground)]',
    onnav:'bg-[var(--nav-active)] text-[var(--nav-active-ink)] hover:brightness-95 shadow-card',
    onnavGhost:'border border-[var(--nav-line)] text-[var(--nav-ink)] hover:bg-[var(--nav-hover)]',
    danger:'border border-[var(--bad)] text-[var(--bad)] hover:bg-[var(--bad)] hover:text-white',
  };
  const El = as;
  return <El className={`${base} ${sizes[size]} ${vars[variant]} ${className}`} {...rest}>
    {icon && <Icon name={icon} className={size==='sm'?'w-4 h-4':'w-[18px] h-[18px]'}/>}
    {children}
    {iconRight && <Icon name={iconRight} className={size==='sm'?'w-4 h-4':'w-[18px] h-[18px]'}/>}
  </El>;
}

export default Button;
