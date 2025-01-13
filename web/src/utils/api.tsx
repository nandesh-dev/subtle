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
                    getSubtitleSegments({ id }) {
                        return {
                            segments: [
                                {
                                    id: 0,
                                    start: { seconds: BigInt(0) },
                                    end: { seconds: BigInt(10) },
                                    text: 'Hey',
                                },
                                {
                                    id: 1,
                                    start: { seconds: BigInt(10) },
                                    end: { seconds: BigInt(35) },
                                    text: 'Whatsup',
                                },
                                {
                                    id: 2,
                                    start: { seconds: BigInt(35) },
                                    end: { seconds: BigInt(40) },
                                    text: 'Nothing',
                                },
                                {
                                    id: 3,
                                    start: { seconds: BigInt(40) },
                                    end: { seconds: BigInt(60) },
                                    text: '[ Birds Singing ]',
                                },
                                {
                                    id: 4,
                                    start: { seconds: BigInt(60) },
                                    end: { seconds: BigInt(65) },
                                    text: `This are segments of subtitle id ${id}`,
                                },
                            ],
                        }
                    },
                    getSubtitle({ id }) {
                        let name = 'English'
                        switch (id) {
                            case 0:
                                name = 'Forced English'
                                break
                            case 1:
                                name = 'Japanese'
                                break
                            case 2:
                                name = 'Full English'
                                break
                        }

                        return { name }
                    },
                    getVideo({ id }) {
                        let filepath = ''
                        switch (id) {
                            case 0:
                                filepath =
                                    '/media/series/Horimiya/Season 1/Horimiya - S01E01 - A Tiny Happenstance Bluray-1080p.mkv'
                                break
                            case 1:
                                filepath =
                                    '/media/series/Horimiya/Season 1/Horimiya - S01E02 - You Wear More Than One Face Bluray-1080p.mkv'
                                break
                            case 2:
                                filepath =
                                    "/media/series/Horimiya/Season 1/Horimiya - S01E03 - That's Why It's Okay Bluray-1080p.mkv"
                                break
                            case 3:
                                filepath =
                                    '/media/series/Horimiya/Season 1/Horimiya - S01E04 - Everybody Loves Somebody Bluray-1080p.mkv'
                                break
                        }

                        return { filepath, subtitleIds: [0, 1, 2] }
                    },
                    getDirectory({ path }) {
                        if (path == '/media/series') {
                            return {
                                childrenDirectoryNames: ['Horimiya'],
                            }
                        }

                        if (path == '/media/series/Horimiya') {
                            return {
                                childrenDirectoryNames: [
                                    'Season 1',
                                    'Season 2',
                                ],
                            }
                        }

                        return {
                            videoIds: [0, 1, 2, 3],
                        }
                    },
                    getMediaDirectories() {
                        return {
                            paths: ['/media/series'],
                        }
                    },
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
