import { useNavigate, useSearchParams } from 'react-router-dom'
import { useProto } from '../../../context/proto'
import { GetVideoRequest } from '../../../../gen/proto/media/media_pb'
import { useQuery } from '@tanstack/react-query'
import { Large, Small } from '../../../utils/react_responsive'
import { useEffect } from 'react'

export function Video() {
    const navigate = useNavigate()
    const { MediaServiceClient } = useProto()
    const [searchParams] = useSearchParams()

    const [directoryPath, name, extension] = [
        searchParams.get('directoryPath'),
        searchParams.get('name'),
        searchParams.get('extension'),
    ]

    if (!directoryPath || !name || !extension) {
        useEffect(() => navigate('/media'))
        return
    }

    const { data } = useQuery({
        queryKey: [],
        queryFn: () =>
            MediaServiceClient?.getVideo(
                new GetVideoRequest({
                    name,
                    extension,
                    directoryPath,
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
                                <div className="grid grid-cols-2 items-center gap-sm rounded-sm bg-gray-80">
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
                            return (
                                <>
                                    <div className="grid grid-cols-6 items-center rounded-sm bg-gray-80 p-xs pl-sm">
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
                                        >
                                            Extract
                                        </button>
                                    </div>
                                </>
                            )
                        })}
                    </section>
                </section>
            </Large>
        </>
    )
}
