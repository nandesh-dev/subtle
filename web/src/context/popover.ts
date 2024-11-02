import { createContext, ReactNode, useContext } from 'react'

export type PopoverContent = {
    node?: ReactNode
    set: (arg0: ReactNode) => void
    reset: () => void
}

export const PopoverContext = createContext<PopoverContent>({
    set: () => {},
    reset: () => {},
})

export const usePopover = () => useContext(PopoverContext)
