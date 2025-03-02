import { useQuery } from '@connectrpc/connect-query'

import {
    getRootDirectoryPaths,
    readDirectory,
} from '@/gen/proto/services/web-WebService_connectquery'

import { LoadingBlock } from '@/src/components/loading_block'

import { useSearchParams } from '@/src/utility/navigation'

import { Directory } from './directory'
import { SearchBar } from './search_bar'
import { Video } from './video'

type UseGetDirectoryQueryOption = {
    path: string
    isRootDirectory: boolean
}

const useGetDirectoryQuery = ({
    path,
    isRootDirectory,
}: UseGetDirectoryQueryOption) => {
    if (isRootDirectory) {
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

export function DirectoryView() {
    const pathSearchParam = useSearchParams().get('path')

    const isRootDirectory = pathSearchParam == null

    const getDirectoryQuery = useGetDirectoryQuery({
        path: pathSearchParam || '',
        isRootDirectory,
    })

    if (getDirectoryQuery.isError) {
        return (
            <section className="w-full grid grid-rows-[auto_3fr_7fr] gap-4">
                <SearchBar />
                <section className="flex flex-col rounded-xl bg-neutral p-4">
                    <h2 className="text-lg mb-2">Folders</h2>
                    <section className=" flex items-center justify-center">
                        <p>{getDirectoryQuery.error.message}</p>
                    </section>
                </section>
                <section className="flex flex-col rounded-xl bg-neutral p-4">
                    <h2 className="text-lg mb-2">Videos</h2>
                    <section className=" flex items-center justify-center">
                        <p>{getDirectoryQuery.error.message}</p>
                    </section>
                </section>
            </section>
        )
    }

    if (getDirectoryQuery.isPending) {
        return (
            <section className="w-full grid grid-rows-[auto_3fr_7fr] gap-4">
                <SearchBar />
                <section className="flex flex-col rounded-xl bg-neutral p-4">
                    <h2 className="text-lg">Folders</h2>
                    <section className="h-full flex items-center justify-center gap-8">
                        <LoadingBlock className="bg-neutral-light size-16" />
                        <LoadingBlock className="bg-neutral-light size-16" />
                        <LoadingBlock className="bg-neutral-light size-16" />
                    </section>
                </section>
                <section className="flex flex-col rounded-xl bg-neutral p-4">
                    <h2 className="text-lg">Videos</h2>
                    <section className="h-full flex items-center justify-center gap-8">
                        <LoadingBlock className="bg-neutral-light size-16" />
                        <LoadingBlock className="bg-neutral-light size-16" />
                        <LoadingBlock className="bg-neutral-light size-16" />
                    </section>
                </section>
            </section>
        )
    }

    return (
        <section className="grid grid-rows-[auto_3fr_7fr] gap-4 grid-cols-1 overflow-hidden">
            <SearchBar />
            <section className="flex min-h-32 flex-col rounded-xl bg-neutral p-4 overflow-hidden">
                <h2 className="text-lg">Folders</h2>
                {getDirectoryQuery.data.childrenDirectoryPaths.length > 0 ? (
                    <section className="h-full overflow-y-auto mt-2">
                        <section className="flex flex-wrap gap-2 w-full h-fit">
                            {getDirectoryQuery.data.childrenDirectoryPaths.map(
                                (childDirectoryPath) => {
                                    return (
                                        <Directory
                                            key={childDirectoryPath}
                                            path={childDirectoryPath}
                                            displayEntirePath={isRootDirectory}
                                        />
                                    )
                                }
                            )}
                        </section>
                    </section>
                ) : (
                    <section className="h-full flex items-center justify-center">
                        <p className="text-sm">No folders found!</p>
                    </section>
                )}
            </section>
            <section className="flex min-h-32 flex-col rounded-xl bg-neutral p-4 overflow-hidden">
                <h2 className="text-lg">Videos</h2>
                {getDirectoryQuery.data.videoPaths.length > 0 ? (
                    <section className="overflow-y-auto mt-2">
                        <section className="flex flex-col gap-2 w-full h-fit">
                            {getDirectoryQuery.data.videoPaths.map(
                                (videoPath) => {
                                    return (
                                        <Video
                                            key={videoPath}
                                            path={videoPath}
                                        />
                                    )
                                }
                            )}
                        </section>
                    </section>
                ) : (
                    <section className="h-full flex items-center justify-center">
                        <p className="text-sm">No videos found!</p>
                    </section>
                )}
            </section>
        </section>
    )
}
