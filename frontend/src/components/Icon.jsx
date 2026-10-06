/** Ikon garis 24x24 bergaya Lucide, ditulis inline agar tidak menambah dependensi runtime.
 *  Ganti dengan `lucide-react` bila ingin katalog ikon penuh. */
const C=(x,y,r)=>({t:'c',x,y,r}), R=(x,y,w,h,rx)=>({t:'r',x,y,w,h,rx});
const PATHS = {
  menu:['M3 6h18','M3 12h18','M3 18h18'],
  x:['M18 6 6 18','M6 6l12 12'],
  chevronDown:['m6 9 6 6 6-6'], chevronUp:['m18 15-6-6-6 6'],
  chevronLeft:['m15 18-6-6 6-6'], chevronRight:['m9 18 6-6-6-6'],
  check:['M20 6 9 17l-5-5'],
  checkCircle:[C(12,12,10),'m8.5 12 2.5 2.5 4.5-5'],
  clock:[C(12,12,10),'M12 6.5V12l3.5 2'],
  users:['M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2',C(9,7,4),'M22 21v-2a4 4 0 0 0-3-3.87','M16 3.13a4 4 0 0 1 0 7.75'],
  user:['M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2',C(12,7,4)],
  target:[C(12,12,9),C(12,12,5),C(12,12,1.5)],
  award:[C(12,8,6),'M15.6 13.4 17 22l-5-2.8L7 22l1.4-8.6'],
  barChart:['M3 3v16a2 2 0 0 0 2 2h16','M7.5 17v-4.5','M12 17V8','M16.5 17v-7'],
  pieChart:['M21.2 15.9A10 10 0 1 1 8.1 2.8','M22 12A10 10 0 0 0 12 2v10z'],
  trendingUp:['m22 7-8.5 8.5-5-5L2 17','M16 7h6v6'],
  fileText:['M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7z','M15 2v5h5','M9 13h6','M9 17h4'],
  message:['M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z'],
  sliders:['M4 21v-7','M4 10V3','M12 21v-9','M12 8V3','M20 21v-5','M20 12V3','M1.5 14h5','M9.5 8h5','M17.5 16h5'],
  search:[C(11,11,7.5),'m20.5 20.5-3.7-3.7'],
  filter:['M21.5 4h-19l7.6 8.6V19l3.8 1.8v-8.2z'],
  download:['M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4','m7 10 5 5 5-5','M12 15V3'],
  plus:['M12 5v14','M5 12h14'],
  pencil:['M12.5 20H21','M16.5 3.5a2.12 2.12 0 0 1 3 3L7.5 18.5l-4 1 1-4z'],
  trash:['M3.5 6h17','M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6','M8.5 6V4.5a2 2 0 0 1 2-2h3a2 2 0 0 1 2 2V6','M10 11v6','M14 11v6'],
  eye:['M2 12s3.6-7 10-7 10 7 10 7-3.6 7-10 7-10-7-10-7Z',C(12,12,3)],
  logOut:['M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4','m16 17 5-5-5-5','M21 12H9'],
  grid:[R(3,3,7,8,1.5),R(14,3,7,5,1.5),R(14,11,7,10,1.5),R(3,15,7,6,1.5)],
  listChecks:['M11 6h10','M11 12h10','M11 18h10','m3 6 1.5 1.5L7.5 4.5','m3 12 1.5 1.5L7.5 10.5','m3 18 1.5 1.5L7.5 16.5'],
  mail:[R(2,4.5,20,15,2),'m22 6.5-9.1 6a2 2 0 0 1-1.8 0L2 6.5'],
  phone:['M6.5 2.5h-3a2 2 0 0 0-2 2.2A18 18 0 0 0 19.3 22.5a2 2 0 0 0 2.2-2v-3a2 2 0 0 0-1.7-2l-2.6-.4a2 2 0 0 0-2 .9l-.6 1a14 14 0 0 1-5.6-5.6l1-.6a2 2 0 0 0 .9-2l-.4-2.6a2 2 0 0 0-2-1.7Z'],
  mapPin:['M20 10.5c0 6-8 11.5-8 11.5S4 16.5 4 10.5a8 8 0 0 1 16 0Z',C(12,10.3,2.8)],
  shield:['M12 22s8-3.8 8-10V5.2L12 2 4 5.2V12c0 6.2 8 10 8 10Z','m9 12 2 2 4-4'],
  arrowRight:['M4.5 12h15','m13 5.5 6.5 6.5-6.5 6.5'],
  arrowLeft:['M19.5 12h-15','m11 5.5-6.5 6.5 6.5 6.5'],
  alert:[C(12,12,10),'M12 7.5v5','M12 16.3h.01'],
  star:['M12 2.6 15 8.7l6.7 1-4.85 4.7 1.15 6.7L12 18l-6 3.1L7.15 14.4 2.3 9.7l6.7-1z'],
  building:[R(4,2.5,16,19,2),'M9.5 7h.01','M14.5 7h.01','M9.5 11.5h.01','M14.5 11.5h.01','M10 21.5V17h4v4.5'],
  lock:[R(4,10.5,16,10.5,2),'M8 10.5V7a4 4 0 0 1 8 0v3.5'],
  calendar:[R(3,4.5,18,17,2),'M16 2.5v4','M8 2.5v4','M3 10h18'],
  refresh:['M21 12a9 9 0 1 1-2.6-6.4','M21 3v5h-5'],
  sparkle:['M12 3v4','M12 17v4','M3 12h4','M17 12h4','m6 6 2.5 2.5','m15.5 15.5 2.5 2.5','m18 6-2.5 2.5','m8.5 15.5L6 18'],
  lightbulb:['M9 18h6','M10 21.5h4','M12 2.5a6 6 0 0 0-3.6 10.8c.6.5 1 1.2 1.1 2h5c.1-.8.5-1.5 1.1-2A6 6 0 0 0 12 2.5Z'],
  clipboard:[R(8,3,8,4,1.5),'M16 5h1.5A2 2 0 0 1 19.5 7v12.5a2 2 0 0 1-2 2h-11a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2H8'],
  send:['m22 2-7.5 20-4-9-9-4z','M22 2 10.5 13'],
  home:['m3 10.5 9-7.5 9 7.5','M5.5 9.5V21h13V9.5','M10 21v-6h4v6'],
  info:[C(12,12,10),'M12 16v-4.5','M12 8h.01'],
  gauge:['M12 15V9.5','M3.5 18a9 9 0 1 1 17 0'],
  arrowUpRight:['M7 17 17 7','M8 7h9v9'],
};
function Icon({name, className='w-5 h-5', sw=1.75}){
  const d = PATHS[name] || [];
  return (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={sw}
         strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      {d.map((it,i)=> typeof it === 'string'
        ? <path key={i} d={it}/>
        : it.t==='c' ? <circle key={i} cx={it.x} cy={it.y} r={it.r}/>
        : <rect key={i} x={it.x} y={it.y} width={it.w} height={it.h} rx={it.rx}/> )}
    </svg>
  );
}

/* ───────────────────────── LOGO ───────────────────────── */

export default Icon;
