import React from "react";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
}

export default function Button({ children, ...props }: ButtonProps) {
  return (
    <button
      {...props}
      className="px-5 py-2 rounded-lg bg-blue-600 text-white font-medium
                 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-400
                 disabled:opacity-50 disabled:cursor-not-allowed transition"
    >
      {children}
    </button>
  );
}
