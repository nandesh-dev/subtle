import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { Root } from './routes/root'
import { Home } from './routes/home/home'
import { Media } from './routes/media/media'

import { ProtoContent, ProtoContext } from './context/proto'
import { createGrpcWebTransport } from '@connectrpc/connect-web'
import { createClient } from '@connectrpc/connect'
import { MediaService } from '../gen/proto/media/media_connect'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Video } from './routes/media/video/video'

const router = createBrowserRouter([
    {
        path: '/',
        element: <Root />,
        children: [
            {
                path: 'home',
                element: <Home />,
            },
            {
                path: 'media',
                element: <Media />,
            },
            {
                path: 'media/video',
                element: <Video />,
            },
        ],
    },
])

const transport = createGrpcWebTransport({
    baseUrl: 'http://localhost:3000',
})

const proto: ProtoContent = {
    MediaServiceClient: createClient(MediaService, transport),
}

const query = new QueryClient()

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <ProtoContext.Provider value={proto}>
            <QueryClientProvider client={query}>
                <RouterProvider router={router} />
            </QueryClientProvider>
        </ProtoContext.Provider>
    </StrictMode>
)
