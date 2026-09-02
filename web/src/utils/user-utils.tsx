export function getInitials(firstName: string, lastName: string): string {
  return `${firstName[0]}${lastName[0]}`.toUpperCase()
}

export function getCapitalise(firstName: string, lastName: string): string {
  const cFirstname = firstName[0].toUpperCase() + firstName.substring(1)
  const cLastname = lastName[0].toUpperCase() + lastName.substring(1)

  const final = `${cFirstname} ${cLastname}`

  return final
}
