import { Route } from "@/src/utility/navigation"
import { Button } from "./button"

export function Navbar() {
    return (
        <section className="flex flex-row rounded-full bg-neutral">
          <Button name="Files" route={Route.Files}/>
          <Button name="Jobs" route={Route.Jobs}/>
          <Button name="Settings" route={Route.Settings}/>
        </section>
    )
}
