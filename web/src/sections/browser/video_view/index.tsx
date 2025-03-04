import { useQuery } from '@connectrpc/connect-query'

import { searchVideo } from '@/gen/proto/services/web-WebService_connectquery'

import { LoadingBlock } from '@/src/components/loading_block'

import { useSearchParams } from '@/src/utility/navigation'

import { Subtitle } from './subtitle'

export function VideoView() {
    const pathSearchParam = useSearchParams().get('path')

    const getVideoQuery = useQuery(searchVideo, { path: pathSearchParam || '' })

    if (getVideoQuery.isError) {
        return (
            <section className="flex flex-col rounded-xl bg-neutral p-4">
                <h2 className="text-lg mb-2">Subtitles</h2>
                <section className=" flex items-center justify-center">
                    <p>{getVideoQuery.error.message}</p>
                </section>
            </section>
        )
    }

    if (getVideoQuery.isPending) {
        return (
            <section className="flex flex-col rounded-xl bg-neutral p-4">
                <h2 className="text-lg">Subtitles</h2>
                <section className="h-full flex items-center justify-center gap-8">
                    <LoadingBlock className="bg-neutral-light size-16" />
                    <LoadingBlock className="bg-neutral-light size-16" />
                    <LoadingBlock className="bg-neutral-light size-16" />
                </section>
            </section>
        )
    }

    return (
        <section className="flex min-h-32 flex-col rounded-xl bg-neutral p-4 overflow-hidden gap-2">
            <h2 className="text-lg">Subtitles</h2>
            {getVideoQuery.data.subtitleIds.length > 0 ? (
                <section className="h-full overflow-y-auto mt-2">
                    <section className="flex flex-col gap-2 w-full h-fit">
                        {getVideoQuery.data.subtitleIds.map((subtitleId) => {
                            return <Subtitle id={subtitleId} key={subtitleId} />
                        })}
                    </section>
                </section>
            ) : (
                <section className="h-full flex items-center justify-center">
                    <p className="text-sm">No subtitles found!</p>
                </section>
            )}
        </section>
    )
}
