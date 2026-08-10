import PasswordInputWithToggle from "@/components/password-toggle-input"
import { Button } from "@/components/ui/button"
import { Field, FieldError, FieldLabel, FieldSet } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { motion, type Variants } from "framer-motion"
import { Mail, Sparkles, User } from "lucide-react"
import { useActionState, useEffect, useState } from "react"
import { type AuthMode } from "@/pages/auth.page"
import { CardDescription } from "../ui/card"
import { UserSchema } from "@/schemas/user.schema"
import axios from "axios"

export const SignUpForm = ({
  variants,
  updateFormError,
  updateAuthMode,
}: {
  variants: Variants
  updateAuthMode: (authMode: AuthMode) => void

  updateFormError: (formError: string | undefined) => void
}) => {
  const [formState, action, isLoading] = useActionState<SignUpForm, FormData>(
    async (_state: SignUpForm, formData: FormData) => {
      const userData = {
        phone: formData.get("phone")?.toString(),
        countryCode: formData.get("countryCode")?.toString(),
        email: formData.get("email")?.toString(),
        name: formData.get("fullName")?.toString(),
        password: formData.get("password")?.toString(),
        confirmPassword: formData.get("confirmPassword")?.toString(),
      }

      const result = UserSchema.safeParse(userData)

      const formError: SignUpFormError = {}

      if (result.error) {
        for (const iss of result.error.issues) {
          formError[iss.path[0] as keyof SignUpFormError] = iss.message
        }

        return {
          state: {
            ...userData,
          },
          errors: formError,
        } as SignUpForm
      }

      const userDetails = result.data

      const res = await axios.post(
        `/api/v1/auth/register`,
        {
          email: userDetails.email,
          password: userDetails.password,
          name: userDetails.name,
          phone: userDetails.phone,
          countryCode: userDetails.countryCode,
        },
        {
          validateStatus: () => true,
        }
      )

      const { status, data } = res

      if (status === 201) {
        updateAuthMode("signin")
        return {
          state: {},
          errors: {},
        }
      }

      const serverError: SignUpFormError = {
        err: data.error || "Something went wrong.",
      }

      return {
        state: {
          ...userData,
        },
        errors: serverError,
      } as SignUpForm
    },

    { state: {}, errors: {} }
  )

  const [errors, setErrors] = useState<SignUpFormError>(formState.errors)

  const [prevStateError, setPrevStateError] = useState<SignUpFormError>(
    formState.errors
  )

  if (formState.errors !== prevStateError) {
    setErrors(formState.errors)
    setPrevStateError(formState.errors)
  }

  useEffect(() => {
    updateFormError(errors.err)
  }, [errors.err, updateFormError])

  return (
    <motion.form
      key="signup"
      initial="hidden"
      variants={variants}
      animate="visible"
      exit="exit"
      action={action}
      className="space-y-3.5"
    >
      {errors.err && <CardDescription>{errors.err}</CardDescription>}
      <FieldSet>
        {/* Name Inputs */}
        <Field className="space-y-1.5">
          <FieldLabel htmlFor="first-name">Full Name</FieldLabel>
          <div className="relative">
            <User className="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              name="fullName"
              id="full-name"
              type="text"
              required
              placeholder="John Dow"
              className="pl-9"
            />
          </div>
          {errors.name && <FieldError>{errors.name}</FieldError>}
        </Field>

        {/* Email Input */}
        <Field className="space-y-1.5">
          <Label htmlFor="signup-email">Email address</Label>
          <div className="relative">
            <Mail className="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              id="signup-email"
              type="email"
              name="email"
              required
              placeholder="alex@example.com"
              className="pl-9"
            />
          </div>

          {errors.email && <FieldError>{errors.email}</FieldError>}
        </Field>

        {/* Password Inputs with Custom PasswordWithToggle */}
        <Field className="space-y-1.5">
          <Label htmlFor="signup-password">Password</Label>
          <PasswordInputWithToggle
            id="signup-password"
            required
            name="password"
            placeholder="Min. 8 parameters"
          />

          {errors.password && <FieldError>{errors.password}</FieldError>}
        </Field>

        <Field className="space-y-1.5">
          <Label htmlFor="signup-confirm-password">Confirm password</Label>
          <Input
            type="password"
            id="signup-confirm-password"
            required
            name="confirmPassword"
            placeholder="Repeat password"
          />

          {errors.confirmPassword && (
            <FieldError>{errors.confirmPassword}</FieldError>
          )}
        </Field>

        {/* Submit Button */}
        <Button disabled={isLoading} type="submit" className="w-full gap-1.5">
          <Sparkles className="h-3.5 w-3.5" />
          <span>
            {isLoading ? <span>Signing up ...</span> : <span>Sign Up</span>}
          </span>
        </Button>
      </FieldSet>
    </motion.form>
  )
}

export interface SignUpForm {
  errors: SignUpFormError
  state: SignUpFormFields
}

export interface SignUpFormFields {
  email?: string
  name?: string
  password?: string
  confirmPassword?: string
  phone?: string
  countryCode?: string
}

export interface SignUpFormError {
  err?: string
  email?: string
  password?: string
  confirmPassword?: string
  name?: string
  phone?: string
  countryCode?: string
}
