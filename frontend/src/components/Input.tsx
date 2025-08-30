import React from "react";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {}

export default function Input(props: InputProps) {
  return (
    <input
      {...props}
      className="flex-1 px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-700
                 focus:ring-2 focus:ring-blue-500 focus:outline-none
                 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
    />
  );
}
