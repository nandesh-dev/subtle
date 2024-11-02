import { useNavigate, useSearchParams } from 'react-router-dom'
import { FolderIcon, SearchIcon } from '../../../assets'
import { Large, Small } from '../../utils/react_responsive'
import { useProto } from '../../context/proto'
import {
    GetDirectoryRequest,
    GetVideoRequest,
} from '../../../gen/proto/media/media_pb'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'

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
        data?.videos.map((video) => {
            return (
                <File
                    id={video.id}
                    key={video.baseName}
                    baseName={video.baseName}
                    extension={video.extension}
                />
            )
        })

    return (
        <>
            <Small>
                <section className="flex h-full flex-col gap-sm">
                    <section className="flex flex-col gap-sm">
                        <SearchBar />
                        <div className="flex flex-row items-center gap-lg">
                            <h2 className="text-md text-gray-830">Media</h2>
                            <p className="text-sm text-gray-520">/{path}</p>
                        </div>
                    </section>
                    <section className="grid grid-flow-row gap-sm overflow-y-auto pb-xxl">
                        {data?.directories.map((directory) => {
                            return (
                                <Folder
                                    key={directory.path}
                                    name={directory.name}
                                    path={directory.path}
                                    subtitle={{ present: 10, total: 20 }}
                                />
                            )
                        })}
                        <div
                            className="h-[4px] rounded-sm bg-gray-80"
                            content="2"
                        />
                        <Videos />
                    </section>
                </section>
            </Small>
            <Large>
                <section className="flex h-full flex-col gap-sm px-lg py-xxl">
                    <section className="grid grid-cols-[1fr_20rem]">
                        <div className="flex flex-row items-center gap-lg">
                            <h2 className="text-md text-gray-830">Media</h2>
                            <p className="text-sm text-gray-520">/{path}</p>
                        </div>
                        <SearchBar />
                    </section>
                    <section className="grid w-full grid-flow-row gap-sm overflow-y-auto">
                        <div className="grid grid-cols-[repeat(auto-fill,minmax(18rem,1fr))] gap-sm">
                            {data?.directories.map((directory) => {
                                return (
                                    <Folder
                                        key={directory.path}
                                        name={directory.name}
                                        path={directory.path}
                                        subtitle={{
                                            present: 10,
                                            total: 20,
                                        }}
                                    />
                                )
                            })}
                        </div>
                        <div className="h-[4px] rounded-sm bg-gray-80" />
                        <Videos />
                    </section>
                </section>
            </Large>
        </>
    )
}

type FolderProp = {
    name: string
    path: string
    subtitle: {
        present: number
        total: number
    }
}

function Folder({ name, subtitle, path }: FolderProp) {
    const queryClient = useQueryClient()
    const { MediaServiceClient } = useProto()
    const [searchParams, setSearchParams] = useSearchParams()
    const [isPrefetched, setIsPrefetched] = useState(false)

    return (
        <button
            className="grid grid-cols-[auto_1fr] gap-md rounded-sm bg-gray-80 p-md"
            onMouseEnter={
                isPrefetched
                    ? undefined
                    : () => {
                          queryClient.prefetchQuery({
                              queryKey: ['get-directory', path],
                              queryFn: () =>
                                  MediaServiceClient?.getDirectory(
                                      new GetDirectoryRequest({ path })
                                  ),
                          })
                          setIsPrefetched(true)
                      }
            }
            onClick={() => {
                const newSearchParam = new URLSearchParams(searchParams)
                newSearchParam.set('path', path)
                setSearchParams(newSearchParam)
            }}
        >
            <FolderIcon className="h-full fill-red" />
            <div className="">
                <p className="text-start text-sm text-gray-830">{name}</p>
                <div className="flex w-full flex-row justify-between">
                    <p className="text-xs text-gray-520">Subtitle</p>
                    <p className="text-xs text-gray-520">
                        {subtitle.present}/{subtitle.total}
                    </p>
                </div>
            </div>
        </button>
    )
}

type FileProp = {
    id: number
    baseName: string
    extension: string
}

function File({ id, baseName, extension }: FileProp) {
    const queryClient = useQueryClient()
    const { MediaServiceClient } = useProto()
    const [isPrefetched, setIsPrefetched] = useState(false)
    const navigate = useNavigate()

    const onClick = () => {
        const newSearchParam = new URLSearchParams({ id: id.toString() })

        navigate('/media/video?' + newSearchParam.toString(), {})
    }

    const onMouseEnter = isPrefetched
        ? undefined
        : () => {
              queryClient.prefetchQuery({
                  queryKey: ['get-video', id],
                  queryFn: () =>
                      MediaServiceClient?.getVideo(new GetVideoRequest({ id })),
              })

              setIsPrefetched(true)
          }

    return (
        <>
            <Small>
                <button
                    className="grid grid-rows-2 gap-sm rounded-sm bg-gray-80 p-sm"
                    onClick={onClick}
                    onMouseEnter={onMouseEnter}
                >
                    <p className="text-start text-sm text-gray-830">
                        {baseName}
                    </p>
                    <div className="flex flex-row justify-between">
                        <p className="text-sm text-gray-520">{extension}</p>
                    </div>
                </button>
            </Small>
            <Large>
                <button
                    className="grid grid-cols-2 rounded-sm bg-gray-80 p-sm"
                    onClick={onClick}
                    onMouseEnter={onMouseEnter}
                >
                    <p className="text-start text-sm text-gray-830">
                        {baseName}
                    </p>
                    <div className="grid grid-cols-3">
                        <p className="text-sm text-gray-520">{extension}</p>
                    </div>
                </button>
            </Large>
        </>
    )
}

function SearchBar() {
    return (
        <div className="flex h-full w-full flex-row items-center gap-sm rounded-md bg-gray-120 px-sm py-xs">
            <SearchIcon className="fill-gray-520" />
            <input
                className="w-full text-sm text-gray-830 placeholder:text-sm placeholder:text-gray-190"
                placeholder="Search"
            />
        </div>
    )
}
