import { History } from './history'
import { Status } from './status'

export function Jobs() {
    return (
        <section className="h-full grid grid-cols-[1fr_auto] overflow-hidden gap-4">
            <Status />
            <History />
        </section>
    )
}
