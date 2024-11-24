import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { Root } from './routes/root'
import { Home } from './routes/home/home'
import { Media } from './routes/media/media'
import { Video } from './routes/video/video'
import { Subtitle } from './routes/subtitle/subtitle'

import { ProtoContent, ProtoContext } from './context/proto'
import { createGrpcWebTransport } from '@connectrpc/connect-web'
import { createClient } from '@connectrpc/connect'
import { MediaService } from '../gen/proto/media/media_connect'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { SubtitleService } from '../gen/proto/subtitle/subtitle_connect'
import { Routines } from './routes/routines/routines'
import { RoutineService } from '../gen/proto/routine/routine_connect'

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
                path: 'video',
                element: <Video />,
            },
            {
                path: 'subtitle',
                element: <Subtitle />,
            },
            {
                path: 'routines',
                element: <Routines />,
            },
        ],
    },
])

const transport = createGrpcWebTransport({
    baseUrl: `http://localhost:3000`,
})

const proto: ProtoContent = {
    MediaServiceClient: createClient(MediaService, transport),
    SubtitleServiceClient: createClient(SubtitleService, transport),
    RoutineServiceClient: createClient(RoutineService, transport),
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
