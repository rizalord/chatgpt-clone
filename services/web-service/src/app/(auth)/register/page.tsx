import { RegisterForm } from '@/components/form/RegisterForm'

export default function Register() {
    const apiUrl = process.env.PUBLIC_API_URL || process.env.API_URL || ''

    return (
        <div className="flex h-dvh w-screen items-start pt-12 md:pt-0 md:items-center justify-center bg-background">
            <div className="w-full max-w-md overflow-hidden rounded-2xl gap-12 flex flex-col">
                <div className="flex flex-col items-center justify-center gap-2 px-4 text-center sm:px-16">
                    <h3 className="text-xl font-semibold dark:text-zinc-50">Sign Up</h3>
                    <p className="text-sm text-gray-500 dark:text-zinc-400">
                        Create an account with your email and password
                    </p>
                </div>

                <RegisterForm apiUrl={apiUrl} />
            </div>
        </div>
    )
}