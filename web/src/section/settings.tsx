import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
    getConfig,
    updateConfig,
} from '../../gen/proto/web/web-WebService_connectquery'

export function Settings() {
    const getConfigQuery = useQuery(getConfig)
    const updateConfigMutation = useMutation(updateConfig)

    return (
        <div className="relative flex h-full flex-col gap-md overflow-hidden rounded-md bg-neutral-2 p-xl">
            <h2 className="text-lg text-text-1">Settings</h2>
            {!getConfigQuery.isSuccess ? (
                <div />
            ) : (
                <textarea
                    className="h-full w-full rounded-sm bg-neutral-1 p-md text-text-1"
                    defaultValue={getConfigQuery.data.config}
                    onChange={(e) => {
                        const updatedConfig = e.target.value
                        updateConfigMutation.mutate({ updatedConfig })
                    }}
                >
                </textarea>
            )}
        </div>
    )
}
