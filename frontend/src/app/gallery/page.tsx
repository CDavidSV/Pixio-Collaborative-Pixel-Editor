import CanvasCard from "@/components/canvasCard";
import { Input } from "@/components/ui/input";
import PageHeader from "@/components/ui/page-header";
import { Sorter } from "@/components/ui/sort";
import { ViewToggle } from "@/components/ui/viewToggle";
import { CanvasCardData } from "@/types/common";
import { SearchIcon } from "lucide-react";

export const metadata = {
    title: "My Gallery | Pixio",
    description:
        "Manage your pixel art creations. Rename, duplicate, or delete your canvases with ease.",
};

export default function GalleryPage() {
    const canvases: CanvasCardData[] = [
        {
            id: "1",
            title: "Sunset Overdrive",
            stars: 120,
            owner: {
                id: "user1",
                username: "pixelartist",
                avatar: "",
                url: "",
            },
            description: "A vibrant pixel art representation of a sunset.",
            image: "/img/canvas1.webp",
            createdAt: new Date("2023-10-01T10:00:00Z"),
            updatedAt: new Date("2023-10-02T11:00:00Z"),
        },
        {
            id: "2",
            title: "Retro Racer",
            stars: 85,
            owner: {
                id: "user2",
                username: "retrolover",
                avatar: "",
                url: "",
            },
            description:
                "A pixel art tribute to classic racing games of the 80s.",
            image: "/img/canvas1.webp",
            createdAt: new Date("2023-09-15T08:30:00Z"),
            updatedAt: new Date("2023-09-16T09:45:00Z"),
        },
        {
            id: "3",
            title: "Fantasy Forest",
            stars: 2384,
            owner: {
                id: "user3",
                username: "fantasyfan",
                avatar: "",
                url: "",
            },
            description:
                "An enchanting pixel art scene depicting a magical forest.",
            image: "/img/canvas1.webp",
            createdAt: new Date("2023-08-20T14:20:00Z"),
            updatedAt: new Date("2023-08-21T15:30:00Z"),
        }
    ];

    return (
        <div>
            <PageHeader
                title='My Gallery'
                description='Manage your pixel art creations. Rename, duplicate, or delete your canvases with ease.'
            />
            <div className='flex items-center justify-between'>
                <div className='relative flex items-center gap-2 w-full md:w-auto md:min-w-[25%]'>
                    <SearchIcon className='absolute h-4 left-3 text-muted-foreground' />
                    <Input className='pl-10' placeholder='Search by Title' />
                </div>
                <div className='flex justify-end gap-2'>
                    <Sorter />
                    <ViewToggle />
                </div>
            </div>
            <div className='grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 2xl:grid-cols-7 gap-3 mt-5'>
                {canvases.map((canvas) => (
                    <CanvasCard key={canvas.id} {...canvas} />
                ))}
            </div>
        </div>
    );
}
