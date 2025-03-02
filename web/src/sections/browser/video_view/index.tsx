import { useQuery } from '@connectrpc/connect-query'

import { searchVideo } from '@/gen/proto/services/web-WebService_connectquery'

import { useSearchParams } from '@/src/utility/navigation'

import { Subtitle } from './subtitle'

export function VideoView() {
    const pathSearchParam = useSearchParams().get('path')

    const getVideoQuery = useQuery(searchVideo, { path: pathSearchParam || '' })

    return (
        <div className="relative flex flex-col gap-lg rounded-md bg-neutral-2 p-xl">
            <section className="flex flex-col gap-md">
                <h2 className="text-lg text-text-1">Subtitles</h2>
            </section>
            <section className="flex flex-col gap-md">
                {!getVideoQuery.isSuccess ? (
                    <>
                        <div className="h-2xl animate-pulse rounded-sm bg-neutral-1" />
                        <div className="h-2xl animate-pulse rounded-sm bg-neutral-1" />
                    </>
                ) : (
                    getVideoQuery.data.subtitleIds.map((subtitleId) => {
                        return <Subtitle key={subtitleId} id={subtitleId} />
                    })
                )}
            </section>
            <div className="absolute bottom-0 right-0 rounded-tl-md bg-neutral-2 p-md">
                {!getVideoQuery.isSuccess ? (
                    <div className="h-lg w-[24rem] animate-pulse rounded-sm bg-neutral-1" />
                ) : (
                    <p className="text-xs text-text-1"> {pathSearchParam}</p>
                )}
            </div>
        </div>
    )
}
