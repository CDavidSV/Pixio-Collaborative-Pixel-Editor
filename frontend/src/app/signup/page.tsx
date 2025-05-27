import { Button } from "@/components/ui/button";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PInput } from "@/components/ui/p-input";
import Link from "next/link";
import Image from "next/image";

export const metadata = {
    title: "Signup - Snipari",
    description:
        "Create a Snipari account to save and share your code snippets.",
};

export default function SignupPage() {
    return (
        <div className='flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10'>
            <div className='flex w-full max-w-sm flex-col gap-6'>
                <a
                    href='#'
                    className='flex items-center gap-2 self-center font-medium font-(family-name:--font-pixelify)'
                >
                    <Image src={"/logo/Pixio_Logo.webp"} width={24} height={24} alt="Pixio" />
                    Pixio
                </a>
                <div className='flex flex-col gap-6'>
                    <Card>
                        <CardHeader className='text-center'>
                            <CardTitle className='text-xl'>
                                Create Account
                            </CardTitle>
                            <CardDescription>
                                Sign up with your Google account
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <form>
                                <div className='grid gap-6'>
                                    <div className='flex flex-col gap-4'>
                                        <Button
                                            type='button'
                                            variant='outline'
                                            className='w-full'
                                        >
                                            <svg
                                                xmlns='http://www.w3.org/2000/svg'
                                                viewBox='0 0 24 24'
                                            >
                                                <path
                                                    d='M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z'
                                                    fill='currentColor'
                                                />
                                            </svg>
                                            Sign up with Google
                                        </Button>
                                    </div>
                                    <div className='relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border'>
                                        <span className='relative z-10 bg-card px-2 text-muted-foreground'>
                                            Or continue with
                                        </span>
                                    </div>
                                    <div className='grid gap-6'>
                                        <div className='grid gap-2'>
                                            <Label htmlFor='email'>Email</Label>
                                            <Input
                                                id='email'
                                                type='email'
                                                placeholder='email@example.com'
                                                required
                                            />
                                        </div>
                                        <div className='grid gap-2'>
                                            <Label htmlFor='username'>
                                                Username
                                            </Label>
                                            <Input
                                                id='username'
                                                type='text'
                                                placeholder='Your username'
                                                required
                                            />
                                            <p className='text-muted-foreground text-sm'>
                                                This is your public display name
                                            </p>
                                        </div>
                                        <div className='grid gap-2'>
                                            <div className='flex items-center'>
                                                <Label htmlFor='password'>
                                                    Password
                                                </Label>
                                            </div>
                                            <PInput id='password' required />
                                            <p className='text-muted-foreground text-sm'>
                                                Password must be at least 8
                                                characters long
                                            </p>
                                        </div>
                                        <div className='grid gap-2'>
                                            <div className='flex items-center'>
                                                <Label htmlFor='confirm-password'>
                                                    Confirm Password
                                                </Label>
                                            </div>
                                            <PInput
                                                id='confirm-password'
                                                required
                                            />
                                        </div>
                                        <Button
                                            type='submit'
                                            className='w-full'
                                        >
                                            Signup
                                        </Button>
                                    </div>
                                    <div className='text-center text-sm'>
                                        Already have an account?{" "}
                                        <Link
                                            href='/login'
                                            className='underline underline-offset-4'
                                        >
                                            Log in
                                        </Link>
                                    </div>
                                </div>
                            </form>
                        </CardContent>
                    </Card>
                    {/* <div className='text-balance text-center text-xs text-muted-foreground [&_a]:underline [&_a]:underline-offset-4 [&_a]:hover:text-primary  '>
            By clicking continue, you agree to our{" "}
            <a href='#'>Terms of Service</a> and <a href='#'>Privacy Policy</a>.
          </div> */}
                </div>
            </div>
        </div>
    );
}
