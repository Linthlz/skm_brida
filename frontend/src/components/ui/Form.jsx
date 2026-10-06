import Icon from '../Icon.jsx';

function Field({label, hint, required, error, children, className=''}){
  return <label className={`block ${className}`}>
    <span className="block text-[13.5px] font-semibold mb-1.5">{label}{required && <span className="text-[var(--bad)]"> *</span>}{!required && <span className="text-muted font-normal"> (opsional)</span>}</span>
    {children}
    {hint && !error && <span className="block mt-1.5 text-[12.5px] text-muted">{hint}</span>}
    {error && <span className="mt-1.5 flex items-center gap-1.5 text-[12.5px] font-semibold text-[var(--bad)]"><Icon name="alert" className="w-4 h-4"/>{error}</span>}
  </label>;
}
const inputCls = 'w-full bg-surface border border-line rounded-[8px] px-3 py-2.5 text-[15px] text-ink placeholder:text-muted/70 focus:border-[var(--brand)] transition-colors';
const Input    = (p)=> <input {...p} className={`${inputCls} ${p.className||''}`}/>;
const Textarea = (p)=> <textarea {...p} className={`${inputCls} min-h-[120px] resize-y ${p.className||''}`}/>;
/** `options` berisi string, atau {value, label} bila nilai berbeda dari teksnya (mis. id layanan).
 *  placeholder={null} menghilangkan pilihan kosong. */
function Select({options=[], placeholder='Pilih salah satu', ...p}){
  return <div className="relative">
    <select {...p} className={`${inputCls} appearance-none pr-10 ${p.className||''}`}>
      {placeholder!==null && <option value="">{placeholder}</option>}
      {options.map(o=> typeof o==='object'
        ? <option key={o.value} value={o.value}>{o.label}</option>
        : <option key={o} value={o}>{o}</option>)}
    </select>
    <Icon name="chevronDown" className="w-4 h-4 absolute right-3 top-1/2 -translate-y-1/2 text-muted pointer-events-none"/>
  </div>;
}

export { Field, Input, Textarea, Select, inputCls };
