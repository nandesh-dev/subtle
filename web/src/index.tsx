import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { Main } from './main'
import { Navigation, NavigationProvider } from './utils/navigation'
import { API, APIProvider } from './utils/api'

const navigation = new Navigation()
const api = new API({ enableMockTransport: import.meta.env.MODE == "development" })

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <NavigationProvider navigation={navigation}>
            <APIProvider api={api}>
                <Main />
            </APIProvider>
        </NavigationProvider>
    </StrictMode>
)
