import { RefCallback, useCallback, useRef, useState } from 'react'

export function useIntersectionObserver() {
    const [isIntersecting, setIsIntersecting] = useState(false)
    const previousObserver = useRef<IntersectionObserver | null>(null)

    const observationElementRef = useCallback<RefCallback<HTMLElement>>(
        (node) => {
            if (previousObserver.current) {
                previousObserver.current.disconnect()
                previousObserver.current = null
            }

            if (node?.nodeType === Node.ELEMENT_NODE) {
                const observer = new IntersectionObserver(([entry]) => {
                    setIsIntersecting(entry.isIntersecting)
                })

                observer.observe(node)
                previousObserver.current = observer
            }
        },
        []
    )

    return { isIntersecting, observationElementRef }
}
