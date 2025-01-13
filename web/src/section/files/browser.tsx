import filepath from 'path-browserify'
import { Routes, useNavigation } from '../../utils/navigation'
import { useQuery } from '@connectrpc/connect-query'
import {
    getDirectory,
    getMediaDirectories,
    getSubtitle,
    getVideo,
} from '../../../gen/proto/web/web-WebService_connectquery'

export function Browser() {
    const navigation = useNavigation()
    const searchParams = navigation.useSearchParams()

    const directoryPathSearchParam = searchParams.get('directory_path')
    const videoIDSearchParam = searchParams.get('video_id')

    if (videoIDSearchParam == null) {
        return <DirectoryView path={directoryPathSearchParam || ''} />
    }

    return <FileView id={parseInt(videoIDSearchParam)} />
}

function DirectoryView({ path }: { path: string }) {
    let isSuccess = false

    let directoryNames: string[] = []
    let videoIDs: number[] = []

    const isRootDirectory = path == ''

    if (isRootDirectory) {
        const getMediaDirectoriesQuery = useQuery(getMediaDirectories)

        if (getMediaDirectoriesQuery.isSuccess) {
            isSuccess = true

            directoryNames = getMediaDirectoriesQuery.data.paths
        }
    } else {
        const getDirectoryQuery = useQuery(getDirectory, { path })

        if (getDirectoryQuery.isSuccess) {
            isSuccess = true

            directoryNames = getDirectoryQuery.data.childrenDirectoryNames
            videoIDs = getDirectoryQuery.data.videoIds
        }
    }

    return (
        <div className="relative h-full overflow-hidden rounded-md bg-neutral-2 p-xl">
            <section className="relative flex h-full flex-col gap-lg overflow-y-scroll">
                {!isSuccess ? (
                    <section className="flex flex-col gap-md">
                        <div className="h-xl w-[12rem] animate-pulse rounded-sm bg-neutral-1" />
                        <section className="flex flex-col gap-sm">
                            <div className="h-2xl w-full animate-pulse rounded-sm bg-neutral-1" />
                            <div className="h-2xl w-full animate-pulse rounded-sm bg-neutral-1" />
                        </section>
                    </section>
                ) : (
                    <>
                        {directoryNames.length ? (
                            <section className="flex flex-col gap-md">
                                <h2 className="text-lg text-text-1">Folders</h2>
                                <section className="flex flex-col gap-sm">
                                    {directoryNames.map((directoryName) => {
                                        return (
                                            <Directory
                                                key={directoryName}
                                                path={filepath.join(
                                                    path,
                                                    directoryName
                                                )}
                                                displayEntirePath={
                                                    isRootDirectory
                                                }
                                            />
                                        )
                                    })}
                                </section>
                            </section>
                        ) : null}
                        {videoIDs.length ? (
                            <section className="flex flex-col gap-md">
                                <h2 className="text-lg text-text-1">Files</h2>
                                <section className="flex flex-col gap-sm">
                                    {videoIDs.map((videoIDs) => {
                                        return (
                                            <File
                                                key={videoIDs}
                                                id={videoIDs}
                                            />
                                        )
                                    })}
                                </section>
                            </section>
                        ) : null}
                    </>
                )}
            </section>
            <div className="absolute bottom-0 right-0 rounded-tl-md bg-neutral-2 p-md">
                <p className="text-xs text-text-1">{path}</p>
            </div>
        </div>
    )
}

function FileView({ id }: { id: number }) {
    const getVideoQuery = useQuery(getVideo, { id })

    return (
        <div className="relative flex flex-col gap-lg rounded-md bg-neutral-2 p-xl">
            <section className="flex flex-col gap-md">
                <h2 className="text-lg text-text-1">Subtitles</h2>
            </section>
            <section className="flex flex-col gap-md">
                {!getVideoQuery.isSuccess ? (
                    <>
                        <div className="h-2xl animate-pulse rounded-sm bg-neutral-1" />
                        <div className="h-2xl animate-pulse rounded-sm bg-neutral-1" />
                    </>
                ) : (
                    getVideoQuery.data.subtitleIds.map((subtitleID) => {
                        return <Subtitle key={subtitleID} id={subtitleID} />
                    })
                )}
            </section>
            <div className="absolute bottom-0 right-0 rounded-tl-md bg-neutral-2 p-md">
                {!getVideoQuery.isSuccess ? (
                    <div className="h-lg animate-pulse bg-neutral-1 rounded-sm w-[24rem]"/>
                ) : (
                    <p className="text-xs text-text-1">
                        {filepath.basename(getVideoQuery.data.filepath)}
                    </p>
                )}
            </div>
        </div>
    )
}

function Subtitle({ id }: { id: number }) {
    const navigation = useNavigation()

    const getSubtitleQuery = useQuery(getSubtitle, {
        id,
    })

    if (!getSubtitleQuery.isSuccess) {
        return <div className="h-2xl animate-pulse rounded-sm bg-neutral-1" />
    }

    return (
        <button
            className="group grid min-h-2xl grid-cols-[0.25rem_1fr] items-center gap-md rounded-sm bg-neutral-1 p-md"
            onClick={() => {
                navigation?.navigate(
                    Routes.Editor,
                    new URLSearchParams({
                        subtitle_id: id.toString(),
                    })
                )
            }}
        >
            <div className="h-full rounded-sm bg-primary-2 group-hover:bg-primary-1" />
            <section className="grid grid-flow-row gap-sm lg:grid-flow-col">
                <p className="text-sm text-text-1">
                    {getSubtitleQuery.data.name}
                </p>
                <div className="text-xs text-text-1">Info</div>
            </section>
        </button>
    )
}

function Directory({
    path,
    displayEntirePath,
}: {
    path: string
    displayEntirePath?: boolean
}) {
    const navigation = useNavigation()

    return (
        <button
            className="group grid min-h-2xl grid-cols-[0.25rem_1fr] items-center gap-md rounded-sm bg-neutral-1 p-md"
            onClick={() => {
                navigation?.navigate(
                    Routes.Files,
                    new URLSearchParams({ directory_path: path })
                )
            }}
        >
            <div className="h-full rounded-sm bg-secondary-2 group-hover:bg-secondary-1" />
            <p className="text-start text-sm text-text-1">
                {displayEntirePath ? path : filepath.basename(path)}
            </p>
        </button>
    )
}

function File({ id }: { id: number }) {
    const navigation = useNavigation()
    const getVideoQuery = useQuery(getVideo, { id })

    if (!getVideoQuery.isSuccess) {
        return <div className="h-2xl animate-pulse rounded-sm bg-neutral-1" />
    }

    return (
        <button
            className="group grid min-h-2xl grid-cols-[0.25rem_1fr] items-center gap-md rounded-sm bg-neutral-1 p-md"
            onClick={() => {
                navigation?.navigate(
                    Routes.Files,
                    new URLSearchParams({ video_id: id.toString() })
                )
            }}
        >
            <div className="h-full rounded-sm bg-primary-2 group-hover:bg-primary-1" />
            <p className="text-start text-sm text-text-1">
                {filepath.basename(getVideoQuery.data.filepath)}
            </p>
        </button>
    )
}
