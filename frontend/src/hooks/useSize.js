import { useEffect, useRef, useState } from 'react';

/** Mengukur lebar container agar grafik SVG bisa responsif tanpa distorsi teks. */
export function useSize(initial = 680) {
  const ref = useRef(null);
  const [w, setW] = useState(initial);
  useEffect(() => {
    if (!ref.current) return;
    const ro = new ResizeObserver((entries) => {
      const cw = entries[0].contentRect.width;
      if (cw > 0) setW(cw);
    });
    ro.observe(ref.current);
    return () => ro.disconnect();
  }, []);
  return [ref, w];
}
