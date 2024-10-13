import { ReactNode } from "react";
import { useMediaQuery } from "react-responsive";

export function Mobile({ children }: { children: ReactNode }) {
  const isDesktop = useMediaQuery({ maxWidth: "640px" });
  return isDesktop ? children : null;
}

export function Tablet({ children }: { children: ReactNode }) {
  const isDesktop = useMediaQuery({ minWidth: "640px", maxWidth: "1024px" });
  return isDesktop ? children : null;
}

export function Desktop({ children }: { children: ReactNode }) {
  const isDesktop = useMediaQuery({ minWidth: "1024px" });
  return isDesktop ? children : null;
}

export function Small({ children }: { children: ReactNode }) {
  const isDesktop = useMediaQuery({ maxWidth: "640px" });
  return isDesktop ? children : null;
}

export function Large({ children }: { children: ReactNode }) {
  const isDesktop = useMediaQuery({ minWidth: "640px" });
  return isDesktop ? children : null;
}
