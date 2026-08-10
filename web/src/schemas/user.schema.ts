import * as z from "zod"

export const UserCredentialSchema = z.object({
  email: z.email("Invalid email address"),
  password: z.string().min(6, "Password must be at least 6 characters"),
})

export const UserSchema = z
  .object({
    email: z.email("Invalid email address"),
    countryCode: z.string().min(2, "Country code required"),
    phone: z.string().min(5, "Phone number required"),
    password: z.string().min(6, "Password must be at least 6 characters"),
    confirmPassword: z.string().min(6, "Confirm password required"),
    name: z.string().min(2, "Name required"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  })

export const UserCredentials = z.object({
  email: z.email("Invalid email address"),
  password: z
    .string()
    .min(6, "Password must be at least 6 characters")
    .max(100, "Password too long"),
})
