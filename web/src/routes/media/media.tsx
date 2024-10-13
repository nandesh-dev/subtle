import { useParams } from "react-router-dom";
import { FolderIcon, SearchIcon } from "../../../assets";
import { Large, Small } from "../../utils/react_responsive";

export function Media() {
  const { "*": mediaPath } = useParams();

  return (
    <>
      <Small>
        <section className="flex flex-col h-full gap-sm">
          <section className="flex flex-col gap-sm">
            <SearchBar />
            <div className="flex flex-row items-center gap-lg">
              <h2 className="text-md text-gray-830">Media</h2>
              <p className="text-sm text-gray-520">/media/{mediaPath}</p>
            </div>
          </section>
          <section className="grid overflow-y-auto grid-flow-row gap-sm pb-xxl">
            <Folder name="Movies" subtitle={{ present: 10, total: 20 }} />
            <div className="rounded-sm h-[4px] bg-gray-80" content="2" />
          </section>
        </section>
      </Small>
      <Large>
        <section className="flex flex-col h-full gap-sm py-xxl px-lg">
          <section className="grid grid-cols-[1fr_20rem]">
            <div className="flex flex-row items-center gap-lg">
              <h2 className="text-md text-gray-830">Media</h2>
              <p className="text-sm text-gray-520">/media/{mediaPath}</p>
            </div>
            <SearchBar />
          </section>
          <section className="grid overflow-y-auto grid-flow-row w-full gap-sm">
            <div className="grid grid-cols-[repeat(auto-fill,minmax(18rem,1fr))] gap-sm">
              <Folder name="Movies" subtitle={{ present: 10, total: 20 }} />
            </div>
            <div className="rounded-sm h-[4px] bg-gray-80" />
          </section>
        </section>
      </Large>
    </>
  );
}

type FolderProp = {
  name: string;
  subtitle: {
    present: number;
    total: number;
  };
};

function Folder({ name, subtitle }: FolderProp) {
  return (
    <div className="grid rounded-sm grid-cols-[auto_1fr] p-md bg-gray-80 gap-md">
      <FolderIcon className="h-full fill-red" />
      <div className="">
        <p className="text-sm text-gray-830">{name}</p>
        <div className="flex flex-row justify-between w-full">
          <p className="text-xs text-gray-520">Subtitle</p>
          <p className="text-xs text-gray-520">
            {subtitle.present}/{subtitle.total}
          </p>
        </div>
      </div>
    </div>
  );
}

function SearchBar() {
  return (
    <div className="flex flex-row items-center w-full h-full rounded-md bg-gray-120 gap-sm px-sm py-xs">
      <SearchIcon className="fill-gray-520" />
      <input
        className="w-full text-sm placeholder:text-gray-190 text-gray-830 placeholder:text-sm"
        placeholder="Search"
      />
    </div>
  );
}
