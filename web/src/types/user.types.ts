export interface Authorization {
  token: string
}

export interface User {
  firstName: string
  lastName: string
  email: string
  avatar: string
}

export interface UserCredentials {
  email: string
  password: string
}

export interface UserDetails {
  email: string
  firstName: string
  lastName: string
  password: string
}
