import {
    ReactNode,
    createContext,
    useContext,
    useEffect,
    useState,
} from 'react'

export enum Route {
    Files = '',
    Settings = 'settings',
    Jobs = 'jobs',
    Editor = 'editor',
    Unknown = 'unknown',
}

function mapPathnameToRoute(pathname: string) {
    switch (pathname) {
        case `/${Route.Files}`:
            return Route.Files
        case `/${Route.Settings}`:
            return Route.Settings
        case `/${Route.Jobs}`:
            return Route.Jobs
        case `/${Route.Editor}`:
            return Route.Editor
        default:
            return Route.Unknown
    }
}

//TODO Allow opening it new tab while pressing ctrl

const HomeRoute = Route.Files

type EventListener = (route: Route, searchParams: URLSearchParams) => void

export class Navigation {
    private eventListeners: EventListener[]
    private routeLastSearchParamsMap: Map<Route, URLSearchParams>
    constructor() {
        this.eventListeners = []
        this.routeLastSearchParamsMap = new Map()

        window.addEventListener('popstate', () => {
            this.updateEventListeners()
        })
    }

    public back() {
        window.history.back()
        this.updateEventListeners()
    }

    public navigate(route: Route, searchParams?: URLSearchParams) {
        const oldURL = new URL(window.location.href)
        this.routeLastSearchParamsMap.set(
            mapPathnameToRoute(oldURL.pathname),
            oldURL.searchParams
        )

        const newURL = new URL(window.location.href)

        newURL.pathname = `/${route}`
        newURL.search = searchParams?.toString() || ''

        window.history.pushState({}, '', newURL)

        this.updateEventListeners()
    }

    public getRouteLastSearchParams(route: Route) {
        return this.routeLastSearchParamsMap.get(route) || new URLSearchParams()
    }

    public useRoute() {
        const route = this.parseRoute()

        const [state, setState] = useState(route)

        useEffect(() => {
            const eventListener: EventListener = (route, _) => {
                setState(route)
            }

            this.eventListeners.push(eventListener)

            return () => {
                this.eventListeners.splice(
                    this.eventListeners.indexOf(eventListener),
                    1
                )
            }
        }, [])

        return state
    }

    public useSearchParams() {
        const searchParams = this.parseSearchParams()

        const [state, setState] = useState(searchParams)

        useEffect(() => {
            const eventListener: EventListener = (_, searchParams) => {
                setState(searchParams)
            }

            this.eventListeners.push(eventListener)

            return () => {
                this.eventListeners.splice(
                    this.eventListeners.indexOf(eventListener),
                    1
                )
            }
        }, [])

        return state
    }

    private updateEventListeners() {
        const route = this.parseRoute()
        const searchParams = this.parseSearchParams()

        for (let eventListener of this.eventListeners) {
            eventListener(route, searchParams)
        }
    }

    private parseSearchParams() {
        const searchParams = new URLSearchParams(window.location.search)
        return searchParams
    }

    private parseRoute() {
        const url = new URL(window.location.href)
        const pathname = url.pathname.replace(/^\//, '')

        if (Object.values(Route).includes(pathname as Route)) {
            return pathname as Route
        }

        return HomeRoute
    }
}

const NavigationContext = createContext<Navigation | null>(null)

export function NavigationProvider({
    navigation,
    children,
}: {
    navigation: Navigation
    children?: ReactNode
}) {
    return <NavigationContext.Provider value={navigation} children={children} />
}

export function useNavigation() {
    let navigation = useContext(NavigationContext)
    if (navigation == null) {
        throw new Error('Navigation is not defined yet!')
    }

    return navigation
}

export function useRoute() {
    const navigation = useNavigation()
    return navigation.useRoute()
}

export function useSearchParams() {
    const navigation = useNavigation()
    return navigation.useSearchParams()
}
