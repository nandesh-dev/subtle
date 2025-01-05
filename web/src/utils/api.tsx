import { ReactNode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createRouterTransport, Transport } from '@connectrpc/connect'
import { WebService } from '../../gen/proto/web/web_pb'
import { TransportProvider } from '@connectrpc/connect-query'
import { createGrpcWebTransport } from '@connectrpc/connect-web'

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
                    getGlobalStatistics() {
                        return {
                            Exported: 92,
                            Formated: 41,
                            Extracted: 38,
                            Total: 102,
                        }
                    },
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
