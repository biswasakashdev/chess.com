export interface Authorization {
  token: string
}

export interface User {
  id: string
  username: string
  firstName: string
  lastName: string
}

export interface UserCredentials {
  username: string
  password: string
}

export interface UserDetails {
  username: string
  firstName: string
  lastName: string
  password: string
}
