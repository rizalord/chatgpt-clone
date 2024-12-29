'use client' // Error boundaries must be Client Components

import { useEffect } from 'react'

export default function Error({
    error,
}: {
    error: Error & { digest?: string }
    reset: () => void
}) {
    useEffect(() => {
        // Log the error to an error reporting service
        console.error(error)
    }, [error])

    return (
        <div className="flex flex-col h-screen justify-center items-center px-4 gap-2">
            <h1 className="uppercase tracking-widest text-gray-500">Something went wrong!</h1>
            <p className="text-lg text-gray-500">{error.message}</p>
            {error.digest && <p className="text-sm text-gray-500">Digest: {error.digest}</p>}
        </div>
    )
}