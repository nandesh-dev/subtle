import { CSSProperties } from 'react'

export function TickIcon({
    className,
    style,
}: {
    className?: string
    style?: CSSProperties
}) {
    return (
        <svg
            width="17"
            height="13"
            viewBox="0 0 17 13"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className={className}
            style={style}
        >
            <path
                d="M5.17253 10.1625L1.70253 6.6925C1.31253 6.3025 0.682531 6.3025 0.292531 6.6925C-0.0974695 7.0825 -0.0974695 7.7125 0.292531 8.1025L4.47253 12.2825C4.86253 12.6725 5.49253 12.6725 5.88253 12.2825L16.4625 1.7025C16.8525 1.3125 16.8525 0.6825 16.4625 0.2925C16.0725 -0.0975 15.4425 -0.0975 15.0525 0.2925L5.17253 10.1625Z"
                fill="#838683"
            />
        </svg>
    )
}
