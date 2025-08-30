import React from "react";

export default function Card({ children }: { children: React.ReactNode }) {
  return (
    <div className="w-full max-w-2xl bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mt-6 overflow-x-auto">
      {children}
    </div>
  );
}
