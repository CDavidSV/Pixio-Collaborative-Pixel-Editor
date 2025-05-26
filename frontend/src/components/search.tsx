import { SearchIcon } from "lucide-react";
import { Input } from "@/components/ui/input"

export function Search() {
    return (
        <div className='relative flex items-center gap-2 w-full md:w-auto md:min-w-[25%]'>
            <SearchIcon className='absolute h-4 left-3 text-muted-foreground' />
            <Input className="pl-10" placeholder="Search Pixio" />
        </div>
    );
}
