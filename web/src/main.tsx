import { useQuery } from '@connectrpc/connect-query'
import { Routes, useNavigation } from './utils/navigation'
import { Style } from './utils/style'
import { getGlobalStatistics } from '../gen/proto/web/web-WebService_connectquery'

export function Main() {
    return (
        <section className="flex flex-col gap-md p-xl">
            <section className="flex flex-row justify-between">
                <section>
                    <h1 className="text-xl text-text-1">Subtle</h1>
                </section>
                <NavigationBar />
            </section>
            <Stats />
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
        <section className="rounded-xl bg-neutral-2">
            {NavigationButtons.map(({ route, name }) => {
                return (
                    <button
                        className="rounded-xl px-xl py-md text-lg"
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
    const { data: globalStatistics } = useQuery(getGlobalStatistics, {})

    const Stat = (data: { name: string; value?: number; total?: number }) => {
        const value = data.value != undefined ? data.value : '..'
        const total = data.total != undefined ? data.total : '..'

        return (
            <div className="flex flex-col">
                <p className="text-lg text-text-1">{data.name}</p>
                <p className="text-2xl font-light text-text-1">
                    {value}/{total}
                </p>
            </div>
        )
    }

    return (
        <section className="flex flex-row justify-center gap-2xl p-2xl">
            <Stat
                name="Exported"
                value={globalStatistics?.Exported}
                total={globalStatistics?.Total}
            />
            <Stat
                name="Formated"
                value={globalStatistics?.Formated}
                total={globalStatistics?.Total}
            />
            <Stat
                name="Extracted"
                value={globalStatistics?.Extracted}
                total={globalStatistics?.Total}
            />
        </section>
    )
}
