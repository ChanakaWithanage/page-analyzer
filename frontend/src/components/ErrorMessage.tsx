import React from "react";

export default function ErrorMessage({ message }: { message: string }) {
  return (
    <div
      className="w-full max-w-xl bg-red-100 border border-red-400 text-red-700
                 px-4 py-3 rounded-lg mb-4"
      role="alert"
    >
      <span className="block font-medium">Error:</span> {message}
    </div>
  );
}
