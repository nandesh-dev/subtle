type LoadingBlockOptions = {
    className: string
}

export function LoadingBlock({ className }: LoadingBlockOptions) {
    return (
        <div className={`motion-safe:animate-pulse rounded-md ${className}`} />
    )
}
