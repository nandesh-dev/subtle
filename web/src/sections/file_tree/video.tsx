import filepath from 'path-browserify'

import { Route, useNavigation } from '@/src/utility/navigation'

export function Video({ path }: { path: string }) {
    const navigation = useNavigation()

    const pathSearchParam = navigation.useSearchParams().get('path')
    const isSelected = !!pathSearchParam && pathSearchParam == path

    const onClick = () => {
        if (isSelected) {
            navigation.navigate(
                Route.Files,
                new URLSearchParams({ path: filepath.dirname(path) })
            )
            return
        }
        navigation.navigate(Route.Files, new URLSearchParams({ path }))
    }

    return (
        <section className="group grid grid-cols-[auto_1fr] gap-1">
            <div
                className={`mx-2 w-1 rounded-full ${isSelected ? 'bg-primary' : 'bg-primary-light'} group-hover:bg-primary transition-colors`}
            />
            <button className="text-start text-sm" onClick={onClick}>
                {filepath.basename(path)}
            </button>
        </section>
    )
}
