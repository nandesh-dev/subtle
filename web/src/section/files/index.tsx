import { Browser } from "./browser";
import { Tree } from "./tree";

export function Files() {
    return (
        <section className="grid grid-cols-[auto_1fr] h-full max-h-full overflow-hidden">
            <Tree />
            <Browser/>
        </section>
    )
}
