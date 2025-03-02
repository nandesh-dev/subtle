import filepath from 'path-browserify'

import { Route, useNavigation } from '@/src/utility/navigation'

type VideoOptions = { path: string }

export function Video({ path }: VideoOptions) {
    const navigation = useNavigation()

    const onClick = () => {
        navigation?.navigate(Route.Files, new URLSearchParams({ path: path }))
    }

    return (
        <button
            className="group grid grid-cols-[var(--spacing)_1fr] bg-neutral-light rounded-lg p-4 gap-4"
            onClick={onClick}
        >
            <div className="h-full rounded-full bg-primary-light group-hover:bg-primary transition-colors" />
            <p>{filepath.basename(path)}</p>
        </button>
    )
}
