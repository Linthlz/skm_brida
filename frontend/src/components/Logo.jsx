/**
 * tone="nav" dipakai di atas permukaan merah (navbar & footer): perisai jadi
 * putih dan detailnya merah, supaya lambang tidak lebur dengan latarnya.
 */
function Logo({className='w-10 h-10', tone='default'}){
  const nav = tone==='nav';
  const perisai = nav ? 'var(--nav-ink)' : 'var(--brand)';
  const garis   = nav ? 'var(--nav)' : '#fff';
  const aksen   = nav ? 'var(--nav-hover)' : 'var(--gold)';
  return (
    <svg className={className} viewBox="0 0 48 48" role="img" aria-label="Lambang BRIDA Kabupaten Buleleng">
      <path d="M24 2.5 43 9.6v15.2c0 11.1-8 17.9-19 21.7C13 42.7 5 35.9 5 24.8V9.6z" fill={perisai}/>
      <path d="M24 6.2 39.4 12v12.8c0 9.2-6.6 15-15.4 18.2C15.2 39.8 8.6 34 8.6 24.8V12z" fill="none" stroke={aksen} strokeWidth="1.4"/>
      <path d="M15 30c1.6-4.4 4.8-7 9-7s7.4 2.6 9 7" fill="none" stroke={garis} strokeWidth="2.1" strokeLinecap="round"/>
      <circle cx="24" cy="17.5" r="4.2" fill={aksen}/>
      <path d="M17.5 34.5h13" stroke={garis} strokeWidth="2.1" strokeLinecap="round"/>
    </svg>
  );
}

export default Logo;
