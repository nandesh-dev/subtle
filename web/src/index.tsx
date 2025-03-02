import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import { Api, ApiProvider } from '@/src/utility/api'
import { Navigation, NavigationProvider } from '@/src/utility/navigation'

import './index.css'
import { Main } from './main'

const navigation = new Navigation()
const api = new Api({
    enableMockTransport: import.meta.env.MODE == 'development',
})

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <NavigationProvider navigation={navigation}>
            <ApiProvider api={api}>
                <Main />
            </ApiProvider>
        </NavigationProvider>
    </StrictMode>
)
