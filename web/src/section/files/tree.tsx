import filepath from 'path-browserify'
import { useQuery } from '@connectrpc/connect-query'
import { Routes, useNavigation } from '../../utils/navigation'
import { Style } from '../../utils/style'
import {
    getRootDirectoryPaths,
    readDirectory,
} from '../../../gen/proto/services/web-WebService_connectquery'

export function Tree() {
    const getRootDirectoryPathsQuery = useQuery(getRootDirectoryPaths)

    return (
        <section className="h-full w-[16vw] min-w-[16rem] overflow-y-scroll pr-md">
            {!getRootDirectoryPathsQuery.isSuccess ? (
                <section className="flex flex-col gap-sm">
                    <div className="h-[8rem] animate-pulse rounded-sm bg-neutral-2" />
                    <div className="h-[6rem] animate-pulse rounded-sm bg-neutral-2" />
                </section>
            ) : (
                getRootDirectoryPathsQuery.data.paths.map((path) => {
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

function Directory({
    path,
    displayEntirePath,
}: {
    path: string
    displayEntirePath?: boolean
}) {
    const navigation = useNavigation()
    const searchParams = navigation.useSearchParams()
    const pathSearchParam = searchParams.get('path')

    const readDirectoryQuery = useQuery(readDirectory, { path })

    let isOpen = !!pathSearchParam && pathSearchParam.startsWith(path)
    let isActive = pathSearchParam == filepath.dirname(path)

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
                        {!readDirectoryQuery.isSuccess ? (
                            <section className="flex flex-col gap-sm">
                                <div className="h-[4rem] animate-pulse rounded-sm bg-neutral-2" />
                                <div className="h-[2rem] w-3/4 animate-pulse rounded-sm bg-neutral-2" />
                            </section>
                        ) : (
                            <>
                                {readDirectoryQuery.data.childrenDirectoryPaths.map(
                                    (childrenDirectoryPath) => {
                                        return (
                                            <Directory
                                                key={childrenDirectoryPath}
                                                path={childrenDirectoryPath}
                                            />
                                        )
                                    }
                                )}
                                {readDirectoryQuery.data.videoPaths.map(
                                    (videoPath) => {
                                        return (
                                            <Video
                                                key={videoPath}
                                                path={videoPath}
                                            />
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

function Video({ path }: { path: string }) {
    const navigation = useNavigation()
    const searchParams = navigation.useSearchParams()

    const pathSearchParam = searchParams.get('path')
    const isSelected = !!pathSearchParam && pathSearchParam == path

    return (
        <section className="mt-sm flex min-h-sm flex-row gap-sm">
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
                        new URLSearchParams({ path })
                    )
                }}
            >
                {filepath.basename(path)}
            </button>
        </section>
    )
}
