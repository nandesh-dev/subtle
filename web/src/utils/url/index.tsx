import { useState } from 'react'

function parse<T>(value: string | null): T | null {
    if (value == null) return null

    try {
        return JSON.parse(value)
    } catch {
        return value as T
    }
}

function serialize<T>(value: T): string {
    if (typeof value == 'object') {
        return JSON.stringify(value)
    } else if (typeof value == 'number') {
        return value.toString()
    } else if (typeof value == 'string') {
        return value
    }

    throw new Error(`cannot serialize variable: ${value}`)
}

type SetParameterState<T> = (value: T | null) => void

export function useParameterState<T>(
    key: string,
    defaultValue: T | null = null
): [T | null, SetParameterState<T>] {
    const value =
        parse<T>(new URLSearchParams(window.location.search).get(key)) ||
        defaultValue

    const [state, setState] = useState<T | null>(value)

    const updateFunction: SetParameterState<T> = (value) => {
        setState(value)

        const url = new URLSearchParams(window.location.search)
        url.set(key, serialize(value))

        window.history.replaceState(
            null,
            '',
            `${window.location.origin}${window.location.pathname}?${url.toString()}`
        )
    }

    return [state, updateFunction]
}
