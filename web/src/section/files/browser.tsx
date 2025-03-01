import filepath from 'path-browserify'
import { Routes, useNavigation } from '../../utils/navigation'
import { useQuery } from '@connectrpc/connect-query'
import {
    getRootDirectoryPaths,
    getSubtitle,
    readDirectory,
    searchVideo,
} from '../../../gen/proto/services/web-WebService_connectquery'
import {
    SubtitleOriginalFormat,
    SubtitleStage,
} from '../../../gen/proto/messages/subtitle_pb'

export function Browser() {
    const navigation = useNavigation()
    const searchParams = navigation.useSearchParams()
    const pathSearchParam = searchParams.get('path')

    if (filepath.extname(pathSearchParam || '') == '') {
        return <DirectoryView />
    }

    return <FileView />
}

function DirectoryView() {
    const navigation = useNavigation()
    const searchParams = navigation.useSearchParams()
    const pathSearchParam = searchParams.get('path')

    const isRootDirectory = pathSearchParam == null

    const useGetDirectoryQuery = ({ path }: { path: string | null }) => {
        if (path == null) {
            const query = useQuery(getRootDirectoryPaths)

            if (query.isSuccess) {
                return {
                    ...query,
                    data: {
                        childrenDirectoryPaths: query.data.paths,
                        videoPaths: [],
                    },
                }
            }

            return {
                ...query,
                data: undefined,
            }
        }

        const query = useQuery(readDirectory, { path })

        if (query.isSuccess) {
            return {
                ...query,
                data: {
                    childrenDirectoryPaths: query.data.childrenDirectoryPaths,
                    videoPaths: query.data.videoPaths,
                },
            }
        }

        return {
            ...query,
            data: undefined,
        }
    }

    const getDirectoryQuery = useGetDirectoryQuery({ path: pathSearchParam })

    return (
        <div className="relative h-full overflow-hidden rounded-md bg-neutral-2 p-xl">
            <section className="relative flex h-full flex-col gap-lg overflow-y-scroll">
                {!getDirectoryQuery.isSuccess ? (
                    <section className="flex flex-col gap-md">
                        <div className="h-xl w-[12rem] animate-pulse rounded-sm bg-neutral-1" />
                        <section className="flex flex-col gap-sm">
                            <div className="h-2xl w-full animate-pulse rounded-sm bg-neutral-1" />
                            <div className="h-2xl w-full animate-pulse rounded-sm bg-neutral-1" />
                        </section>
                    </section>
                ) : (
                    <>
                        {getDirectoryQuery.data.childrenDirectoryPaths
                            .length ? (
                            <section className="flex flex-col gap-md">
                                <h2 className="text-lg text-text-1">Folders</h2>
                                <section className="flex flex-col gap-sm">
                                    {getDirectoryQuery.data.childrenDirectoryPaths.map(
                                        (childDirectoryPath) => {
                                            return (
                                                <Directory
                                                    key={childDirectoryPath}
                                                    path={childDirectoryPath}
                                                    displayEntirePath={
                                                        isRootDirectory
                                                    }
                                                />
                                            )
                                        }
                                    )}
                                </section>
                            </section>
                        ) : null}
                        {getDirectoryQuery.data.videoPaths.length ? (
                            <section className="flex flex-col gap-md">
                                <h2 className="text-lg text-text-1">Files</h2>
                                <section className="flex flex-col gap-sm">
                                    {getDirectoryQuery.data.videoPaths.map(
                                        (videoPath) => {
                                            return (
                                                <File
                                                    key={videoPath}
                                                    path={videoPath}
                                                />
                                            )
                                        }
                                    )}
                                </section>
                            </section>
                        ) : null}
                    </>
                )}
            </section>
            <div className="absolute bottom-0 right-0 rounded-tl-md bg-neutral-2 p-md">
                <p className="text-xs text-text-1">{pathSearchParam}</p>
            </div>
        </div>
    )
}

function FileView() {
    const navigation = useNavigation()
    const searchParams = navigation.useSearchParams()
    const pathSearchParam = searchParams.get('path')

    const getVideoQuery = useQuery(searchVideo, { path: pathSearchParam || '' })

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
                    getVideoQuery.data.subtitleIds.map((subtitleId) => {
                        return <Subtitle key={subtitleId} id={subtitleId} />
                    })
                )}
            </section>
            <div className="absolute bottom-0 right-0 rounded-tl-md bg-neutral-2 p-md">
                {!getVideoQuery.isSuccess ? (
                    <div className="h-lg w-[24rem] animate-pulse rounded-sm bg-neutral-1" />
                ) : (
                    <p className="text-xs text-text-1"> {pathSearchParam}</p>
                )}
            </div>
        </div>
    )
}

function Subtitle({ id }: { id: string }) {
    const navigation = useNavigation()

    const getSubtitleQuery = useQuery(getSubtitle, { id })

    if (!getSubtitleQuery.isSuccess) {
        return <div className="h-2xl animate-pulse rounded-sm bg-neutral-1" />
    }

    let status = 'unknown'
    switch (getSubtitleQuery.data.stage) {
        case SubtitleStage.DETECTED:
            status = 'detected'
            break
        case SubtitleStage.EXTRACTED:
            status = 'extracted'
            break
        case SubtitleStage.FORMATED:
            status = 'formated'
            break
        case SubtitleStage.EXPORTED:
            status = 'exported'
            break
    }
    if (getSubtitleQuery.data.isProcessing) {
        switch (getSubtitleQuery.data.stage) {
            case SubtitleStage.DETECTED:
                status = 'extracting'
                break
            case SubtitleStage.EXTRACTED:
                status = 'formating'
                break
            case SubtitleStage.FORMATED:
                status = 'exporting'
                break
        }
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
            <section className="grid grid-flow-row auto-rows-fr gap-sm lg:auto-cols-fr lg:grid-flow-col">
                <p className="text-sm text-text-1">
                    {getSubtitleQuery.data.title}
                </p>
                <section className="grid auto-cols-fr grid-flow-col">
                    <p className="text-center text-xs text-text-1">
                        {getSubtitleQuery.data.language}
                    </p>
                    <p className="text-center text-xs text-text-1">{format}</p>
                    <p className="text-center text-xs text-text-1">{status}</p>
                </section>
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
                    new URLSearchParams({ path })
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

function File({ path }: { path: string }) {
    const navigation = useNavigation()

    return (
        <button
            className="group grid min-h-2xl grid-cols-[0.25rem_1fr] items-center gap-md rounded-sm bg-neutral-1 p-md"
            onClick={() => {
                navigation?.navigate(
                    Routes.Files,
                    new URLSearchParams({ path: path })
                )
            }}
        >
            <div className="h-full rounded-sm bg-primary-2 group-hover:bg-primary-1" />
            <p className="text-start text-sm text-text-1">
                {filepath.basename(path)}
            </p>
        </button>
    )
}
