import { useNavigate, useSearchParams } from 'react-router-dom'
import { useProto } from '../../../context/proto'
import {
    ExtractRawStreamRequest,
    GetVideoRequest,
} from '../../../../gen/proto/media/media_pb'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { Large, Small } from '../../../utils/react_responsive'
import React, { FormEventHandler, useEffect, useState } from 'react'
import { usePopover } from '../../../context/popover'

export function Video() {
    const navigate = useNavigate()
    const popover = usePopover()
    const { MediaServiceClient } = useProto()
    const [searchParams] = useSearchParams()

    const idString = searchParams.get('id')

    if (!idString) {
        useEffect(() => navigate('/media'))
        return
    }

    const id = parseInt(idString)

    const { data } = useQuery({
        queryKey: ['get-video', id],
        queryFn: () =>
            MediaServiceClient?.getVideo(
                new GetVideoRequest({
                    id,
                })
            ),
    })

    return (
        <>
            <Small>
                <section className="flex h-full flex-col gap-sm">
                    <section className="flex flex-col gap-sm">
                        <div className="flex flex-row items-center gap-lg">
                            <h2 className="text-md text-gray-830">Video</h2>
                        </div>
                    </section>
                    <section className="grid grid-flow-row gap-sm overflow-y-auto pb-xxl"></section>
                </section>
            </Small>
            <Large>
                <section className="flex h-full flex-col gap-sm px-lg py-xxl">
                    <section className="grid grid-cols-[1fr_20rem]">
                        <div className="flex flex-row items-center gap-lg">
                            <h2 className="text-md text-gray-830">Video</h2>
                        </div>
                    </section>
                    <section className="grid w-full grid-flow-row gap-sm overflow-y-auto">
                        {data?.subtitles.map((subtitle) => {
                            const isExported = subtitle.exportPath != ''
                            return (
                                <div
                                    key={subtitle.id}
                                    className="grid grid-cols-2 items-center gap-sm rounded-sm bg-gray-80"
                                >
                                    <div className="p-sm">
                                        <p className="text-sm text-gray-830">
                                            {subtitle.title}
                                        </p>
                                    </div>

                                    <div className="grid grid-cols-3 items-center p-xs">
                                        <p className="text-center text-gray-520">
                                            {subtitle.importIsExternal
                                                ? 'External'
                                                : 'Internal'}
                                        </p>
                                        <p className="text-center text-gray-520">
                                            {subtitle.language}
                                        </p>
                                        <button
                                            className={`rounded-xs px-sm py-xs text-xs font-medium text-gray-830 ${isExported ? 'bg-red' : 'bg-primary'}`}
                                        >
                                            {isExported ? 'Delete' : 'Export'}
                                        </button>
                                    </div>
                                </div>
                            )
                        })}
                        <div className="flex w-full flex-row items-center gap-md">
                            <h3 className="text-nowrap text-md text-gray-830">
                                Available Video Subtitles
                            </h3>
                            <div className="h-[4px] w-full rounded-sm bg-gray-80" />
                        </div>
                        {data?.rawStreams.map((rawStream) => {
                            const exportStreamPopover = () => {
                                popover.set(
                                    <ExtractPopover
                                        videoId={id}
                                        rawStreamIndex={rawStream.index}
                                        defaultTitle={rawStream.title}
                                    />
                                )
                            }
                            return (
                                <div
                                    className="grid grid-cols-6 items-center rounded-sm bg-gray-80 p-xs pl-sm"
                                    key={rawStream.index}
                                >
                                    <p className="col-span-3 text-gray-830">
                                        {rawStream.title}
                                    </p>
                                    <p className="text-center text-gray-520">
                                        {rawStream.format}
                                    </p>
                                    <p className="text-center text-gray-520">
                                        {rawStream.language}
                                    </p>
                                    <button
                                        className={`rounded-xs bg-orange px-sm py-xs text-xs font-medium text-gray-830`}
                                        onClick={exportStreamPopover}
                                    >
                                        Extract
                                    </button>
                                </div>
                            )
                        })}
                    </section>
                </section>
            </Large>
        </>
    )
}

type ExportPopoverProp = {
    rawStreamIndex: number
    videoId: number
    defaultTitle: string
}

function ExtractPopover({
    defaultTitle,
    videoId,
    rawStreamIndex,
}: ExportPopoverProp) {
    const queryClient = useQueryClient()
    const popover = usePopover()
    const { MediaServiceClient } = useProto()
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState('')

    const extractSubtitle: FormEventHandler = (
        e: React.FormEvent<HTMLFormElement>
    ) => {
        e.preventDefault()
        const formData = new FormData(e.currentTarget)
        const title = formData.get('title') as string

        const query = queryClient.fetchQuery({
            queryKey: ['extract-raw-stream', rawStreamIndex, videoId],
            queryFn: () =>
                MediaServiceClient?.extractRawStream(
                    new ExtractRawStreamRequest({
                        title,
                        rawStreamIndex,
                        videoId,
                    })
                ),
        })

        setIsLoading(true)

        query
            .then(() => {
                setIsLoading(false)
                queryClient.invalidateQueries({
                    queryKey: ['get-video', videoId],
                })
                popover.reset()
            })
            .catch((err: Error) => {
                setIsLoading(false)
                setError(err.message)
            })
    }

    return (
        <div className="flex h-full flex-col gap-sm">
            <h4 className="text-md text-gray-830">Extract Subtitle</h4>
            <form
                onSubmit={extractSubtitle}
                className="flex h-full flex-col justify-between"
            >
                <label className="flex flex-col gap-xs">
                    <span className="text-sm text-gray-830">Title</span>
                    <input
                        name="title"
                        defaultValue={defaultTitle}
                        placeholder="Subtitle Title"
                        type="text"
                        className="rounded-xxs px-sm py-xs text-gray-520 outline outline-2 outline-gray-190 placeholder:text-gray-190"
                    />
                </label>

                <div className="flex gap-sm sm:flex-col md:flex-row md:justify-between">
                    <button
                        type="submit"
                        disabled={isLoading}
                        className="w-fit rounded-xs bg-orange px-sm py-xs text-gray-830 disabled:bg-gray-190"
                    >
                        Extract
                    </button>
                    {error && <p className="text-gray-520">{error}</p>}
                </div>
            </form>
        </div>
    )
}
