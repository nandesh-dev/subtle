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
}

//TODO Allow opening it new tab while pressing ctrl

const HomeRoute = Route.Files

type EventListener = (route: Route, searchParams: URLSearchParams) => void

export class Navigation {
    private eventListeners: EventListener[]
    constructor() {
        this.eventListeners = []

        window.addEventListener('popstate', () => {
            this.updateEventListeners()
        })
    }

    back() {
        window.history.back()
        this.updateEventListeners()
    }

    navigate(route: Route, searchParams?: URLSearchParams) {
        const newURL = new URL(window.location.href)

        newURL.pathname = `/${route}`
        newURL.search = searchParams?.toString() || ''

        window.history.pushState({}, '', newURL)

        this.updateEventListeners()
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
