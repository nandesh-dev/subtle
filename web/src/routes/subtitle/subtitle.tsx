import { useNavigate, useSearchParams } from 'react-router-dom'
import { CrossIcon, ProcessingIcon, SubtitleIcon } from '../../../assets'
import { useProto } from '../../context/proto'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import {
    GetSegmentRequest,
    GetSegmentResponse,
    GetSubtitleRequest,
    UpdateSegmentRequest,
} from '../../../gen/proto/subtitle/subtitle_pb'
import { Large, Small } from '../../utils/react_responsive'
import { useIntersectionObserver } from '../../utils/useIntersectionObserver'
import { GetVideoRequest } from '../../../gen/proto/media/media_pb'

export function Subtitle() {
    const { SubtitleServiceClient, MediaServiceClient } = useProto()
    const [searchParams] = useSearchParams()

    const navigate = useNavigate()

    const rawId = searchParams.get('id')
    if (!rawId || rawId == null) {
        useEffect(() => navigate('/media'))
        return <></>
    }
    const id = parseInt(rawId)

    const { data: subtitleData } = useQuery({
        queryKey: ['get-subtitle', id],
        queryFn: () =>
            SubtitleServiceClient?.getSubtitle(new GetSubtitleRequest({ id })),
    })

    const { data: videoData } = useQuery({
        queryKey: ['get-video', subtitleData?.videoId],
        queryFn: () =>
            MediaServiceClient?.getVideo(
                new GetVideoRequest({ id: subtitleData?.videoId })
            ),
        enabled: !!subtitleData?.videoId,
    })

    let statusText = 'Detected'
    if (subtitleData?.export) {
        const MAX_SIZE = 24
        const exportPath = subtitleData.export.path

        if (exportPath.length > MAX_SIZE) {
            statusText = '...' + exportPath.slice(-MAX_SIZE + 3)
        } else {
            statusText = exportPath
        }
    } else if (subtitleData?.isProcessing) {
        if (subtitleData.segmentIds.length > 0) statusText = 'Exporting'
        else {
            if (subtitleData.import?.isExternal) statusText = 'Importing'
            else statusText = 'Extracting'
        }
    } else if ((subtitleData?.segmentIds.length || 0) > 0) {
        if (subtitleData?.import?.isExternal) statusText = 'Imported'
        else statusText = 'Extracted'
    }

    const Segments = () =>
        subtitleData?.segmentIds.map((id) => <Segment id={id} key={id} />)

    return (
        <section className="flex h-full flex-col gap-sm md:px-lg md:py-xxl">
            <div className="flex flex-row items-center gap-lg md:min-h-[4rem]">
                <h2 className="text-md text-gray-830">Subtitle</h2>
                <p className="text-sm text-gray-520">{`${videoData?.directoryPath ?? ''}/${videoData?.baseName ?? ''}${videoData?.extension ?? ''}`}</p>
            </div>
            <section className="grid grid-cols-[auto_1fr] gap-md rounded-sm bg-gray-80 p-md">
                <SubtitleIcon className="h-full fill-yellow" />
                <div className="">
                    <input
                        className="text-sm text-gray-830"
                        defaultValue={subtitleData?.title}
                    />
                    <p className="text-xs text-gray-520">{statusText}</p>
                </div>
            </section>
            <section className="overflow-y-auto">
                {(subtitleData?.segmentIds.length || 0) > 0 && (
                    <>
                        <div className="mb-sm flex w-full flex-row items-center gap-md">
                            <h3 className="text-nowrap text-md text-gray-520">
                                Segments
                            </h3>
                            <div className="h-[4px] w-full rounded-sm bg-gray-80" />
                        </div>
                        <div className="flex flex-col gap-xs">
                            <Segments />
                        </div>
                    </>
                )}
            </section>
        </section>
    )
}

type SegmentProp = {
    id: number
}

function Segment({ id }: SegmentProp) {
    const { isIntersecting, observationElementRef } = useIntersectionObserver()

    const queryClient = useQueryClient()
    const { SubtitleServiceClient } = useProto()
    const [data, setData] = useState<GetSegmentResponse | undefined>()
    const isLoading = data == undefined

    const loadData = () => {
        queryClient
            .fetchQuery({
                queryKey: ['get-segment', id],
                queryFn: () =>
                    SubtitleServiceClient?.getSegment(
                        new GetSegmentRequest({ id })
                    ),
            })
            .then(setData)
    }

    const updateTextMutation = useMutation({
        mutationFn: async (newText: string) => {
            await SubtitleServiceClient?.updateSegment(
                new UpdateSegmentRequest({
                    id,
                    start: data?.start,
                    end: data?.end,
                    new: {
                        text: newText,
                    },
                })
            )
        },
    })

    if (isLoading) {
        if (isIntersecting) loadData()
        return (
            <div
                className="flex min-h-[16rem] items-center justify-center rounded-sm bg-gray-80 md:min-h-[8rem]"
                ref={observationElementRef}
            >
                <p className="text-sm text-gray-830">...</p>
            </div>
        )
    }

    console.log(data)

    const start = new Date(
        Number(data.start?.seconds || 0) * 1000 +
            Number(data.start?.nanos || 0) / 1000
    )

    const end = new Date(
        Number(data.end?.seconds || 0) * 1000 +
            Number(data.end?.nanos || 0) / 1000
    )

    return (
        <>
            <Small>
                <div className="flex flex-col gap-sm rounded-sm bg-gray-80 p-sm">
                    <div className="flex flex-row justify-between">
                        <p className="text-xs text-gray-520">
                            {start.toLocaleTimeString(undefined, {
                                hour12: false,
                            })}
                        </p>
                        <p className="text-xs text-gray-520">
                            {end.toLocaleTimeString(undefined, {
                                hour12: false,
                            })}
                        </p>
                    </div>
                    <div className="flex h-full flex-col">
                        <div className="relative flex flex-col items-center p-md">
                            {data.original?.text ? (
                                <div className="flex items-center justify-center p-sm pt-md">
                                    <p className="text-center text-gray-830">
                                        {data.original?.text}
                                    </p>
                                </div>
                            ) : (
                                data.original?.image && (
                                    <img
                                        src={URL.createObjectURL(
                                            new Blob([data.original?.image], {
                                                type: 'image/png',
                                            })
                                        )}
                                    />
                                )
                            )}
                            <h5 className="absolute left-0 top-0 text-xs text-gray-520">
                                Original
                            </h5>
                        </div>
                        <div className="relative">
                            <div className="flex flex-col items-center justify-center p-md pb-sm">
                                <textarea
                                    className="w-full text-center text-gray-830"
                                    onChange={(e) =>
                                        updateTextMutation.mutate(
                                            e.target.value
                                        )
                                    }
                                    defaultValue={data.new?.text}
                                />
                            </div>
                            <h5 className="absolute left-0 top-0 text-xs text-gray-520">
                                New
                            </h5>
                        </div>
                    </div>
                </div>
            </Small>
            <Large>
                <div className="flex h-fit min-h-[8rem] flex-row gap-sm rounded-sm bg-gray-80 p-sm">
                    <div className="grid w-full grid-cols-2">
                        <div className="relative">
                            <div className="flex h-full items-center justify-center p-sm pt-md">
                                {data.original?.text ? (
                                    <p className="text-center text-gray-830">
                                        {data.original?.text}
                                    </p>
                                ) : (
                                    data.original?.image && (
                                        <img
                                            src={URL.createObjectURL(
                                                new Blob(
                                                    [data.original?.image],
                                                    { type: 'image/png' }
                                                )
                                            )}
                                        />
                                    )
                                )}
                            </div>
                            <h5 className="absolute left-0 top-0 text-xs text-gray-520">
                                Original
                            </h5>
                        </div>
                        <div className="relative">
                            <div className="flex h-full items-center justify-center p-md pb-sm">
                                <textarea
                                    className="w-full text-center text-gray-830"
                                    onChange={(e) =>
                                        updateTextMutation.mutate(
                                            e.target.value
                                        )
                                    }
                                    defaultValue={data.new?.text}
                                />
                            </div>
                            <h5 className="absolute left-0 top-0 text-xs text-gray-520">
                                New
                            </h5>
                            {updateTextMutation.isPending && (
                                <ProcessingIcon className="absolute right-0 top-0 fill-gray-520" />
                            )}
                            {updateTextMutation.isError && (
                                <CrossIcon className="absolute right-0 top-0 fill-gray-520" />
                            )}
                        </div>
                    </div>
                    <div className="flex flex-col justify-between">
                        <p className="text-xs text-gray-520">
                            {start.toLocaleTimeString(undefined, {
                                hour12: false,
                            })}
                        </p>
                        <p className="text-xs text-gray-520">
                            {end.toLocaleTimeString(undefined, {
                                hour12: false,
                            })}
                        </p>
                    </div>
                </div>
            </Large>
        </>
    )
}
