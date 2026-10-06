import { useCallback, useEffect, useRef, useState } from 'react';

/**
 * Memanggil fungsi API dan melacak status muat. `deps` menentukan kapan data
 * diambil ulang. Respons lama yang datang terlambat diabaikan.
 *
 *   const { data, error, loading, reload } = useApi(() => hasil({ periode }), [periode]);
 */
export function useApi(fn, deps = []) {
  const [state, setState] = useState({ data: null, error: null, loading: true });
  const giliran = useRef(0);

  const muat = useCallback(() => {
    const ke = ++giliran.current;
    setState((s) => ({ ...s, error: null, loading: true }));
    fn().then(
      (data) => ke === giliran.current && setState({ data, error: null, loading: false }),
      (error) => ke === giliran.current && setState((s) => ({ data: s.data, error, loading: false })),
    );
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, deps);

  useEffect(() => {
    muat();
    return () => { giliran.current++; };
  }, [muat]);

  const setData = useCallback((f) => setState((s) => ({ ...s, data: typeof f === 'function' ? f(s.data) : f })), []);

  return { ...state, reload: muat, setData };
}
