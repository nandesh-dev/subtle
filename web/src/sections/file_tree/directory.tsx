import { useQuery } from '@connectrpc/connect-query'
import filepath from 'path-browserify'

import { readDirectory } from '@/gen/proto/services/web-WebService_connectquery'

import { LoadingBlock } from '@/src/components/loading_block'

import { Route, useNavigation } from '@/src/utility/navigation'

import { Video } from './video'

type DirectoryOptions = {
    path: string
    displayEntirePath?: boolean
    isRootDirectory?: boolean
}

export function Directory({
    path,
    displayEntirePath,
    isRootDirectory,
}: DirectoryOptions) {
    const navigation = useNavigation()
    const searchParams = navigation.useSearchParams()
    const pathSearchParam = searchParams.get('path')

    const readDirectoryQuery = useQuery(readDirectory, { path })

    const isOpen = !!pathSearchParam && pathSearchParam.startsWith(path)
    const isSelected = pathSearchParam == path

    const onClick = () => {
        if (isSelected) {
            if (isRootDirectory) {
                navigation.navigate(Route.Files, new URLSearchParams())
                return
            }

            navigation.navigate(
                Route.Files,
                new URLSearchParams({ path: filepath.dirname(path) })
            )
            return
        }
        navigation.navigate(Route.Files, new URLSearchParams({ path }))
    }

    if (readDirectoryQuery.isError) {
        return
    }

    if (readDirectoryQuery.isPending) {
        return (
            <section className="flex flex-row gap-2">
                <LoadingBlock className="bg-neutral size-6" />
                <LoadingBlock className="bg-neutral size-6" />
                <LoadingBlock className="bg-neutral size-6" />
            </section>
        )
    }

    return (
        <section className="grid min-h-4 gap-1 grid-cols-[auto_1fr]">
            {!isRootDirectory ? (
                <div
                    className={`mx-2 w-1 rounded-full ${isSelected ? 'bg-secondary' : 'bg-secondary-light'} hover:bg-secondary transition-colors`}
                />
            ) : null}
            <section>
                <button className="text-sm mb-1" onClick={onClick}>
                    {displayEntirePath ? path : filepath.basename(path)}
                </button>
                {isOpen && (
                    <section className="flex flex-col gap-2">
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
                        {readDirectoryQuery.data.videoPaths.map((videoPath) => {
                            return <Video key={videoPath} path={videoPath} />
                        })}
                    </section>
                )}
            </section>
        </section>
    )
}
