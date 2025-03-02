import { useQuery } from '@connectrpc/connect-query'

import {
    SubtitleOriginalFormat,
    SubtitleStage,
} from '@/gen/proto/messages/subtitle_pb'
import { getSubtitle } from '@/gen/proto/services/web-WebService_connectquery'

import { Route, useNavigation } from '@/src/utility/navigation'

export function Subtitle({ id }: { id: string }) {
    const navigation = useNavigation()

    const getSubtitleQuery = useQuery(getSubtitle, { id })

    if (!getSubtitleQuery.isSuccess) {
        return <div className="h-2xl animate-pulse rounded-sm bg-neutral-1" />
    }

    let status = 'unknown'
    switch (getSubtitleQuery.data.stage) {
        case SubtitleStage.DETECTED:
            status = 'detected'
            break
        case SubtitleStage.EXTRACTED:
            status = 'extracted'
            break
        case SubtitleStage.FORMATED:
            status = 'formated'
            break
        case SubtitleStage.EXPORTED:
            status = 'exported'
            break
    }
    if (getSubtitleQuery.data.isProcessing) {
        switch (getSubtitleQuery.data.stage) {
            case SubtitleStage.DETECTED:
                status = 'extracting'
                break
            case SubtitleStage.EXTRACTED:
                status = 'formating'
                break
            case SubtitleStage.FORMATED:
                status = 'exporting'
                break
        }
    }

    let format = ''
    switch (getSubtitleQuery.data.originalFormat) {
        case SubtitleOriginalFormat.SRT:
            format = 'srt'
            break
        case SubtitleOriginalFormat.ASS:
            format = 'ass'
            break
        case SubtitleOriginalFormat.PGS:
            format = 'pgs'
            break
    }

    return (
        <button
            className="group grid min-h-2xl grid-cols-[0.25rem_1fr] items-center gap-md rounded-sm bg-neutral-1 p-md"
            onClick={() => {
                navigation?.navigate(
                    Route.Editor,
                    new URLSearchParams({
                        subtitle_id: id.toString(),
                    })
                )
            }}
        >
            <div className="h-full rounded-sm bg-primary-2 group-hover:bg-primary-1" />
            <section className="grid grid-flow-row auto-rows-fr gap-sm lg:auto-cols-fr lg:grid-flow-col">
                <p className="text-sm text-text-1">
                    {getSubtitleQuery.data.title}
                </p>
                <section className="grid auto-cols-fr grid-flow-col">
                    <p className="text-center text-xs text-text-1">
                        {getSubtitleQuery.data.language}
                    </p>
                    <p className="text-center text-xs text-text-1">{format}</p>
                    <p className="text-center text-xs text-text-1">{status}</p>
                </section>
            </section>
        </button>
    )
}
