"use client";

import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Button } from "./button";
import { ChevronDown, ChevronUp } from "lucide-react";
import { useState } from "react";

export enum SortOptions {
    Name = "name",
    LastModified = "lastModified",
    LastOpened = "lastOpened",
}

export enum SortDirection {
    Newest = "newest",
    Oldest = "oldest",
}

export function Sorter() {
    const [sortOption, setSortOption] = useState<SortOptions>(SortOptions.Name);
    const [sort, setSort] = useState<SortDirection>(SortDirection.Newest);

    return (
        <div className='flex flex-row gap-2'>
            <Button
                variant='secondary'
                onClick={() =>
                    setSort(
                        sort === SortDirection.Newest
                            ? SortDirection.Oldest
                            : SortDirection.Newest
                    )
                }
            >
                {sort === SortDirection.Newest ? (
                    <ChevronDown />
                ) : (
                    <ChevronUp />
                )}
            </Button>
            <Select defaultValue="name" onValueChange={(v) => setSortOption(v as SortOptions)}>
                <SelectTrigger className='w-[150px]'>
                    <SelectValue placeholder='Sort by' />
                </SelectTrigger>
                <SelectContent>
                    <SelectGroup>
                        <SelectLabel>Sort by</SelectLabel>
                        <SelectItem value='name'>
                            Name
                        </SelectItem>
                        <SelectItem value='lastModified'>
                            Last Modified
                        </SelectItem>
                        <SelectItem value='lastOpened'>
                            Last Opened
                        </SelectItem>
                    </SelectGroup>
                </SelectContent>
            </Select>
        </div>
    );
}
