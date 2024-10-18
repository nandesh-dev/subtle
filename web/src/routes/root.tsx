import { Link, Outlet, useLocation } from 'react-router-dom'
import { FileIcon, HomeIcon } from '../../assets'
import { Desktop, Mobile, Tablet } from '../utils/react_responsive'

export function Root() {
    let location = useLocation()

    return (
        <>
            <Mobile>
                <section className="h-dvh w-dvw bg-gray-50">
                    <section className="h-full w-full px-sm pt-sm">
                        <Outlet />
                    </section>
                    <section className="absolute bottom-0 w-full p-sm">
                        <nav className="w-full rounded-sm bg-gray-120 px-md py-sm">
                            <ul className="flex flex-row justify-between">
                                <li>
                                    <Link to={'/home'}>
                                        <HomeIcon
                                            className={
                                                location.pathname.startsWith(
                                                    '/home'
                                                )
                                                    ? 'fill-primary'
                                                    : 'fill-gray-830'
                                            }
                                        />
                                    </Link>
                                </li>
                                <li>
                                    <Link to={'/media'}>
                                        <FileIcon
                                            className={
                                                location.pathname.startsWith(
                                                    '/media'
                                                )
                                                    ? 'fill-primary'
                                                    : 'fill-gray-830'
                                            }
                                        />
                                    </Link>
                                </li>
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
                                <li>
                                    <Link to={'/home'}>
                                        <HomeIcon
                                            className={
                                                location.pathname.startsWith(
                                                    '/home'
                                                )
                                                    ? 'fill-primary'
                                                    : 'fill-gray-830'
                                            }
                                        />
                                    </Link>
                                </li>
                                <li>
                                    <Link to={'/media'}>
                                        <FileIcon
                                            className={
                                                location.pathname.startsWith(
                                                    '/media'
                                                )
                                                    ? 'fill-primary'
                                                    : 'fill-gray-830'
                                            }
                                        />
                                    </Link>
                                </li>
                            </ul>
                        </nav>
                        <p className="text-vertical text-sm text-gray-190">
                            Drop Your Subtitle
                        </p>
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
                                    stroke-width="16"
                                    fill="none"
                                    stroke-linecap="round"
                                    stroke-dasharray="654"
                                    stroke-dashoffset="100"
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
                                <li>
                                    <Link
                                        to={'/home'}
                                        className={
                                            'text-sm ' +
                                            (location.pathname.startsWith(
                                                '/home'
                                            )
                                                ? 'text-primary'
                                                : 'text-gray-830')
                                        }
                                    >
                                        Home
                                    </Link>
                                </li>
                                <li>
                                    <Link
                                        to={'/media'}
                                        className={
                                            'text-sm ' +
                                            (location.pathname.startsWith(
                                                '/media'
                                            )
                                                ? 'text-primary'
                                                : 'text-gray-830')
                                        }
                                    >
                                        Media
                                    </Link>
                                </li>
                            </ul>
                        </nav>
                        <div className="outline-sm flex aspect-[3/4] w-full flex-col items-center justify-center rounded-sm outline-dashed outline-gray-190">
                            <p className="text-sm text-gray-190">Drop</p>
                            <p className="text-sm text-gray-190">Your</p>
                            <p className="text-sm text-gray-190">Subtitle</p>
                        </div>
                    </section>
                    <section className="overflow-hidden">
                        <Outlet />
                    </section>
                </section>
            </Desktop>
        </>
    )
}
