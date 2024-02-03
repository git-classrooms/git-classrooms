import axios from "axios"
import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export async function isAuthenticated() {
  try {
    await axios
      .get("/api/me", { withCredentials: true })

    return true
  } catch (e) {
    return false
  }

}
