import { useQuery } from '@connectrpc/connect-query'
import { Routes, useNavigation } from '../utils/navigation'
import {
    getSubtitle,
    getSubtitleSegments,
} from '../../gen/proto/web/web-WebService_connectquery'
import { Segment } from '../../gen/proto/web/web_pb'
import { useEffect, useState } from 'react'
import { Style } from '../utils/style'

export function Editor() {
    const navigation = useNavigation()
    const searchParam = navigation.useSearchParams()

    const subtitleIDSearchParam = searchParam.get('subtitle_id')

    if (subtitleIDSearchParam == null) {
        navigation.navigate(Routes.Files)
        return
    }

    const getSubtitleQuery = useQuery(getSubtitle, {
        id: parseInt(subtitleIDSearchParam),
    })

    const getSubtitleSegmentsQuery = useQuery(getSubtitleSegments, {
        id: parseInt(subtitleIDSearchParam),
    })

    const [scale, _setScale] = useState(2)

    const [data, setData] = useState<
        { segment: Segment; isSelected: boolean }[]
    >([])

    useEffect(() => {
        if (!getSubtitleSegmentsQuery.isSuccess) return
        if (data.length > 0) return

        setData(
            getSubtitleSegmentsQuery.data.segments.map((segment) => {
                return { segment, isSelected: false }
            })
        )
    }, [getSubtitleSegmentsQuery.isSuccess])

    let duration = 0
    for (const { segment } of data) {
        if (segment.end !== undefined && segment.end.seconds > duration) {
            duration = Number(segment.end.seconds)
        }
    }
    return (
        <div className="grid h-full w-full grid-cols-[1fr_1fr] grid-rows-[auto_2fr_1fr] gap-md rounded-lg bg-neutral-3 p-lg">
            {!getSubtitleQuery.isSuccess ? (
                <section className="col-span-2 flex h-lg flex-row justify-between">
                    <button
                        className="text-sm text-text-1"
                        onClick={() => {
                            navigation.navigate(Routes.Files)
                        }}
                    >
                        Return
                    </button>
                    <div className="h-lg w-[12rem] animate-pulse rounded-sm bg-neutral-1"></div>
                </section>
            ) : (
                <section className="col-span-2 flex h-lg flex-row justify-between">
                    <button
                        className="text-sm text-text-1"
                        onClick={() => {
                            navigation.navigate(
                                Routes.Files,
                                new URLSearchParams({
                                    video_id:
                                        getSubtitleQuery.data.videoId.toString(),
                                })
                            )
                        }}
                    >
                        Return
                    </button>
                    <p className="text-xs text-text-1">
                        {getSubtitleQuery.data.name}
                    </p>
                </section>
            )}
            <section className="flex flex-col gap-sm rounded-sm p-md outline outline-1 outline-neutral-1">
                {data.map(({ segment, isSelected }) => {
                    return (
                        <div
                            className="rounded-sm p-md"
                            style={{
                                backgroundColor: isSelected
                                    ? Style.colors.secondary[2]
                                    : Style.colors.neutral[1],
                            }}
                        >
                            <p
                                style={{
                                    color: isSelected
                                        ? Style.colors.text[2]
                                        : Style.colors.text[1],
                                }}
                            >
                                {segment.text}
                            </p>
                        </div>
                    )
                })}
            </section>
            <section className="rounded-sm outline outline-1 outline-neutral-1"></section>
            <section className="col-span-2 overflow-scroll rounded-sm p-md outline outline-1 outline-neutral-1">
                <section className="relative h-full">
                    {Array.from(
                        { length: Math.round(duration / 10) + 1 },
                        (_, i) => {
                            const time = i * 10
                            const seconds = time % 60
                            const minutes = Math.floor(time / 60) % 60
                            const hours = Math.floor(i / 3600)

                            return (
                                <div
                                    className="absolute h-full"
                                    style={{ left: time * scale + 'rem' }}
                                >
                                    <p className="h-lg text-text-1">{`${hours}:${minutes}:${seconds}`}</p>
                                    <div className="h-full w-[1px] bg-neutral-1" />
                                </div>
                            )
                        }
                    )}
                    {data.map(({ segment, isSelected }, i) => {
                        const start = Number(segment.start?.seconds || 0)
                        const end = Number(segment.end?.seconds || 0)

                        return (
                            <button
                                className="absolute top-xl rounded-sm p-md"
                                onClick={() => {
                                    data[i].isSelected = !isSelected
                                    setData([...data])
                                }}
                                style={{
                                    left: start * scale + 'rem',
                                    width: (end - start) * scale + 'rem',
                                    backgroundColor: isSelected
                                        ? Style.colors.secondary[2]
                                        : Style.colors.neutral[1],
                                }}
                            >
                                <p
                                    className="break-words"
                                    style={{
                                        color: isSelected
                                            ? Style.colors.text[2]
                                            : Style.colors.text[1],
                                    }}
                                >
                                    {segment.text}
                                </p>
                            </button>
                        )
                    })}
                </section>
            </section>
        </div>
    )
}
