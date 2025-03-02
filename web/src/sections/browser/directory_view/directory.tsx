import filepath from 'path-browserify'

import { Route, useNavigation } from '@/src/utility/navigation'

type DirectoryOptions = {
    path: string
    displayEntirePath?: boolean
}

export function Directory({ path, displayEntirePath }: DirectoryOptions) {
    const navigation = useNavigation()

    const navigateToDirectory = () => {
        navigation.navigate(Route.Files, new URLSearchParams({ path }))
    }

    return (
        <button
            className="group h-fit min-w-48 flex flex-col bg-neutral-light p-4 gap-4 rounded-lg"
            onClick={navigateToDirectory}
        >
            <section className="h-fit flex flex-row items-center pr-8 gap-4">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    className="size-10 fill-secondary-light group-hover:fill-secondary transition-colors"
                >
                    <path d="M2.4 21.6C1.74 21.6 1.175 21.365 0.705 20.895C0.235 20.425 0 19.86 0 19.2V4.80002C0 4.14002 0.235 3.57502 0.705 3.10502C1.175 2.63502 1.74 2.40002 2.4 2.40002H9.6L12 4.80002H21.6C22.26 4.80002 22.825 5.03502 23.295 5.50502C23.765 5.97502 24 6.54002 24 7.20002V19.2C24 19.86 23.765 20.425 23.295 20.895C22.825 21.365 22.26 21.6 21.6 21.6H2.4Z" />
                </svg>
                <p className="text-nowrap">
                    {displayEntirePath ? path : filepath.basename(path)}
                </p>
            </section>
        </button>
    )
}
