import { Status } from './status'

export function Jobs() {
    return (
        <section className="h-full grid grid-cols-[1fr_auto] overflow-hidden gap-4">
            <Status />
            <section className="w-96 flex flex-col rounded-xl bg-neutral p-4">
                <h2 className="text-lg">Job History</h2>
            </section>
        </section>
    )
}
