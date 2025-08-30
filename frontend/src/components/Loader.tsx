import React from "react";

export default function Loader({ size = "md" }: { size?: "sm" | "md" }) {
  const sizeClasses = size === "sm" ? "h-4 w-4 border-2" : "h-6 w-6 border-3";

  return (
    <div
      className={`${sizeClasses} border-white border-t-transparent
                  rounded-full animate-spin`}
    />
  );
}
