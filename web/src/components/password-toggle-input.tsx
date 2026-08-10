import { Eye, EyeClosed } from "lucide-react"
import { Input } from "./ui/input"
import { useState, type ComponentProps } from "react"

export default function PasswordInputWithToggle({
  ...props
}: ComponentProps<"input">) {
  const [inputType, setInputType] = useState<"password" | "text">("password")

  const updateInputType = () => {
    if (inputType === "password") {
      setInputType("text")
    } else {
      setInputType("password")
    }
  }
  return (
    <div className="relative w-full">
      <Input {...props} type={inputType} />
      <button
        onClick={updateInputType}
        type="button"
        className="absolute top-1/2 right-3 -translate-y-1/2 cursor-pointer"
      >
        {inputType === "password" ? (
          <Eye onClick={updateInputType} />
        ) : (
          <EyeClosed onClick={updateInputType} />
        )}
      </button>
    </div>
  )
}
