"use client"

export default function GlobalError({ error }: { error: Error & { digest?: string } }) {
  return (
    <html>
      <body>
        <div className="flex flex-col items-center justify-center min-h-screen bg-red-100 text-red-800">
          <h1 className="text-4xl font-bold mb-4">Something went wrong!</h1>
          <p className="text-lg mb-2">{error.message}</p>
          {error.digest && <p className="text-sm">Digest: {error.digest}</p>}
        </div>
      </body>
    </html>
  )
}
