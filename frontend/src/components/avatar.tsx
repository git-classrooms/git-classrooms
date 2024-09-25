import React, { useMemo } from "react";
import { AvatarFallback, AvatarImage, Avatar as AvatarShad } from "@/components/ui/avatar";
import { Avatar as AvatarRadix } from "@radix-ui/react-avatar";

export interface AvatarProps {
  avatarUrl: string | undefined;
  fallbackUrl: string | undefined;
  name: string;
}

export const Avatar: React.FC<AvatarProps> = ({ avatarUrl, fallbackUrl, name }) => {
  const fallbackInitials = useMemo(() => {
    const parts = name.split(" ");
    if (parts.length === 1) {
      return parts[0][0] + parts[0][1];
    } else {
      return parts[0][0] + parts.pop()![0];
    }
  }, [name]);

  return (
    <AvatarShad>
      <AvatarImage src={avatarUrl ?? ""} />
      <AvatarFallback>
        <AvatarRadix>
          <AvatarImage src={fallbackUrl ?? ""} />
          <AvatarFallback>{fallbackInitials}</AvatarFallback>
        </AvatarRadix>
      </AvatarFallback>
    </AvatarShad>
  );
};
