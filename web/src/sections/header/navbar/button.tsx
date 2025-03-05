import { Route, useNavigation } from '@/src/utility/navigation'

export function Button({ name, route }: { name: string; route: Route }) {
    const navigation = useNavigation()
    const currentRoute = navigation.useRoute()

    const isSelectedRouteButton = currentRoute == route

    const navigate = () => {
        navigation.navigate(route, navigation.getRouteLastSearchParams(route))
    }

    return (
        <button
            className={`py-4 px-6 rounded-full transition ${isSelectedRouteButton ? 'bg-primary-light hover:bg-primary text-tertiary-light' : 'hover:bg-neutral-light'}`}
            onClick={navigate}
        >
            {name}
        </button>
    )
}
