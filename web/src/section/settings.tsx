import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
    getConfig,
    updateConfig,
} from '../../gen/proto/web/web-WebService_connectquery'
import { useRef } from 'react'

export function Settings() {
    const getConfigQuery = useQuery(getConfig)
    const updateConfigMutation = useMutation(updateConfig)
    const textAreaRef = useRef<HTMLTextAreaElement>(null)

    return (
        <div className="relative flex h-full flex-col gap-md overflow-hidden rounded-md bg-neutral-2 p-xl">
            <h2 className="text-lg text-text-1">Settings</h2>
            {!getConfigQuery.isSuccess ? (
                <div className="h-full animate-pulse rounded-sm bg-neutral-1" />
            ) : (
                <textarea
                    className="h-full rounded-sm bg-neutral-1 p-md text-text-1"
                    defaultValue={getConfigQuery.data.config}
                    ref={textAreaRef}
                />
            )}
            <section className="flex flex-row gap-md items-center justify-between">
                <button
                    className="w-fit rounded-sm bg-primary-1 px-md py-sm text-xs text-text-2 hover:bg-primary-2 disabled:bg-primary-2"
                    disabled={updateConfigMutation.isPending}
                    onClick={() => {
                        if (textAreaRef.current == null) return
                        updateConfigMutation.mutate({
                            updatedConfig: textAreaRef.current.value,
                        })
                    }}
                >
                    Update Config
                </button>
                {updateConfigMutation.isError && (
                    <p className="text-sm text-text-1 ">{updateConfigMutation.error.rawMessage}</p>
                )}
            </section>
        </div>
    )
}
