import { useInfiniteQuery, useQuery } from '@connectrpc/connect-query'
import { useEffect, useRef } from 'react'

import { getJobLogs } from '@/gen/proto/services/web-WebService_connectquery'

import { LoadingBlock } from '@/src/components/loading_block'

import { Log } from './log'

const HISTORY_PAGE_LIMIT = 30
const NEWER_LOGS_REFETCH_INTERVAL = 5000

export function History() {
    const getOlderLogsQuery = useInfiniteQuery(
        getJobLogs,
        { limit: HISTORY_PAGE_LIMIT, olderThanLogId: '' },
        {
            pageParamKey: 'olderThanLogId',
            getNextPageParam: (lastPage) =>
                lastPage.ids[lastPage.ids.length - 1],
        }
    )

    const getNewerLogsQuery = useQuery(
        getJobLogs,
        { newerThanLogId: getOlderLogsQuery.data?.pages[0]?.ids[0] || '0' },
        {
            enabled: getOlderLogsQuery.isSuccess,
            refetchInterval: NEWER_LOGS_REFETCH_INTERVAL,
        }
    )

    const scrollContainerRef = useRef<HTMLElement>(null)

    useEffect(() => {
        const scrollContainer = scrollContainerRef.current
        if (!scrollContainer) return
        if (getOlderLogsQuery.isPending) return

        const onScroll = () => {
            if (
                scrollContainer.scrollHeight - scrollContainer.scrollTop <=
                scrollContainer.clientHeight + 30
            ) {
                getOlderLogsQuery.fetchNextPage()
            }
        }

        scrollContainer.addEventListener('scroll', onScroll)
        return () => scrollContainer.removeEventListener('scroll', onScroll)
    }, [
        scrollContainerRef.current,
        getOlderLogsQuery.isSuccess,
        getNewerLogsQuery.isSuccess,
    ])

    if (getOlderLogsQuery.isError || getNewerLogsQuery.isError) {
        return (
            <section className="w-96 flex flex-col rounded-xl bg-neutral p-4">
                <h2 className="text-lg">Job History</h2>
                <section className="h-full flex items-center justify-center">
                    <p className="text-center">
                        {getOlderLogsQuery.isError &&
                            getOlderLogsQuery.error.message}
                        {getNewerLogsQuery.isError &&
                            getNewerLogsQuery.error.message}
                    </p>
                </section>
            </section>
        )
    }

    if (getOlderLogsQuery.isPending || getNewerLogsQuery.isPending) {
        return (
            <section className="w-96 flex flex-col rounded-xl bg-neutral p-4">
                <h2 className="text-lg">Job History</h2>
                <section className="h-full flex flex-col items-center justify-center gap-8">
                    <LoadingBlock className="bg-neutral-light size-16" />
                    <LoadingBlock className="bg-neutral-light size-16" />
                    <LoadingBlock className="bg-neutral-light size-16" />
                </section>
            </section>
        )
    }

    return (
        <section className="w-96 flex flex-col rounded-xl bg-neutral p-4 overflow-hidden">
            <h2 className="text-lg">Job History</h2>
            <section
                className="h-full overflow-y-auto mt-2"
                ref={scrollContainerRef}
            >
                <section className="flex flex-col gap-2 w-full h-fit">
                    {getNewerLogsQuery.data.ids.map((logId) => {
                        return <Log id={logId} key={logId} />
                    })}
                    {getOlderLogsQuery.data.pages.flatMap((page) =>
                        page.ids.map((logId) => {
                            return <Log id={logId} key={logId} />
                        })
                    )}
                </section>
            </section>
        </section>
    )
}
