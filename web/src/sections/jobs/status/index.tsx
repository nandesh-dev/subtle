import { useQuery } from '@connectrpc/connect-query'

import { getJobs } from '@/gen/proto/services/web-WebService_connectquery'
import { LoadingBlock } from '@/src/components/loading_block'
import { Job } from './job'

export function Status() {
    const getJobsQuery = useQuery(getJobs)

    if (getJobsQuery.isError) {
        return (
            <section className="flex flex-col rounded-xl bg-neutral p-4">
                <h2 className="text-lg">Job Status</h2>
                <section className="h-full flex items-center justify-center">
                    <p>{getJobsQuery.error.message}</p>
                </section>
            </section>
        )
    }

    if (getJobsQuery.isPending) {
        return (
            <section className="flex flex-col rounded-xl bg-neutral p-4">
                <h2 className="text-lg">Job Status</h2>
                <section className="h-full flex flex-row items-center justify-center gap-8">
                    <LoadingBlock className="bg-neutral-light size-16" />
                    <LoadingBlock className="bg-neutral-light size-16" />
                    <LoadingBlock className="bg-neutral-light size-16" />
                </section>
            </section>
        )
    }

    return (
        <section className="flex flex-col rounded-xl bg-neutral p-4 overflow-hidden">
            <h2 className="text-lg">Job Status</h2>
                <section className="h-full overflow-y-auto mt-2">
                    <section className="flex flex-col gap-2 w-full h-fit">
                        {getJobsQuery.data.jobCodes.map((jobCode) => {
                            return <Job code={jobCode} key={jobCode} />
                        })}
                    </section>
                </section>
        </section>
    )
}
