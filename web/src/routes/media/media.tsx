import { useNavigate, useSearchParams } from 'react-router-dom'
import { FolderIcon, SearchIcon } from '../../../assets'
import { Large, Small } from '../../utils/react_responsive'
import { useProto } from '../../context/proto'
import {
    GetDirectoryRequest,
    GetVideoRequest,
} from '../../../gen/proto/media/media_pb'
import { useQuery } from '@tanstack/react-query'

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
                        <div className="grid grid-cols-[repeat(auto-fill,minmax(18rem,1fr))] gap-sm">
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
                <div className="">
                    <p className="text-start text-sm text-gray-830">
                        Loading...
                    </p>
                    <div className="flex w-full flex-row justify-between">
                        <p className="text-xs text-gray-520">Subtitle</p>
                        <p className="text-xs text-gray-520">../..</p>
                    </div>
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
            <div className="">
                <p className="text-start text-sm text-gray-830">{data?.name}</p>
                <div className="flex w-full flex-row justify-between">
                    <p className="text-xs text-gray-520">Subtitle</p>
                    <p className="text-xs text-gray-520">10/20</p>
                </div>
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

    const { data, isLoading } = useQuery({
        queryKey: ['get-video', id],
        queryFn: () =>
            MediaServiceClient?.getVideo(new GetVideoRequest({ id })),
    })

    if (isLoading) {
        return (
            <>
                <Small>
                    <div className="grid grid-rows-2 gap-sm rounded-sm bg-gray-80 p-sm">
                        <p className="text-start text-sm text-gray-830">...</p>
                        <div className="flex flex-row justify-between">
                            <p className="text-sm text-gray-520">...</p>
                        </div>
                    </div>
                </Small>
                <Large>
                    <div className="grid grid-cols-2 rounded-sm bg-gray-80 p-sm">
                        <p className="text-start text-sm text-gray-830">...</p>
                        <div className="grid grid-cols-3">
                            <p className="text-sm text-gray-520">...</p>
                        </div>
                    </div>
                </Large>
            </>
        )
    }

    const onClick = () => {
        const newSearchParam = new URLSearchParams({ id: id.toString() })
        navigate('/video?' + newSearchParam.toString(), {})
    }

    return (
        <>
            <Small>
                <button
                    className="grid grid-rows-2 gap-sm rounded-sm bg-gray-80 p-sm"
                    onClick={onClick}
                >
                    <p className="text-start text-sm text-gray-830">
                        {data?.baseName}
                    </p>
                    <div className="flex flex-row justify-between">
                        <p className="text-sm text-gray-520">
                            {data?.extension}
                        </p>
                    </div>
                </button>
            </Small>
            <Large>
                <button
                    className="grid grid-cols-2 rounded-sm bg-gray-80 p-sm"
                    onClick={onClick}
                >
                    <p className="text-start text-sm text-gray-830">
                        {data?.baseName}
                    </p>
                    <div className="grid grid-cols-3 justify-items-end">
                        <p className="text-sm text-gray-520">
                            {data?.extension}
                        </p>
                    </div>
                </button>
            </Large>
        </>
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
