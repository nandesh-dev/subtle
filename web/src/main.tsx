import { useQuery } from '@connectrpc/connect-query'
import { Routes, useNavigation } from './utils/navigation'
import { Style } from './utils/style'
import { getGlobalStatistics } from '../gen/proto/web/web-WebService_connectquery'
import { Editor, Files, Jobs, Settings } from './section'
import { useIsMutating } from '@tanstack/react-query'

export function Main() {
    const navigation = useNavigation()
    const currentRoute = navigation.useRoute()

    const isMutating = useIsMutating()

    return (
        <section className="grid h-dvh w-dvw grid-rows-[auto_auto_1fr] gap-md p-xl">
            <section className="flex max-w-full flex-row items-center gap-xl">
                <section>
                    <h1 className="text-xl text-text-1">Subtle</h1>
                </section>
                {isMutating ? (
                    <div className="h-xs w-full animate-loader rounded-sm bg-gradient-to-r from-primary-2 to-primary-2 bg-[length:60%_100%] bg-no-repeat" />
                ) : (
                    <div className="h-xs w-full" />
                )}
                <NavigationBar />
            </section>
            <Stats />
            {(currentRoute == Routes.Files ||
                currentRoute == Routes.Editor) && <Files />}
            {currentRoute == Routes.Settings && <Settings />}
            {currentRoute == Routes.Jobs && <Jobs />}
            {currentRoute == Routes.Editor && (
                <section className="absolute bottom-0 left-0 right-0 top-0 bg-[rgba(0,0,0,0.1)] p-xl">
                    <Editor />
                </section>
            )}
        </section>
    )
}

const NavigationButtons = [
    {
        route: Routes.Files,
        name: 'Files',
    },
    {
        route: Routes.Settings,
        name: 'Settings',
    },
    {
        route: Routes.Jobs,
        name: 'Jobs',
    },
]

function NavigationBar() {
    const navigation = useNavigation()
    const currentRoute = navigation?.useRoute()

    return (
        <section className="flex flex-row rounded-xl bg-neutral-2">
            {NavigationButtons.map(({ route, name }) => {
                return (
                    <button
                        className="rounded-xl px-xl py-md text-md"
                        style={
                            route == currentRoute
                                ? {
                                      background: Style.colors.primary[1],
                                      color: Style.colors.text[2],
                                  }
                                : { color: Style.colors.text[1] }
                        }
                        key={name}
                        onClick={() => {
                            navigation?.navigate(route)
                        }}
                    >
                        {name}
                    </button>
                )
            })}
        </section>
    )
}

function Stats() {
    const { data: globalStatistics, isLoading } = useQuery(
        getGlobalStatistics,
        {}
    )

    const Stat = (data: { name: string; value?: number }) => {
        if (isLoading) {
            return (
                <div className="h-[8rem] w-full max-w-[16rem] animate-pulse rounded-sm bg-neutral-2" />
            )
        }

        return (
            <div className="flex h-[8rem] w-full max-w-[16rem] flex-col">
                <p className="text-lg text-text-1">{data.name}</p>
                <p className="text-2xl font-light text-text-1">
                    {data?.value}/{globalStatistics?.Total}
                </p>
            </div>
        )
    }

    return (
        <section className="flex flex-row justify-center gap-2xl py-2xl">
            <Stat name="Exported" value={globalStatistics?.Exported} />
            <Stat name="Formated" value={globalStatistics?.Formated} />
            <Stat name="Extracted" value={globalStatistics?.Extracted} />
        </section>
    )
}
