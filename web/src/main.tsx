import { Browser } from '@/src/sections/browser'
import { Editor } from '@/src/sections/editor'
import { FileTree } from '@/src/sections/file_tree'
import { Header } from '@/src/sections/header'
import { Jobs } from '@/src/sections/jobs'
import { Settings } from '@/src/sections/settings'
import { Statistics } from '@/src/sections/statistics'

import { Route, useRoute } from './utility/navigation'

export function Main() {
    const currentRoute = useRoute()

    return (
        <section className="h-dvh w-dvw bg-neutral-dark overflow-y-hidden">
            <section
                className={`h-fit w-full flex flex-col transition-transform duration-200 ${currentRoute == Route.Editor && 'translate-y-[-100dvh]'}`}
            >
                <section className="h-dvh w-full p-8 grid grid-rows-[auto_auto_1fr] grid-cols-1">
                    <Header />
                    <Statistics />
                    <section className="overflow-y-hidden overflow-x-hidden">
                        <section
                            className={`h-full w-fit gap-4 grid grid-rows-1 grid-cols-[repeat(3,calc(100dvw-var(--spacing)*8*2))] transition-transform duration-200 ${currentRoute == Route.Jobs && 'translate-x-[calc(-100dvw+var(--spacing)*8*2)]'} ${currentRoute == Route.Settings && 'translate-x-[calc(-200dvw+var(--spacing)*8*3)]'}`}
                        >
                            <section className="grid grid-cols-[calc(var(--spacing)*96)_1fr] gap-4">
                                <FileTree />
                                <Browser />
                            </section>
                            <section className="">
                                <Jobs />
                            </section>
                            <section className="">
                                <Settings />
                            </section>
                        </section>
                    </section>
                </section>
                <section className="h-dvh p-8 pt-0">
                    <Editor />
                </section>
            </section>
        </section>
    )
}
