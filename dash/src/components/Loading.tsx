import React from 'react';

const Loading: React.FC = () => {
  return (
    <div className="flex h-full w-full flex-col items-center justify-center h-full bg-gray-100 dark:bg-gray-900">
      <div className="w-16 h-16 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
      <p className="mt-4 text-xl font-semibold text-gray-700 dark:text-gray-300">
        Loading...
      </p>
      <p className="mt-2 text-lg text-gray-500 dark:text-gray-400">
        Please wait while we make some api calls
      </p>
    </div>
  );
};

export default Loading;
