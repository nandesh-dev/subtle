import { Transport, createRouterTransport } from '@connectrpc/connect'
import { TransportProvider } from '@connectrpc/connect-query'
import { createGrpcWebTransport } from '@connectrpc/connect-web'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactNode } from 'react'

import { WebService } from '../../../gen/proto/services/web_pb'
import { getConfig, updateConfig } from './config'
import { getRootDirectoryPaths, readDirectory, searchVideo } from './filesystem'
import { mockApiDelayInterceptor } from './mock_api_delay_interceptor'
import {
    calculateSubtitleStatistics,
    getSubtitle,
    getSubtitleCue,
    getSubtitleCueOriginalData,
    getSubtitleCueSegment,
} from './subtitle'

type ApiOptions = {
    enableMockTransport?: boolean
}

export class Api {
    public rpcTransport: Transport
    public tanstackQueryClient: QueryClient

    constructor(options?: ApiOptions) {
        this.tanstackQueryClient = new QueryClient()

        this.rpcTransport = createGrpcWebTransport({
            baseUrl: `${window.location.origin}`,
        })

        if (options?.enableMockTransport) {
            this.rpcTransport = createRouterTransport(
                ({ service }) => {
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
                },
                { transport: { interceptors: [mockApiDelayInterceptor] } }
            )
        }
    }
}

export function ApiProvider({
    api,
    children,
}: {
    api: Api
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
