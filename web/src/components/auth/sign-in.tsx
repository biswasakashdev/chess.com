import { useActionState, useEffect, useState } from "react"
import { motion, type Variants } from "framer-motion"
import { Mail, ArrowRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import PasswordInputWithToggle from "@/components/password-toggle-input"
import { Field, FieldError } from "@/components/ui/field"
import { UserCredentialSchema } from "@/schemas/user.schema"
import axios from "axios"
import type { AuthMode } from "@/pages/auth.page"
import { useNavigate } from "react-router"

export const SignInForm = ({
  variants,
  updateFormError,
}: {
  variants: Variants
  updateAuthMode: (authMode: AuthMode) => void
  updateFormError: (formError: string | undefined) => void
}) => {
  const navigate = useNavigate()

  const [state, action, isLoading] = useActionState<SignInForm, FormData>(
    async (_prevState: SignInForm, formData: FormData) => {
      const formFields = {
        username: formData.get("username")?.toString() || "",
        password: formData.get("password")?.toString() || "",
      }

      const result = UserCredentialSchema.safeParse(formFields)
      const rememberMe = formData.get("rememberMe") ? true : false

      if (result.error) {
        const formError: SignInFormError = {}
        for (const iss of result.error.issues) {
          formError[iss.path[0] as keyof SignInFormError] = iss.message
        }
        return {
          state: {
            username: formFields.username,
          },
          errors: formError,
        }
      }

      const res = await axios.post(`/api/v1/auth`, result.data, {
        params: {
          rememberMe,
        },
        validateStatus: () => true,
      })

      const { data, status } = res

      if (status !== 201) {
        const err: SignInFormError = {
          err: data.error || "Something went wrong.",
        }
        return {
          state: {
            ...formFields,
            rememberMe,
          },
          errors: err,
        }
      }

      navigate("/home")

      return {
        state: {},
        errors: {},
      }
    },
    {
      state: {},
      errors: {},
    }
  )

  const [errors, setErrors] = useState<SignInFormError>(state.errors)

  const [prevStateError, setPrevStateError] = useState<SignInFormError>(
    state.errors
  )

  if (state.errors !== prevStateError) {
    setErrors(state.errors)
    setPrevStateError(state.errors)
  }

  useEffect(() => {
    updateFormError(errors.err)
  }, [errors.err, updateFormError])

  return (
    <motion.form
      key="signin"
      initial="hidden"
      animate="visible"
      variants={variants}
      exit="exit"
      action={action}
      className="space-y-4"
    >
      {/* Email */}
      <Field className="space-y-1.5">
        <Label htmlFor="signin-email">Email</Label>
        <div className="relative">
          <Mail className="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            id="signin-email"
            type="text"
            name="email"
            required
            defaultValue={state.state.username}
            placeholder="name@example.com"
            className="pl-9"
            onFocus={() =>
              setErrors((pre) => ({ ...pre, email: undefined, err: undefined }))
            }
          />
        </div>
        {errors.username && <FieldError>{errors.username}</FieldError>}
      </Field>

      {/* Password Field */}
      <div className="space-y-1.5">
        <div className="flex items-center justify-between">
          <Label htmlFor="signin-password">Password</Label>
          <a
            href="#"
            className="text-xs text-muted-foreground underline-offset-2 hover:text-foreground hover:underline"
          >
            Forgot?
          </a>
        </div>
        <PasswordInputWithToggle
          id="signin-password"
          name="password"
          required
          placeholder="••••••••"
          onFocus={() =>
            setErrors((pre) => ({
              ...pre,
              password: undefined,
              err: undefined,
            }))
          }
        />
        {errors.password && <FieldError>{errors.password}</FieldError>}
      </div>

      {/* Submit Button */}
      <Button type="submit" className="w-full gap-2">
        <span>{isLoading ? "Signing In ... " : "Sign In"}</span>
        <ArrowRight className="h-3.5 w-3.5" />
      </Button>
    </motion.form>
  )
}

export interface SignInForm {
  state: SignInFormFields
  errors: SignInFormError
}

export interface SignInFormFields {
  username?: string
  rememberMe?: boolean
}

export interface SignInFormError {
  err?: string
  username?: string
  password?: string
}
