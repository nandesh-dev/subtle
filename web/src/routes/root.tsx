import { Link, Outlet, useLocation } from "react-router-dom";
import { FileIcon, HomeIcon } from "../../assets";
import { Desktop, Mobile, Tablet } from "../utils/react_responsive";

export function Root() {
  let location = useLocation();

  return (
    <>
      <Mobile>
        <section className="bg-gray-50 h-dvh w-dvw">
          <section className="w-full h-full px-sm pt-sm">
            <Outlet />
          </section>
          <section className="absolute bottom-0 w-full p-sm">
            <nav className="w-full rounded-sm py-sm px-md bg-gray-120">
              <ul className="flex flex-row justify-between">
                <li>
                  <Link to={"/home"}>
                    <HomeIcon
                      className={
                        location.pathname.startsWith("/home")
                          ? "fill-primary"
                          : "fill-gray-830"
                      }
                    />
                  </Link>
                </li>
                <li>
                  <Link to={"/media"}>
                    <FileIcon
                      className={
                        location.pathname.startsWith("/media")
                          ? "fill-primary"
                          : "fill-gray-830"
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
        <section className="grid bg-gray-50 w-dvw h-dvh p-sm grid-cols-[auto_1fr]">
          <section className="flex flex-col justify-between items-center h-full rounded-sm bg-gray-120 px-xs py-xl">
            <h1 className="text-lg text-gray-830">S</h1>
            <div className="rounded-sm w-fit h-[8rem] pt-lg bg-gray-190">
              <div className="h-full rounded-sm bg-primary w-xs"></div>
            </div>
            <nav>
              <ul className="flex flex-col justify-between items-center gap-sm">
                <li>
                  <Link to={"/home"}>
                    <HomeIcon
                      className={
                        location.pathname.startsWith("/home")
                          ? "fill-primary"
                          : "fill-gray-830"
                      }
                    />
                  </Link>
                </li>
                <li>
                  <Link to={"/media"}>
                    <FileIcon
                      className={
                        location.pathname.startsWith("/media")
                          ? "fill-primary"
                          : "fill-gray-830"
                      }
                    />
                  </Link>
                </li>
              </ul>
            </nav>
            <p className="text-sm text-vertical text-gray-190">
              Drop Your Subtitle
            </p>
          </section>
          <section className="overflow-hidden">
            <Outlet />
          </section>
        </section>
      </Tablet>
      <Desktop>
        <section className="grid bg-gray-50 w-dvw h-dvh p-sm grid-cols-[16rem_1fr]">
          <section className="flex flex-col justify-between h-full rounded-sm bg-gray-120 p-xl">
            <h1 className="text-lg text-center text-gray-830">Subtle</h1>
            <div className="flex relative justify-center items-center w-fill aspect-square">
              <svg
                viewBox="0 0 240 240"
                className="w-full aspect-square stroke-primary"
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
              <p className="absolute text-gray-830 text-md">83%</p>
              <div className="bg-primary"></div>
            </div>
            <nav className="">
              <ul className="flex flex-col justify-between gap-sm">
                <li>
                  <Link
                    to={"/home"}
                    className={
                      "text-sm " +
                      (location.pathname.startsWith("/home")
                        ? "text-primary"
                        : "text-gray-830")
                    }
                  >
                    Home
                  </Link>
                </li>
                <li>
                  <Link
                    to={"/media"}
                    className={
                      "text-sm " +
                      (location.pathname.startsWith("/media")
                        ? "text-primary"
                        : "text-gray-830")
                    }
                  >
                    Media
                  </Link>
                </li>
              </ul>
            </nav>
            <div className="flex flex-col justify-center items-center w-full rounded-sm aspect-[3/4] outline-dashed outline-gray-190 outline-sm">
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
  );
}
