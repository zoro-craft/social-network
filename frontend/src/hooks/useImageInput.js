'use client';

import { useState } from 'react';

export function useImageInput() {
  const [file, setFile] = useState(null);

  const handleChange = (e) => {
    setFile(e.target.files[0] || null);
  };

  return { file, setFile, handleChange };
}
