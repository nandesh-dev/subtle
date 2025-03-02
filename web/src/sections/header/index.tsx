import { Route, useNavigation } from '@/src/utility/navigation'

import { Navbar } from './navbar'

export function Header() {
    const navigation = useNavigation()

    const navigateToHome = () => {
        navigation.navigate(Route.Files)
    }

    return (
        <section className="flex flex-row justify-between">
            <button onClick={navigateToHome}>
                <h1 className="text-3xl">Subtle</h1>
            </button>
            <Navbar />
        </section>
    )
}
