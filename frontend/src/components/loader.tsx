import { Loader2 } from "lucide-react";

export function Loader() {

  return (
    <div className="w-full h-screen flex items-center justify-center">
      <Loader2 className="animate-spin h-8 w-8" />
    </div>
  )
}
