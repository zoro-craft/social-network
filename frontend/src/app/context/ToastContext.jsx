'use client';

import { createContext, useContext, useState } from 'react';
import styles from '../styles/toast.module.css';

const ToastContext = createContext(null);

export function ToastProvider({ children }) {
  const [toast, setToast] = useState(null);

  const showToast = (message, isError = true) => {
    setToast({ message, isError });
    setTimeout(() => setToast(null), 3000);
  };

  return (
    <ToastContext.Provider value={showToast}>
      {children}
      {toast && (
        <div
          className={`${styles.toast} ${toast.isError ? styles.error : styles.success}`}
        >
          {toast.message}
        </div>
      )}
    </ToastContext.Provider>
  );
}

export function useToast() {
  return useContext(ToastContext);
}
