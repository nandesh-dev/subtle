import { useNavigate, useSearchParams } from 'react-router-dom'
import { FolderIcon, ProcessingIcon, TickIcon } from '../../../assets'
import { useProto } from '../../context/proto'
import { GetVideoRequest } from '../../../gen/proto/media/media_pb'
import { useQuery } from '@tanstack/react-query'
import { useEffect } from 'react'
import { GetSubtitleRequest } from '../../../gen/proto/subtitle/subtitle_pb'
import { Large, Small } from '../../utils/react_responsive'

export function Video() {
    const { MediaServiceClient } = useProto()
    const [searchParams] = useSearchParams()
    const navigate = useNavigate()

    const rawId = searchParams.get('id')
    if (!rawId || rawId == null) {
        useEffect(() => navigate('/media'))
        return <></>
    }
    const id = parseInt(rawId)

    const { data } = useQuery({
        queryKey: ['get-video', id],
        queryFn: () =>
            MediaServiceClient?.getVideo(new GetVideoRequest({ id })),
    })

    const Subtitles = () =>
        data?.subtitleIds.map((id) => <Subtitle key={id} id={id} />)

    return (
        <section className="flex h-full flex-col gap-sm md:px-lg md:py-xxl">
            <div className="flex flex-row items-center gap-lg md:min-h-[4rem]">
                <h2 className="text-md text-gray-830">Video</h2>
            </div>
            <section className="grid grid-cols-[auto_1fr] gap-md rounded-sm bg-gray-80 p-md">
                <FolderIcon className="h-full fill-orange" />
                <div className="">
                    <div className="flex w-full flex-row items-end">
                        <p className="text-sm text-gray-830">
                            {data?.baseName}
                        </p>
                        <p className="text-xs text-gray-520">
                            {data?.extension}
                        </p>
                    </div>
                    <p className="text-xs text-gray-520">
                        {data?.directoryPath}
                    </p>
                </div>
            </section>
            {(data?.subtitleIds.length || 0) > 0 && (
                <>
                    <div className="flex w-full flex-row items-center gap-md">
                        <h3 className="text-nowrap text-md text-gray-520">
                            Subtitles
                        </h3>
                        <div className="h-[4px] w-full rounded-sm bg-gray-80" />
                    </div>
                    <Subtitles />
                </>
            )}
        </section>
    )
}

type SubtitleProp = {
    id: number
}

function Subtitle({ id }: SubtitleProp) {
    const navigate = useNavigate()
    const { SubtitleServiceClient } = useProto()
    const { isLoading, data } = useQuery({
        queryKey: ['get-subtitle', id],
        queryFn: () =>
            SubtitleServiceClient?.getSubtitle(new GetSubtitleRequest({ id })),
    })

    if (isLoading) {
        return (
            <>
                <Small>
                    <div className="grid grid-rows-2 gap-xs rounded-sm bg-gray-80 p-sm">
                        <div className="flex flex-row justify-between">
                            <p className="text-sm text-gray-830">...</p>
                            <ProcessingIcon className="opacity-0" />
                        </div>
                        <p className="text-start text-xs text-gray-520">...</p>
                    </div>
                </Small>
                <Large>
                    <div className="grid grid-cols-2 rounded-sm bg-gray-80 p-sm">
                        <p className="text-start text-sm text-gray-830">...</p>
                        <div className="grid grid-cols-3 items-center justify-items-end">
                            <p className="col-span-2 text-sm text-gray-520">
                                ...
                            </p>
                            <ProcessingIcon className="opacity-0" />
                        </div>
                    </div>
                </Large>
            </>
        )
    }

    let statusText = 'Detected'
    if (data?.export) {
        const MAX_SIZE = 24
        const exportPath = data.export.path

        if (exportPath.length > MAX_SIZE) {
            statusText = '...' + exportPath.slice(-MAX_SIZE + 3)
        } else {
            statusText = exportPath
        }
    } else if (data?.isProcessing) {
        if (data.segmentIds.length > 0) statusText = 'Exporting'
        else {
            if (data.import?.isExternal) statusText = 'Importing'
            else statusText = 'Extracting'
        }
    } else if ((data?.segmentIds.length || 0) > 0) {
        if (data?.import?.isExternal) statusText = 'Imported'
        else statusText = 'Extracted'
    }

    const navigateToSubtitle = () => {
        const newSearchParam = new URLSearchParams({ id: id.toString() })
        navigate('/subtitle?' + newSearchParam.toString())
    }

    return (
        <>
            <Small>
                <button
                    onClick={navigateToSubtitle}
                    className="grid grid-rows-2 gap-xs rounded-sm bg-gray-80 p-sm"
                >
                    <div className="flex flex-row justify-between">
                        <p className="text-sm text-gray-830">{data?.title}</p>
                        {data?.isProcessing ? <ProcessingIcon /> : <TickIcon />}
                    </div>
                    <p className="text-start text-xs text-gray-520">
                        {statusText}
                    </p>
                </button>
            </Small>
            <Large>
                <button
                    onClick={navigateToSubtitle}
                    className="grid grid-cols-2 rounded-sm bg-gray-80 p-sm"
                >
                    <p className="text-start text-sm text-gray-830">
                        {data?.title}
                    </p>
                    <div className="grid grid-cols-3 items-center justify-items-end">
                        <p className="col-span-2 text-sm text-gray-520">
                            {statusText}
                        </p>
                        {data?.isProcessing ? <ProcessingIcon /> : <TickIcon />}
                    </div>
                </button>
            </Large>
        </>
    )
}
