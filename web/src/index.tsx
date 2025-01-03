import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { Main } from './main'
import { Navigation, NavigationProvider } from './utils/navigation'

const navigation = new Navigation()

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <NavigationProvider navigation={navigation}>
            <Main />
        </NavigationProvider>
    </StrictMode>
)
