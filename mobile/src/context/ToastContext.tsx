import { createContext, useCallback, useContext, useMemo, useState } from 'react';
import { StyleSheet } from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import Animated, { FadeInUp, FadeOutUp } from 'react-native-reanimated';
import { colors, fonts, radii, spacing } from '@/lib/theme';

interface ToastItem {
  id: number;
  message: string;
  isError: boolean;
}

type ShowToast = (message: string, isError?: boolean) => void;

const ToastContext = createContext<ShowToast>(() => {});

let nextId = 0;

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toast, setToast] = useState<ToastItem | null>(null);
  const insets = useSafeAreaInsets();

  const showToast = useCallback<ShowToast>((message, isError = true) => {
    const id = ++nextId;
    setToast({ id, message, isError });
    setTimeout(() => {
      setToast((current) => (current?.id === id ? null : current));
    }, 3000);
  }, []);

  const value = useMemo(() => showToast, [showToast]);

  return (
    <ToastContext.Provider value={value}>
      {children}
      {toast && (
        <Animated.View
          key={toast.id}
          entering={FadeInUp}
          exiting={FadeOutUp}
          pointerEvents="none"
          style={[
            styles.toast,
            { top: insets.top + spacing.md },
            toast.isError ? styles.toastError : styles.toastSuccess,
          ]}
        >
          <Animated.Text style={styles.text}>{toast.message}</Animated.Text>
        </Animated.View>
      )}
    </ToastContext.Provider>
  );
}

export function useToast(): ShowToast {
  return useContext(ToastContext);
}

const styles = StyleSheet.create({
  toast: {
    position: 'absolute',
    alignSelf: 'center',
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderRadius: radii.base,
    borderWidth: 1,
    maxWidth: '90%',
  },
  toastError: {
    backgroundColor: colors.bgCard,
    borderColor: colors.dangerBorder,
  },
  toastSuccess: {
    backgroundColor: colors.bgCard,
    borderColor: colors.success,
  },
  text: {
    color: colors.textPrimary,
    fontFamily: fonts.sansMedium,
    fontSize: 14,
    textAlign: 'center',
  },
});