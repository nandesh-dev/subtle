import { useQuery } from '@connectrpc/connect-query'

import { getJobLog } from '@/gen/proto/services/web-WebService_connectquery'

import { LoadingBlock } from '@/src/components/loading_block'

type LogOptions = {
    id: string
}

export function Log({ id }: LogOptions) {
    const getLogQuery = useQuery(getJobLog, { id })

    if (getLogQuery.isError) {
        return (
            <section className="h-24 rounded-lg bg-neutral-light p-4 flex items-center justify-center">
                <p className="text-center">{getLogQuery.error.message}</p>
            </section>
        )
    }

    if (getLogQuery.isPending) {
        return (
            <section className="h-24 rounded-lg bg-neutral-light p-4 flex flex-row gap-4 items-center justify-center">
                <LoadingBlock className="bg-neutral-lighter size-8" />
                <LoadingBlock className="bg-neutral-lighter size-8" />
                <LoadingBlock className="bg-neutral-lighter size-8" />
            </section>
        )
    }

    const startTimestamp = new Date(
        Number(getLogQuery.data.startTimestamp?.seconds) * 1000
    )

    return (
        <section className="min-h-24 rounded-lg bg-neutral-light py-4 flex flex-col justify-between gap-2">
            <section className="flex flex-row justify-between px-4">
                <p>{getLogQuery.data.jobName}</p>
                <div className="size-2 rounded-full bg-primary" />
            </section>
            <section>
                <div className="h-px bg-neutral px-4 mb-4" />
                <section className="flex flex-row justify-between px-4">
                    <p className="text-sm">{startTimestamp.toLocaleString()}</p>
                    <p className="text-sm">
                        {formatDuration(
                            Number(getLogQuery.data.duration?.seconds)
                        )}
                    </p>
                </section>
            </section>
        </section>
    )
}

function formatDuration(seconds: number) {
    const minutes = Math.floor(seconds / 60) + 'm'
    const secs = (seconds % 60) + 's'

    return [minutes, secs].join(' ')
}
