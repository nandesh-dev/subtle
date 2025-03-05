import { useQuery } from '@connectrpc/connect-query'

import {
    SubtitleOriginalFormat,
    SubtitleStage,
} from '@/gen/proto/messages/subtitle_pb'
import { getSubtitle } from '@/gen/proto/services/web-WebService_connectquery'

import { LoadingBlock } from '@/src/components/loading_block'

import { Route, useNavigation } from '@/src/utility/navigation'

export function Subtitle({ id }: { id: string }) {
    const navigation = useNavigation()

    const getSubtitleQuery = useQuery(getSubtitle, { id })

    if (getSubtitleQuery.isError) {
        return (
            <section className="h-32 flex bg-neutral-light rounded-lg p-4 items-center justify-center">
                <p>{getSubtitleQuery.error.message}</p>
            </section>
        )
    }

    if (getSubtitleQuery.isPending) {
        return (
            <section className="h-32 flex bg-neutral-light rounded-lg p-4 items-center justify-center flex-row gap-4">
                <LoadingBlock className="bg-neutral-lighter size-8" />
                <LoadingBlock className="bg-neutral-lighter size-8" />
                <LoadingBlock className="bg-neutral-lighter size-8" />
            </section>
        )
    }

    let progress = 0
    switch (getSubtitleQuery.data.stage) {
        case SubtitleStage.DETECTED:
            progress = 1
            break
        case SubtitleStage.EXTRACTED:
            progress = 3
            break
        case SubtitleStage.FORMATED:
            progress = 5
            break
        case SubtitleStage.EXPORTED:
            progress = 7
            break
    }
    if (getSubtitleQuery.data.isProcessing) {
        progress += 1
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

    const language = getSubtitleQuery.data.language || 'unknown language'

    const isExtracting =
        getSubtitleQuery.data.isProcessing &&
        getSubtitleQuery.data.stage == SubtitleStage.DETECTED
    const isFormating =
        getSubtitleQuery.data.isProcessing &&
        getSubtitleQuery.data.stage == SubtitleStage.EXTRACTED
    const isExporting =
        getSubtitleQuery.data.isProcessing &&
        getSubtitleQuery.data.stage == SubtitleStage.FORMATED

    const onClick = () => {
        navigation.navigate(
            Route.Editor,
            new URLSearchParams({
                subtitle_id: id.toString(),
            })
        )
    }

    return (
        <button
            className="group flex flex-col bg-neutral-light rounded-lg py-4"
            onClick={onClick}
        >
            <section className="w-full flex flex-row justify-between px-4">
                <p>{getSubtitleQuery.data.title}</p>
                <section className="flex flex-row gap-8">
                    <p>{language}</p>
                    <p>{format}</p>
                </section>
            </section>
            <div className="h-px bg-neutral w-full my-4"/>
            <section className="w-full grid grid-cols-[auto_1fr_auto_1fr_auto_1fr_auto] items-center gap-2 px-6 py-2">
                <div
                    className={`bg-neutral-light size-4 rounded-full border-4 ${progress >= 1 ? 'border-primary-light' : 'border-neutral-lighter'}`}
                />
                <div
                    className={`h-1 rounded-full ${progress >= 2 ? 'bg-primary-light' : 'bg-neutral-lighter'} ${isExtracting && 'animate-pulse'}`}
                />
                <div
                    className={`bg-neutral-light size-4 rounded-full border-4 ${progress >= 3 ? 'border-primary-light' : 'border-neutral-lighter'}`}
                />
                <div
                    className={`h-1 rounded-full ${progress >= 4 ? 'bg-primary-light' : 'bg-neutral-lighter'} ${isFormating && 'animate-pulse'}`}
                />
                <div
                    className={`bg-neutral-light size-4 rounded-full border-4 ${progress >= 5 ? 'border-primary-light' : 'border-neutral-lighter'}`}
                />
                <div
                    className={`h-1 rounded-full ${progress >= 6 ? 'bg-primary-light' : 'bg-neutral-lighter'} ${isExporting && 'animate-pulse'}`}
                />
                <div
                    className={`bg-neutral-light size-4 rounded-full border-4 ${progress >= 7 ? 'border-primary-light' : 'border-neutral-lighter'}`}
                />
            </section>
            <section className="w-full grid grid-cols-[6fr_11fr_11fr_6fr] px-4">
                <p className="text-sm text-start">
                    {getSubtitleQuery.data.importIsExternal
                        ? 'Imported'
                        : 'Detected'}
                </p>
                <p className="text-sm text-center">Extracted</p>
                <p className="text-sm text-center">Formated</p>
                <p className="text-sm text-end">Exported</p>
            </section>
        </button>
    )
}
