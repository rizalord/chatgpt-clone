"use client"

import { FaGoogle } from 'react-icons/fa'
import Link from 'next/link'
import * as Yup from 'yup'
import { ErrorMessage, Field, Form, Formik, FormikHelpers } from 'formik'
import { useState } from 'react'
import { signIn } from 'next-auth/react'
import { useRouter } from 'next/navigation'
import { Response } from '@/types/dto/default.dto'
import { AuthResponse } from '@/types/dto/auth.dto'
import { toast } from 'react-toastify'

interface RegisterFormProps {
  apiUrl: string
}

export function RegisterForm({ apiUrl }: RegisterFormProps) {
  const router = useRouter()
  const [errorMessage, setErrorMessage] = useState<string | null>(null)

  const initialValues = {
    name: '',
    email: '',
    password: '',
  }

  const validationSchema = Yup.object({
    name: Yup.string().required('Name is required'),
    email: Yup.string().email('Invalid email address').required('Email is required'),
    password: Yup.string().required('Password is required').min(6, 'Password must be at least 6 characters long'),
  })

  const onSubmit = async (values: typeof initialValues, { setSubmitting }: FormikHelpers<typeof initialValues>) => {
    setSubmitting(true)
    setErrorMessage(null)

    try {
      const response = await fetch(`${apiUrl}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(values),
      })

      const data: Response<AuthResponse> = await response.json()

      if (response.ok) {
        const signInResponse = await signIn('login', {
          email: values.email,
          password: values.password,
          redirect: false,
        })

        if (signInResponse?.error) {
          setErrorMessage('An error occurred. Please try again later.')
        } else {
          toast.success('Signed up successfully')
          await new Promise(resolve => setTimeout(resolve, 1000))
          router.push('/')
        }
      } else {
        setErrorMessage(data.message || 'An error occurred. Please try again later.')
      }
    } catch (_) {
      setErrorMessage('An error occurred. Please try again later.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <>
      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema}
        onSubmit={onSubmit}
      >
        {
          (props) => (
            <Form className="flex flex-col gap-4 px-4 sm:px-16">
              {
                errorMessage && (
                  <div className="px-3 py-3 flex text-red-600 bg-red-100 rounded-md text-sm">
                    {errorMessage}
                  </div>
                )
              }

              {/* Name */}
              <div className="flex flex-col gap-2">
                <label
                  htmlFor="name"
                  className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 text-zinc-600 font-normal dark:text-zinc-400"
                >
                  Name
                </label>

                <Field
                  id="name"
                  name="name"
                  type="text"
                  placeholder="Enter your name"
                  autoComplete="name"
                  className="flex h-10 w-full rounded-md border border-input px-3 py-2 text-base ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 bg-muted text-md md:text-sm"
                />

                <ErrorMessage name="name" component="div" className="text-red-600 text-sm" />
              </div>
              {/* End Name */}

              {/* Email */}
              <div className="flex flex-col gap-2">
                <label
                  htmlFor="email"
                  className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 text-zinc-600 font-normal dark:text-zinc-400"
                >
                  Email Address
                </label>

                <Field
                  id="email"
                  name="email"
                  type="email"
                  placeholder="user@acme.com"
                  autoComplete="email"
                  autoFocus
                  className="flex h-10 w-full rounded-md border border-input px-3 py-2 text-base ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 bg-muted text-md md:text-sm"
                />

                <ErrorMessage name="email" component="div" className="text-red-600 text-sm" />
              </div>
              {/* End Email */}

              {/* Password */}
              <div className="flex flex-col gap-2">
                <label
                  htmlFor="password"
                  className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 text-zinc-600 font-normal dark:text-zinc-400"
                >
                  Password
                </label>

                <Field
                  id="password"
                  name="password"
                  type="password"
                  placeholder="Enter your password"
                  autoComplete="current-password"
                  className="flex h-10 w-full rounded-md border border-input px-3 py-2 text-base ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 bg-muted text-md md:text-sm"
                />

                <ErrorMessage name="password" component="div" className="text-red-600 text-sm" />
              </div>
              {/* End Password */}

              {/* Submit */}
              <button
                type="submit"
                disabled={props.isSubmitting || !props.isValid}
                className="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2 relative"
              >
                Sign Up
              </button>
              {/* End Submit */}

              {/* Login link */}
              <p className="text-center text-sm text-gray-600 mt-4 dark:text-zinc-400">
                {'Already have an account? '}
                <Link
                  href="/login"
                  className="font-semibold text-gray-800 hover:underline dark:text-zinc-200"
                >
                  Sign in
                </Link>
                {' instead.'}
              </p>
              {/* End Login link */}

              {/* Divider */}
              <div className="flex items-center my-4">
                <div className="flex-grow border-t border-gray-300"></div>
                <span className="mx-4 text-gray-500">or</span>
                <div className="flex-grow border-t border-gray-300"></div>
              </div>
              {/* End Divider */}

              {/* Google Sign Up */}
              <div className="flex flex-col gap-2">
                <button
                  type="button"
                  onClick={() => signIn('google')}
                  className="flex items-center justify-center gap-3 px-4 py-2 bg-[#4285f4] text-white rounded-md"
                >
                  <FaGoogle />

                  <span>Sign up with Google</span>
                </button>
              </div>
              {/* End Google Sign Up */}
            </Form>
          )
        }
      </Formik>
    </>
  )
}
