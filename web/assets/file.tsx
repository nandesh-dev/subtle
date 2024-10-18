import { CSSProperties } from 'react'

export function FileIcon({
    className,
    style,
}: {
    className?: string
    style?: CSSProperties
}) {
    return (
        <svg
            width="23"
            height="18"
            viewBox="0 0 23 18"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className={className}
            style={style}
        >
            <path d="M2.3 18C1.6675 18 1.12623 17.7799 0.6762 17.3396C0.2254 16.8986 0 16.3688 0 15.75V2.25C0 1.63125 0.2254 1.10175 0.6762 0.6615C1.12623 0.2205 1.6675 0 2.3 0H8.25125C8.55792 0 8.8504 0.0562501 9.1287 0.16875C9.40623 0.28125 9.65042 0.440625 9.86125 0.646875L11.5 2.25H20.7C21.3325 2.25 21.8741 2.4705 22.3249 2.9115C22.775 3.35175 23 3.88125 23 4.5V15.75C23 16.3688 22.775 16.8986 22.3249 17.3396C21.8741 17.7799 21.3325 18 20.7 18H2.3ZM2.3 2.25V15.75H20.7V4.5H10.5513L8.25125 2.25H2.3Z" />
        </svg>
    )
}
