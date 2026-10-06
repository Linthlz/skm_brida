function Card({children, className='', pad='p-5 sm:p-6'}){
  return <div className={`bg-surface border border-line rounded-card min-w-0 ${pad} ${className}`}>{children}</div>;
}

export default Card;
