import { useNavigate, useSearchParams } from 'react-router-dom'
import { FolderIcon, SearchIcon } from '../../../assets'
import { Large, Small } from '../../utils/react_responsive'
import { useProto } from '../../context/proto'
import {
    GetDirectoryRequest,
    GetVideoRequest,
} from '../../../gen/proto/media/media_pb'
import { useQueries, useQuery } from '@tanstack/react-query'
import { GetSubtitleRequest } from '../../../gen/proto/subtitle/subtitle_pb'

export function Media() {
    const { MediaServiceClient } = useProto()
    const [searchParams] = useSearchParams()

    const path = searchParams.get('path') || ''

    const { data } = useQuery({
        queryKey: ['get-directory', path],
        queryFn: () =>
            MediaServiceClient?.getDirectory(new GetDirectoryRequest({ path })),
    })

    const Videos = () =>
        data?.videoIds.map((id) => {
            return <Video id={id} key={id} />
        })

    const Folders = () =>
        data?.childrenPaths.map((path) => {
            return <Folder key={path} path={path} />
        })

    return (
        <>
            <Small>
                <section className="flex h-full flex-col gap-sm">
                    <section className="flex flex-col gap-sm">
                        <div className="flex flex-row items-center gap-lg">
                            <h2 className="text-md text-gray-830">Media</h2>
                            <p className="text-sm text-gray-520">/{path}</p>
                        </div>
                        <SearchBar />
                    </section>
                    <section className="grid grid-flow-row gap-sm overflow-y-auto pb-xxl">
                        <Folders />
                        {(data?.videoIds.length || 0) > 0 && (
                            <>
                                <div className="flex w-full flex-row items-center gap-md">
                                    <h3 className="text-nowrap text-md text-gray-520">
                                        Videos
                                    </h3>
                                    <div className="h-[4px] w-full rounded-sm bg-gray-80" />
                                </div>
                                <Videos />
                            </>
                        )}
                    </section>
                </section>
            </Small>
            <Large>
                <section className="flex h-full flex-col gap-sm px-lg py-xxl">
                    <section className="grid min-h-[4rem] grid-cols-[1fr_20rem] items-center">
                        <div className="flex flex-row gap-lg">
                            <h2 className="text-md text-gray-830">Media</h2>
                            <p className="text-sm text-gray-520">/{path}</p>
                        </div>
                        <SearchBar />
                    </section>
                    <section className="grid w-full grid-flow-row gap-sm overflow-y-auto">
                        <div className="grid grid-cols-[repeat(auto-fill,minmax(24rem,1fr))] gap-sm">
                            <Folders />
                        </div>
                        {(data?.videoIds.length || 0) > 0 && (
                            <>
                                <div className="flex w-full flex-row items-center gap-md">
                                    <h3 className="text-nowrap text-md text-gray-520">
                                        Videos
                                    </h3>
                                    <div className="h-[4px] w-full rounded-sm bg-gray-80" />
                                </div>
                                <Videos />
                            </>
                        )}
                    </section>
                </section>
            </Large>
        </>
    )
}

type FolderProp = {
    path: string
}

function Folder({ path }: FolderProp) {
    const { MediaServiceClient } = useProto()
    const [searchParams, setSearchParams] = useSearchParams()

    const { data, isLoading } = useQuery({
        queryKey: ['get-directory', path],
        queryFn: () =>
            MediaServiceClient?.getDirectory(new GetDirectoryRequest({ path })),
    })

    if (isLoading) {
        return (
            <div className="grid grid-cols-[auto_1fr] gap-md rounded-sm bg-gray-80 p-md">
                <FolderIcon className="h-full fill-red" />
                <div className="grid grid-flow-row gap-xs">
                    <div className="h-md animate-pulse rounded-sm bg-gray-190 text-start text-sm text-gray-830" />
                    <div className="h-sm w-1/2 animate-pulse rounded-sm bg-gray-120" />
                </div>
            </div>
        )
    }

    return (
        <button
            className="grid grid-cols-[auto_1fr] gap-md rounded-sm bg-gray-80 p-md"
            onClick={() => {
                const newSearchParam = new URLSearchParams(searchParams)
                newSearchParam.set('path', path)
                setSearchParams(newSearchParam)
            }}
        >
            <FolderIcon className="h-full fill-red" />
            <div className="grid grid-flow-row gap-xs">
                <p className="text-start text-sm text-gray-830">{data?.name}</p>
                <p className="text-start text-xs text-gray-520">
                    {data?.videoIds.length} Videos /{' '}
                    {data?.childrenPaths.length} Folders
                </p>
            </div>
        </button>
    )
}

type FileProp = {
    id: number
}

function Video({ id }: FileProp) {
    const { MediaServiceClient } = useProto()
    const navigate = useNavigate()

    const { data: video, isLoading } = useQuery({
        queryKey: ['get-video', id],
        queryFn: () =>
            MediaServiceClient?.getVideo(new GetVideoRequest({ id })),
    })

    if (isLoading) {
        return (
            <div className="flex flex-col gap-xs rounded-sm bg-gray-80 p-sm lg:grid lg:grid-cols-2 lg:items-center">
                <div className="h-md animate-pulse rounded-sm bg-gray-190" />
                <div className="flex grid-cols-3 flex-row justify-items-end gap-xs lg:grid">
                    <div className="h-sm w-[4rem] animate-pulse rounded-sm bg-gray-120" />
                    <div className="h-sm w-[8rem] animate-pulse rounded-sm bg-gray-120" />
                    <div className="h-sm w-[2rem] animate-pulse rounded-sm bg-gray-120" />
                </div>
            </div>
        )
    }

    const onClick = () => {
        const newSearchParam = new URLSearchParams({ id: id.toString() })
        navigate('/video?' + newSearchParam.toString(), {})
    }

    return (
        <button
            className="flex grid-cols-2 flex-col gap-xs rounded-sm bg-gray-80 p-sm lg:grid lg:items-center"
            onClick={onClick}
        >
            <p className="text-start text-sm text-gray-830">
                {video?.baseName}
            </p>
            <div className="flex grid-cols-3 flex-row justify-items-end gap-xs lg:grid">
                <p className="text-sm text-gray-520">{video?.extension}</p>
            </div>
        </button>
    )
}

function SearchBar() {
    return (
        <div className="flex h-fit w-full flex-row items-center gap-sm rounded-md bg-gray-120 px-sm py-xs">
            <SearchIcon className="fill-gray-520" />
            <input
                className="w-full text-sm text-gray-830 placeholder:text-sm placeholder:text-gray-190"
                placeholder="Search"
            />
        </div>
    )
}
