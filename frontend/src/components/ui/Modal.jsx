import { useEffect } from 'react';
import Icon from '../Icon.jsx';

function Modal({open, onClose, title, desc, children, footer, width='max-w-lg'}){
  useEffect(()=>{
    if(!open) return;
    const h = e=> e.key==='Escape' && onClose();
    window.addEventListener('keydown',h); return ()=>window.removeEventListener('keydown',h);
  },[open,onClose]);
  if(!open) return null;
  return (
    <div className="fixed inset-0 z-[80] flex items-end sm:items-center justify-center p-0 sm:p-4" role="dialog" aria-modal="true" aria-label={title}>
      <div className="absolute inset-0 bg-[#2A0B0B]/55 backdrop-blur-[2px]" onClick={onClose}/>
      <div className={`relative w-full ${width} bg-surface border border-line rounded-t-[14px] sm:rounded-card shadow-pop anim-rise max-h-[90vh] flex flex-col`}>
        <div className="flex items-start justify-between gap-4 px-5 sm:px-6 pt-5 pb-4 border-b border-linesoft">
          <div>
            <h3 className="font-display text-[19px] font-bold">{title}</h3>
            {desc && <p className="text-[13.5px] text-muted mt-1">{desc}</p>}
          </div>
          <button onClick={onClose} aria-label="Tutup" className="shrink-0 p-1.5 rounded-[7px] text-muted hover:text-ink hover:bg-[var(--ground)]"><Icon name="x" className="w-5 h-5"/></button>
        </div>
        <div className="px-5 sm:px-6 py-5 overflow-y-auto">{children}</div>
        {footer && <div className="px-5 sm:px-6 py-4 border-t border-linesoft flex flex-col-reverse sm:flex-row sm:justify-end gap-2.5" style={{paddingBottom:'calc(1rem + env(safe-area-inset-bottom,0px))'}}>{footer}</div>}
      </div>
    </div>
  );
}

export default Modal;
