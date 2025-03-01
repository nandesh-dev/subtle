import { ReactNode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createRouterTransport, Transport } from '@connectrpc/connect'
import { WebService } from '../../../gen/proto/services/web_pb'
import { TransportProvider } from '@connectrpc/connect-query'
import { createGrpcWebTransport } from '@connectrpc/connect-web'
import { getConfig, updateConfig } from './config'
import { getRootDirectoryPaths, readDirectory, searchVideo } from './filesystem'
import {
    getSubtitle,
    getSubtitleCue,
    getSubtitleCueOriginalData,
    getSubtitleCueSegment,
    calculateSubtitleStatistics,
} from './subtitle'

type APIOptions = {
    enableMockTransport?: boolean
}

export class API {
    public rpcTransport: Transport
    public tanstackQueryClient: QueryClient

    constructor(options?: APIOptions) {
        this.tanstackQueryClient = new QueryClient()

        this.rpcTransport = createGrpcWebTransport({
            baseUrl: `${window.location.origin}`,
        })

        if (options?.enableMockTransport) {
            this.rpcTransport = createRouterTransport(({ service }) => {
                service(WebService, {
                    getConfig,
                    updateConfig,

                    getRootDirectoryPaths,
                    readDirectory,
                    searchVideo,

                    calculateSubtitleStatistics,
                    getSubtitle,
                    getSubtitleCue,
                    getSubtitleCueOriginalData,
                    getSubtitleCueSegment,
                })
            })
        }
    }
}

export function APIProvider({
    api,
    children,
}: {
    api: API
    children?: ReactNode
}) {
    return (
        <TransportProvider transport={api.rpcTransport}>
            <QueryClientProvider
                client={api.tanstackQueryClient}
                children={children}
            />
        </TransportProvider>
    )
}
