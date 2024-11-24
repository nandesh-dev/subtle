import { Link, Outlet, useLocation } from 'react-router-dom'
import {
    ClockIcon,
    CrossIcon,
    FileIcon,
    HomeIcon,
    HorizontalSubtitleDrop,
    VerticalSubtitleDrop,
} from '../../assets'
import { Desktop, Mobile, Tablet } from '../utils/react_responsive'
import { ReactNode, useState } from 'react'
import { PopoverContent, PopoverContext } from '../context/popover'

export function Root() {
    let location = useLocation()

    let [popoverNode, setPopoverNode] = useState<ReactNode>(null)
    let popoverValue: PopoverContent = {
        set: (node) => {
            setPopoverNode(node)
        },
        reset: () => {
            setPopoverNode(null)
        },
    }

    const navigations = [
        {
            name: 'Home',
            link: '/home',
            icon: HomeIcon,
            ariaLabel: 'Switch to home page',
        },
        {
            name: 'Media',
            link: '/media',
            icon: FileIcon,
            ariaLabel: 'Switch to media page',
        },
        {
            name: 'Routines',
            link: '/routines',
            icon: ClockIcon,
            ariaLabel: 'Switch to routines page',
        },
    ]

    return (
        <>
            <PopoverContext.Provider value={popoverValue}>
                <Mobile>
                    <section className="h-dvh w-dvw bg-gray-50">
                        <section className="h-full w-full px-sm pt-sm">
                            <Outlet />
                        </section>
                        <section className="absolute bottom-0 left-0 p-sm">
                            <nav className="rounded-sm bg-gray-120 px-md py-sm">
                                <ul className="flex flex-row items-center gap-md">
                                    {navigations.map(
                                        ({ link, icon: Icon, ariaLabel }) => {
                                            return (
                                                <li key={link}>
                                                    <Link
                                                        to={link}
                                                        aria-label={ariaLabel}
                                                    >
                                                        <Icon
                                                            className={
                                                                location.pathname.startsWith(
                                                                    link
                                                                )
                                                                    ? 'fill-primary'
                                                                    : 'fill-gray-830'
                                                            }
                                                        />
                                                    </Link>
                                                </li>
                                            )
                                        }
                                    )}
                                </ul>
                            </nav>
                        </section>
                    </section>
                </Mobile>
                <Tablet>
                    <section className="grid h-dvh w-dvw grid-cols-[auto_1fr] bg-gray-50 p-sm">
                        <section className="flex h-full flex-col items-center justify-between rounded-sm bg-gray-120 px-xs py-xl">
                            <h1 className="text-lg text-gray-830">S</h1>
                            <div className="h-[8rem] w-fit rounded-sm bg-gray-190 pt-lg">
                                <div className="h-full w-xs rounded-sm bg-primary"></div>
                            </div>
                            <nav>
                                <ul className="flex flex-col items-center justify-between gap-sm">
                                    {navigations.map(
                                        ({ link, icon: Icon, ariaLabel }) => {
                                            return (
                                                <li key={link}>
                                                    <Link
                                                        to={link}
                                                        aria-label={ariaLabel}
                                                    >
                                                        <Icon
                                                            className={
                                                                location.pathname.startsWith(
                                                                    link
                                                                )
                                                                    ? 'fill-primary'
                                                                    : 'fill-gray-830'
                                                            }
                                                        />
                                                    </Link>
                                                </li>
                                            )
                                        }
                                    )}
                                </ul>
                            </nav>
                            <VerticalSubtitleDrop />
                        </section>
                        <section className="overflow-hidden">
                            <Outlet />
                        </section>
                    </section>
                </Tablet>
                <Desktop>
                    <section className="grid h-dvh w-dvw grid-cols-[16rem_1fr] bg-gray-50 p-sm">
                        <section className="flex h-full flex-col justify-between rounded-sm bg-gray-120 p-xl">
                            <h1 className="text-center text-lg text-gray-830">
                                Subtle
                            </h1>
                            <div className="w-fill relative flex aspect-square items-center justify-center">
                                <svg
                                    viewBox="0 0 240 240"
                                    className="aspect-square w-full stroke-primary"
                                >
                                    <circle
                                        cx="120"
                                        cy="120"
                                        r="104"
                                        strokeWidth="16"
                                        fill="none"
                                        strokeLinecap="round"
                                        strokeDasharray="654"
                                        strokeDashoffset="100"
                                        transform="rotate(90 120 120)"
                                    />
                                </svg>
                                <p className="absolute text-md text-gray-830">
                                    83%
                                </p>
                                <div className="bg-primary"></div>
                            </div>
                            <nav className="">
                                <ul className="flex flex-col justify-between gap-sm">
                                    {navigations.map(({ link, name }) => {
                                        return (
                                            <li key={link}>
                                                <Link
                                                    to={link}
                                                    className={
                                                        'text-sm ' +
                                                        (location.pathname.startsWith(
                                                            link
                                                        )
                                                            ? 'text-primary'
                                                            : 'text-gray-830')
                                                    }
                                                >
                                                    {name}
                                                </Link>
                                            </li>
                                        )
                                    })}
                                </ul>
                            </nav>
                            <div className="outline-sm flex aspect-[3/4] w-full flex-col items-center justify-center rounded-sm outline-dashed outline-gray-190">
                                <HorizontalSubtitleDrop />
                            </div>
                        </section>
                        <section className="overflow-hidden">
                            <Outlet />
                        </section>
                    </section>
                </Desktop>
                {popoverNode && (
                    <section className="absolute bottom-0 left-0 right-0 top-0 flex items-center justify-center bg-gray-50 bg-opacity-75">
                        <div className="relative h-3/4 w-3/4 rounded-xs bg-gray-80 p-sm">
                            <button
                                className="absolute right-sm top-sm"
                                onClick={() => setPopoverNode(null)}
                            >
                                <CrossIcon className="fill-gray-190" />
                            </button>
                            {popoverNode}
                        </div>
                    </section>
                )}
            </PopoverContext.Provider>
        </>
    )
}
