"use client";

import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarTrigger,
} from "./sidebar";
import { usePathname } from "next/navigation";
import { Separator } from "./separator";
import { NavUser } from "./nav-user";
import { JSX } from "react";
import Link from "next/link";
import {
    BookImage,
    Clock4,
    Code,
    Folder,
    Home,
    Plus,
    Star,
    Trash2,
    UsersRound,
} from "lucide-react";
import { Search } from "../search";
import { Button } from "./button";
import Notifications from "../notifications";

const sidebarHiddedRoutes = ["/login", "/signup"];

interface SidebarItem {
    title: string;
    path: string;
    icon: JSX.ElementType;
}

export function AppSidebarHeader() {
    const pathname = usePathname();

    if (sidebarHiddedRoutes.includes(pathname)) return null;

    return (
        <header className='flex h-16 w-full shrink-0 items-center justify-between px-4'>
            <div className='flex items-center gap-2 w-full'>
                <SidebarTrigger className='-ml-1' />
                <Separator orientation='vertical' className='mr-2 !h-4' />
                <Search />
            </div>
            <div className='flex gap-3 items-center'>
                <Button className='w-40'>
                    <Plus />
                    Create Canvas
                </Button>
                <Notifications />
                <NavUser
                    user={{
                        name: "Viper",
                        email: "example@email.com",
                        avatar: "https://cdn.cdavidsv.dev/mml/avatars/5b433311-67d3-4823-9fb6-5b8dc9fc0810.jpeg",
                    }}
                />
            </div>
        </header>
    );
}

export function AppSidebar() {
    const pathname = usePathname();

    if (sidebarHiddedRoutes.includes(pathname)) return null;

    const sidebarItems: SidebarItem[] = [
        {
            title: "Home",
            path: "/",
            icon: Home,
        },
        {
            title: "My Gallery",
            path: "/gallery",
            icon: BookImage,
        },
        {
            title: "My Collections",
            path: "/collections",
            icon: UsersRound,
        },
        {
            title: "Starred",
            path: "/starred",
            icon: Star,
        },
        {
            title: "Shared with me",
            path: "/shared",
            icon: UsersRound,
        },
        {
            title: "Recent",
            path: "/recent",
            icon: Clock4,
        },
    ];

    return (
        <Sidebar collapsible='icon'>
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <h1 className='truncate text-foreground text-2xl transition-opacity duration-[320ms] ease-in-out group-data-[collapsible=icon]:opacity-0'>
                            Pixio
                        </h1>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupLabel>App</SidebarGroupLabel>
                    <SidebarGroupContent>
                        <SidebarMenu>
                            {sidebarItems.map((item) => (
                                <SidebarMenuItem key={item.title}>
                                    <SidebarMenuButton
                                        tooltip={item.title}
                                        className={`${
                                            pathname === item.path
                                                ? "!bg-primary !text-white"
                                                : ""
                                        }`}
                                        asChild
                                    >
                                        <Link
                                            href={item.path}
                                            className='h-12 group-data-[collapsible=icon]:size-12!'
                                        >
                                            <item.icon className='!h-5 !w-5 ml-[0.362rem]' />
                                            <span className='ml-2 text-sm font-[500] text-nowrap transition-opacity duration-[320ms] group-data-[collapsible=icon]:opacity-0'>
                                                {item.title}
                                            </span>
                                        </Link>
                                    </SidebarMenuButton>
                                </SidebarMenuItem>
                            ))}
                        </SidebarMenu>
                    </SidebarGroupContent>
                </SidebarGroup>
                <SidebarGroup />
            </SidebarContent>
            <SidebarFooter></SidebarFooter>
        </Sidebar>
    );
}
