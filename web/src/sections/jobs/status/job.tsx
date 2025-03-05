import { useQuery } from '@connectrpc/connect-query'

import { getJob } from '@/gen/proto/services/web-WebService_connectquery'

type JobOptions = {
    code: string
}

export function Job({ code }: JobOptions) {
    const getJobQuery = useQuery(getJob, { code })

    if (!getJobQuery.isSuccess) return

    const lastRun = new Date(Number(getJobQuery.data.lastRun?.seconds) * 1000)

    return (
        <section className="group grid grid-cols-[auto_auto_1fr_auto] bg-neutral-light rounded-lg">
            <section className="w-24 flex items-center justify-center">
                <p className="text-4xl">{getJobQuery.data.sequenceNumber}</p>
            </section>
            <div className="bg-neutral w-px h-full" />
            <section className="flex flex-col justify-between gap-4 p-4">
                <section>
                    <h3 className="text-lg mb-1">{getJobQuery.data.name}</h3>
                    <p>{getJobQuery.data.description}</p>
                </section>
                <p className="text-sm">Last Run: {lastRun.toLocaleString()}</p>
            </section>
            <section className="flex flex-col justify-end p-4">
                <button
                    className={`w-32 text-center py-2 rounded-md ${getJobQuery.data.isRunning ? 'disabled:bg-secondary-light' : 'bg-primary hover:bg-primary-light disabled:bg-primary-light'}`}
                    disabled={getJobQuery.data.isRunning}
                >
                    {getJobQuery.data.isRunning ? 'Running' : 'Run'}
                </button>
            </section>
        </section>
    )
}
