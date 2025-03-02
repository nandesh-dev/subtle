import { LoadingBlock } from '@/src/components/loading_block'
import { useQuery } from '@connectrpc/connect-query'

import { calculateSubtitleStatistics } from '@/gen/proto/services/web-WebService_connectquery'

import { Stat } from './stat'

const UPDATE_INTERVAL = 10 * 1000

export function Statistics() {
    const calculateSubtitleStatisticsQuery = useQuery(
        calculateSubtitleStatistics,
        {},
        { refetchInterval: UPDATE_INTERVAL }
    )

    if (calculateSubtitleStatisticsQuery.isError) {
        return (
            <section className="h-64 flex flex-row justify-center items-center">
                {calculateSubtitleStatisticsQuery.error.message}
            </section>
        )
    }

    if (calculateSubtitleStatisticsQuery.isPending) {
        return (
            <section className="h-64 flex flex-row justify-center items-center gap-16">
                <LoadingBlock className="bg-neutral size-16"/>
                <LoadingBlock className="bg-neutral size-16"/>
                <LoadingBlock className="bg-neutral size-16"/>
            </section>
        )
    }

    return (
        <section className="flex flex-row justify-center h-64">
            <Stat
                name="Exported"
                value={
                    calculateSubtitleStatisticsQuery.data
                        .videoWithExportedSubtitleCount
                }
                total={calculateSubtitleStatisticsQuery.data.totalVideoCount}
            />
            <Stat
                name="Formated"
                value={
                    calculateSubtitleStatisticsQuery.data
                        .videoWithFormatedSubtitleCount
                }
                total={calculateSubtitleStatisticsQuery.data.totalVideoCount}
            />
            <Stat
                name="Extracted"
                value={
                    calculateSubtitleStatisticsQuery.data
                        .videoWithExtractedSubtitleCount
                }
                total={calculateSubtitleStatisticsQuery.data.totalVideoCount}
            />
        </section>
    )
}
