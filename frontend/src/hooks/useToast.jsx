import { createContext, useCallback, useContext, useState } from 'react';
import Toast from '../components/ui/Toast.jsx';

const ToastContext = createContext(() => {});

/** toast('Pesan') untuk notifikasi sukses, toast('Pesan', 'galat') untuk kegagalan. */
export function ToastProvider({ children }) {
  const [msg, setMsg] = useState(null);
  const show = useCallback((teks, tone = 'ok') => setMsg({ teks, tone, id: Date.now() }), []);
  const tutup = useCallback(() => setMsg(null), []);
  return (
    <ToastContext.Provider value={show}>
      {children}
      <Toast msg={msg} onDone={tutup} />
    </ToastContext.Provider>
  );
}

export const useToast = () => useContext(ToastContext);
