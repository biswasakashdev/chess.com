import { SignInForm } from "@/components/auth/sign-in"
import { SignUpForm } from "@/components/auth/sign-up"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { AnimatePresence, type Variants } from "framer-motion"
import { useState } from "react"

const formTransition: Variants = {
  hidden: { opacity: 0, x: -10 },
  visible: { opacity: 1, x: 0, transition: { duration: 0.2 } },
  exit: { opacity: 0, x: 10, transition: { duration: 0.15 } },
}

export type AuthMode = "signin" | "signup"

export const AuthPage = () => {
  const [authMode, setAuthMode] = useState<AuthMode>("signin")
  const [formErr, setFormError] = useState<string | undefined>(undefined)

  const updateFormError = (formError: string | undefined) => {
    setFormError(formError)
  }

  const updateAuthMode = (authMode: AuthMode) => {
    setAuthMode(authMode)
  }

  return (
    <div className="flex min-h-screen flex-col justify-between bg-background text-foreground antialiased">
      {/* Top Navbar */}
      <header className="sticky top-0 z-50 flex w-full justify-center border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
        <div className="container flex h-16 items-center justify-between px-4 sm:px-8">
          <img src="/logo.png" alt="logo" className="h-full" />
        </div>
      </header>

      {/* 2. Main Authentication Container */}
      <main className="flex flex-1 items-center justify-center p-4">
        <Card className="w-full max-w-110 shadow-sm">
          <CardHeader className="space-y-1 pb-4 text-left">
            <CardTitle className="text-xl font-semibold">
              {authMode === "signin"
                ? "Login to instance"
                : "Create merchant seat"}
            </CardTitle>
            <CardDescription className="text-xs">
              {formErr ? (
                <span className="font-semibold text-red-600">{formErr}</span>
              ) : (
                <span>
                  {authMode === "signin"
                    ? "New to the ecosystem?"
                    : "Already mapped your handles?"}{" "}
                  <button
                    type="button"
                    onClick={() =>
                      setAuthMode((prev) =>
                        prev === "signin" ? "signup" : "signin"
                      )
                    }
                    className="cursor-pointer font-medium text-foreground underline underline-offset-4 focus:outline-none"
                  >
                    {authMode === "signin" ? "Join" : "Sign in"}
                  </button>
                </span>
              )}
            </CardDescription>
          </CardHeader>

          <CardContent className="space-y-6">
            <AnimatePresence mode="wait">
              {authMode === "signin" ? (
                <SignInForm
                  variants={formTransition}
                  updateAuthMode={updateAuthMode}
                  updateFormError={updateFormError}
                />
              ) : (
                <SignUpForm
                  variants={formTransition}
                  updateAuthMode={updateAuthMode}
                  updateFormError={updateFormError}
                />
              )}
            </AnimatePresence>
          </CardContent>
        </Card>
      </main>

    </div>
  )
}

export default AuthPage
