export interface Authorization {
  token: string
}

export interface User {
  firstName: string
  lastName: string
  username: string
  id: string
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
