import * as z from "zod"

export const UserCredentialSchema = z.object({
  username: z
    .string()
    .min(5, "Too short username")
    .max(15, "Too large username"),
  password: z.string().min(6, "Password must be at least 6 characters"),
})

export const UserSchema = z
  .object({
    username: z
      .string()
      .min(5, "Too short username")
      .max(15, "Too large username"),
    password: z.string().min(6, "Password must be at least 6 characters"),
    confirmPassword: z.string().min(6, "Confirm password required"),
    firstName: z.string().min(3, "Name required"),
    lastName: z.string().min(3, "Name required"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  })

export const UserCredentials = z.object({
  username: z
    .string()
    .min(5, "Too short username")
    .max(15, "Too large username"),
  password: z.string().min(6, "Password must be at least 6 characters"),
})
