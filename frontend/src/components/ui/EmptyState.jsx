import Icon from '../Icon.jsx';

function EmptyState({icon='search', title, desc, action}){
  return <div className="py-14 px-6 text-center">
    <div className="mx-auto w-12 h-12 grid place-items-center rounded-full bg-[var(--ground)] text-muted mb-4"><Icon name={icon} className="w-6 h-6"/></div>
    <p className="font-display font-bold text-[17px]">{title}</p>
    <p className="text-[14px] text-muted mt-1.5 mx-auto" style={{maxWidth:'44ch'}}>{desc}</p>
    {action && <div className="mt-5">{action}</div>}
  </div>;
}

export default EmptyState;
