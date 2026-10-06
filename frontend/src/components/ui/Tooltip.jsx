function Tooltip({x, y, children, w=200}){
  return <div className="pointer-events-none absolute z-30 bg-[#5A1010] text-white text-[12.5px] border border-white/15 leading-snug rounded-[7px] px-2.5 py-2 shadow-pop"
              style={{left:Math.max(4,Math.min(x, 9999)), top:y, transform:'translate(-50%,-115%)', maxWidth:w}}>{children}</div>;
}

export default Tooltip;
