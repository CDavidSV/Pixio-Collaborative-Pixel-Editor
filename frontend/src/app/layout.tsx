import type { Metadata } from "next";
import { ThemeProvider } from "./theme-provider";
import { AppSidebar, AppSidebarHeader } from "@/components/ui/app-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { cookies } from "next/headers";
import { Poppins } from "next/font/google";
import "./globals.css";

export const metadata: Metadata = {
    title: "Pixio",
    description:
        "A modern pixel art editor with a focus on collaboration and simplicity.",
};

const poppins = Poppins({
    subsets: ["latin"],
    weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
    variable: "--font-poppins",
});

export default async function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    const cookieStore = await cookies();
    const defaultOpen = cookieStore.get("sidebar_state")?.value === "true";

    return (
        <>
            <html lang='en' className={poppins.variable} suppressHydrationWarning>
                <head />
                <body>
                    <ThemeProvider
                        attribute='class'
                        defaultTheme='system'
                        enableSystem
                        disableTransitionOnChange
                    >
                        <SidebarProvider defaultOpen={defaultOpen}>
                            <AppSidebar />
                            <SidebarInset>
                                <AppSidebarHeader />
                                <main className='pr-8 pl-12 py-4'>{children}</main>
                            </SidebarInset>
                        </SidebarProvider>
                    </ThemeProvider>
                </body>
            </html>
        </>
    );
}
