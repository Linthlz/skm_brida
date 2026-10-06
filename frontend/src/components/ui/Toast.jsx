import { useEffect } from 'react';
import Icon from '../Icon.jsx';

function Toast({msg, onDone}){
  useEffect(()=>{ if(!msg) return; const t=setTimeout(onDone, msg.tone==='galat'?4200:2600); return ()=>clearTimeout(t); },[msg,onDone]);
  if(!msg) return null;
  return <div role="status" aria-live="polite" className="fixed left-1/2 -translate-x-1/2 z-[90] bg-[#5A1010] text-white text-[14px] border border-white/15 font-semibold px-4 py-3 rounded-[9px] shadow-pop flex items-center gap-2 anim-rise"
              style={{bottom:'calc(1.25rem + env(safe-area-inset-bottom,0px))', maxWidth:'calc(100vw - 32px)'}}>
    <Icon name={msg.tone==='galat'?'alert':'checkCircle'} className="w-[18px] h-[18px] shrink-0"/>{msg.teks}
  </div>;
}

export default Toast;
