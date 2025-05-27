"use client"

import { CanvasCardData } from "@/types/common";
import Link from "next/link";
import Image from "next/image";
import { getRelativeTime } from "@/util/helpers";
import { useState } from "react";
import { Button } from "./ui/button";
import { EllipsisVertical, Star } from "lucide-react";

export default function CanvasCard(canvas: CanvasCardData) {
    const [hovered, setHovered] = useState(false);

    return (
        <Link
            href={`/canvas/${canvas.id}`}
            className='relative'
            onPointerEnter={() => setHovered(true)}
            onPointerLeave={() => setHovered(false)}
        >
            <div className={`absolute top-0 right-0 z-10 m-2 p-1 rounded-md flex items-center gap-2 bg-card/90 transition-all duration-300 ease-in-out ${hovered ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-2'}`}>
                <Button variant="ghost">
                    <Star />
                </Button>
                <Button variant="ghost">
                    <EllipsisVertical />
                </Button>
            </div>
            <div className='rounded-lg overflow-hidden bg-black'>
                <Image
                    className={`w-full h-full aspect-square object-cover transition-all duration-300 ease-in-out ${hovered && 'opacity-80 shadow-2xl'}`}
                    src={canvas.image}
                    width={300}
                    height={300}
                    alt={canvas.title}
                />
            </div>
            <p className='text-lg font-semibold mt-1.5'>{canvas.title}</p>
            <div className='flex items-center justify-start gap-1 text-sm text-muted-foreground'>
                <p>Last edited {getRelativeTime(canvas.updatedAt)}</p>
            </div>
        </Link>
    );
}
