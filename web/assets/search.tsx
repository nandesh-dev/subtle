import { CSSProperties } from "react";

export function SearchIcon({
  className,
  style,
}: {
  className?: string;
  style?: CSSProperties;
}) {
  return (
    <svg
      width="19"
      height="19"
      viewBox="0 0 19 19"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
      style={style}
    >
      <path d="M13.2467 12.2936H12.4126L12.117 12.0085C13.384 10.5303 14.0386 8.51366 13.6796 6.37032C13.1834 3.4351 10.7338 1.09115 7.7775 0.732163C3.31132 0.183129 -0.447455 3.9419 0.101578 8.40808C0.460562 11.3644 2.80452 13.814 5.73974 14.3102C7.88308 14.6692 9.89972 14.0146 11.3779 12.7476L11.663 13.0432V13.8773L16.1503 18.3646C16.5832 18.7975 17.2906 18.7975 17.7235 18.3646C18.1563 17.9317 18.1563 17.2243 17.7235 16.7914L13.2467 12.2936ZM6.91171 12.2936C4.28268 12.2936 2.16046 10.1713 2.16046 7.5423C2.16046 4.91327 4.28268 2.79104 6.91171 2.79104C9.54074 2.79104 11.663 4.91327 11.663 7.5423C11.663 10.1713 9.54074 12.2936 6.91171 12.2936Z" />
    </svg>
  );
}
