import { CSSProperties } from 'react'

export function ClockIcon({
    className,
    style,
}: {
    className?: string
    style?: CSSProperties
}) {
    return (
        <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className={className}
            style={style}
        >
            <path d="M11.988 0C5.364 0 0 5.376 0 12C0 18.624 5.364 24 11.988 24C18.624 24 24 18.624 24 12C24 5.376 18.624 0 11.988 0ZM12 21.6C6.696 21.6 2.4 17.304 2.4 12C2.4 6.696 6.696 2.4 12 2.4C17.304 2.4 21.6 6.696 21.6 12C21.6 17.304 17.304 21.6 12 21.6ZM11.736 6H11.664C11.184 6 10.8 6.384 10.8 6.864V12.528C10.8 12.948 11.016 13.344 11.388 13.56L16.368 16.548C16.776 16.788 17.304 16.668 17.544 16.26C17.796 15.852 17.664 15.312 17.244 15.072L12.6 12.312V6.864C12.6 6.384 12.216 6 11.736 6Z" />
        </svg>
    )
}
