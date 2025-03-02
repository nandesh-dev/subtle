import { useQuery } from '@connectrpc/connect-query'

import { getRootDirectoryPaths } from '@/gen/proto/services/web-WebService_connectquery'

import { LoadingBlock } from '@/src/components/loading_block'

import { Directory } from './directory'

export function FileTree() {
    const getRootDirectoryPathsQuery = useQuery(getRootDirectoryPaths)

    if (getRootDirectoryPathsQuery.isError) {
        return (
            <section className="overflow-y-scroll flex flex-col gap-8 items-center justify-center">
                <p>{getRootDirectoryPathsQuery.error.message}</p>
            </section>
        )
    }

    if (getRootDirectoryPathsQuery.isPending) {
        return (
            <section className="overflow-y-scroll flex flex-col gap-8 items-center justify-center">
                <LoadingBlock className="bg-neutral size-16"/>
                <LoadingBlock className="bg-neutral size-16"/>
                <LoadingBlock className="bg-neutral size-16"/>
            </section>
        )
    }

    return (
        <section className="overflow-y-auto pr-4">
            {getRootDirectoryPathsQuery.data.paths.map((path) => {
                return (
                    <Directory
                        key={path}
                        path={path}
                        displayEntirePath={true}
                        isRootDirectory={true}
                    />
                )
            })}
        </section>
    )
}
