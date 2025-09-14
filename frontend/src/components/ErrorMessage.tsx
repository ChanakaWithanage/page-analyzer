import React from "react";

type Props = {
  message: string;
  className?: string;
};

export default function ErrorMessage({ message, className = "" }: Props) {
  return (
    <div
      role="alert"
      className={
        "w-full rounded-lg border border-red-200/60 bg-red-50 p-4 text-red-900 " +
        "dark:border-red-400/40 dark:bg-red-950/30 dark:text-red-100 " +
        className
      }
    >
      <div className="font-semibold">Error:</div>
      <p className="mt-1">{message}</p>
    </div>
  );
}
