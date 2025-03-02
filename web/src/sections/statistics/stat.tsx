type StatOption = {
    name: string
    value: number
    total: number
}

export function Stat({ name, value, total }: StatOption) {
    return (
        <section className="flex w-full max-w-80 items-center justify-center">
            <div>
                <p className="text-lg">{name}</p>
                <p className="text-6xl font-light">
                    {value} / {total}
                </p>
            </div>
        </section>
    )
}
