import { auth } from '@/auth'

export default auth((req) => {
    const { pathname, origin } = req.nextUrl

    if (!req.auth && (pathname === "/" || pathname.startsWith("/chats/"))) {
        const newUrl = new URL("/login", origin)
        return Response.redirect(newUrl)
    }

    if (req.auth && (pathname === "/login" || pathname === "/register")) {
        const newUrl = new URL("/", origin)
        return Response.redirect(newUrl)
    }
})

// Optionally, don't invoke Middleware on some paths
export const config = {
    matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
}