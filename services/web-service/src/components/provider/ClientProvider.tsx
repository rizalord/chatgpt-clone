"use client"

import React from 'react'
import { Bounce, ToastContainer } from 'react-toastify'

export default function ClientProvider({
    children,
}: Readonly<{
    children: React.ReactNode
}>) {
    return (
        <>
            {children}
            <ToastContainer
                position="bottom-center"
                autoClose={1500}
                hideProgressBar={true}
                newestOnTop={false}
                closeOnClick={true}
                rtl={false}
                draggable
                theme="light"
                transition={Bounce}
            />
        </>
    )
}