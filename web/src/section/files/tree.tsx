import filepath from 'path-browserify'
import { useQuery } from '@connectrpc/connect-query'
import {
    getDirectory,
    getMediaDirectories,
    getVideo,
} from '../../../gen/proto/web/web-WebService_connectquery'
import { Routes, useNavigation } from '../../utils/navigation'
import { Style } from '../../utils/style'

export function Tree() {
    const getMediaDirectoriesQuery = useQuery(getMediaDirectories)

    return (
        <section className="h-full w-[16vw] min-w-[16rem] overflow-y-scroll pr-md">
            {!getMediaDirectoriesQuery.isSuccess ? (
                <section className="flex flex-col gap-sm">
                    <div className="h-[8rem] animate-pulse rounded-sm bg-neutral-2" />
                    <div className="h-[6rem] animate-pulse rounded-sm bg-neutral-2" />
                </section>
            ) : (
                getMediaDirectoriesQuery.data.paths.map((path) => {
                    return (
                        <Directory
                            key={path}
                            path={path}
                            displayEntirePath={true}
                        />
                    )
                })
            )}
        </section>
    )
}

function File({ id }: { id: number }) {
    const navigation = useNavigation()
    const searchParams = navigation.useSearchParams()

    const getVideoQuery = useQuery(getVideo, { id })

    const videoIDSearchParam = searchParams.get('video_id')
    const isSelected =
        videoIDSearchParam !== null && parseInt(videoIDSearchParam) == id

    return (
        <section className="mt-sm flex min-h-sm flex-row gap-sm">
            {!getVideoQuery.isSuccess ? (
                <div className="h-[4rem] w-full animate-pulse rounded-sm bg-neutral-2" />
            ) : (
                <>
                    {isSelected && (
                        <section className="flex w-md min-w-md flex-col items-center">
                            <div className="h-full w-[2px] rounded-sm bg-primary-1" />
                        </section>
                    )}
                    <button
                        className="text-start text-xs text-text-1"
                        onClick={() => {
                            navigation.navigate(
                                Routes.Files,
                                new URLSearchParams({ video_id: id.toString() })
                            )
                        }}
                    >
                        {filepath.basename(getVideoQuery.data.filepath)}
                    </button>
                </>
            )}
        </section>
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
    const searchParams = navigation.useSearchParams()

    const getDirectoryQuery = useQuery(getDirectory, { path })

    const directoryPathSearchParam = searchParams.get('directory_path')
    const videoIDSearchParam = searchParams.get('video_id')

    let isOpen = false
    let isActive = false

    const getVideoQuery = useQuery(
        getVideo,
        {
            id: parseInt(videoIDSearchParam || '-1'),
        },
        { enabled: videoIDSearchParam !== null }
    )

    if (directoryPathSearchParam !== null) {
        isOpen = directoryPathSearchParam.startsWith(path)
        isActive = directoryPathSearchParam == path
    } else if (videoIDSearchParam !== null) {
        if (getVideoQuery.isSuccess) {
            isOpen = getVideoQuery.data.filepath.startsWith(path)
            isActive = filepath.dirname(getVideoQuery.data.filepath) == path
        }
    }

    return (
        <section>
            <button
                className="text-xs text-text-1"
                onClick={() => {
                    navigation.navigate(
                        Routes.Files,
                        new URLSearchParams({ directory_path: path })
                    )
                }}
            >
                {displayEntirePath ? path : filepath.basename(path)}
            </button>
            {isOpen && (
                <section className="mt-sm flex min-h-sm flex-row gap-sm">
                    <section className="flex w-md min-w-md flex-col items-center">
                        <div
                            className="h-full w-[2px] rounded-sm"
                            style={{
                                backgroundColor: isActive
                                    ? Style.colors.secondary[1]
                                    : Style.colors.secondary[2],
                            }}
                        />
                    </section>
                    <section className="flex w-full flex-col gap-sm">
                        {!getDirectoryQuery.isSuccess ? (
                            <section className="flex flex-col gap-sm">
                                <div className="h-[4rem] animate-pulse rounded-sm bg-neutral-2" />
                                <div className="h-[2rem] w-3/4 animate-pulse rounded-sm bg-neutral-2" />
                            </section>
                        ) : (
                            <>
                                {getDirectoryQuery.data.childrenDirectoryNames.map(
                                    (childrenDirectoryName) => {
                                        return (
                                            <Directory
                                                key={childrenDirectoryName}
                                                path={filepath.join(
                                                    path,
                                                    childrenDirectoryName
                                                )}
                                            />
                                        )
                                    }
                                )}
                                {getDirectoryQuery.data.videoIds.map(
                                    (videoId) => {
                                        return (
                                            <File key={videoId} id={videoId} />
                                        )
                                    }
                                )}
                            </>
                        )}
                    </section>
                </section>
            )}
        </section>
    )
}
