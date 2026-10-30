import { useEffect } from 'react';
import { DarkTheme, Stack, ThemeProvider, useRouter } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';
import { StyleSheet } from 'react-native';
import * as SplashScreen from 'expo-splash-screen';
import { useFonts } from 'expo-font';
import {
  Inter_400Regular,
  Inter_500Medium,
  Inter_600SemiBold,
  Inter_700Bold,
  Inter_800ExtraBold,
} from '@expo-google-fonts/inter';
import {
  JetBrainsMono_500Medium,
  JetBrainsMono_600SemiBold,
} from '@expo-google-fonts/jetbrains-mono';
import { AuthProvider, useAuth } from '@/context/AuthContext';
import { RealtimeProvider } from '@/context/RealtimeContext';
import { ToastProvider } from '@/context/ToastContext';
import { registerSessionLost } from '@/lib/session';
import { setSessionToken } from '@/lib/api';
import { colors } from '@/lib/theme';

// Keep the native splash visible until the web-matching fonts are loaded.
SplashScreen.preventAutoHideAsync().catch(() => {
  // ignore — splash may already be hidden in some environments
});

// When the server kills the session (or it expires), clear auth and go to /login.
function SessionBridge() {
  const { setAuth } = useAuth();
  const router = useRouter();

  useEffect(() => {
    const unsub = registerSessionLost(() => {
      void setSessionToken(null);
      setAuth({ authenticated: false, user: null });
      router.replace('/login');
    });
    return unsub;
  }, [setAuth, router]);

  return null;
}

export default function RootLayout() {
  const [fontsLoaded, fontError] = useFonts({
    Inter_400Regular,
    Inter_500Medium,
    Inter_600SemiBold,
    Inter_700Bold,
    Inter_800ExtraBold,
    JetBrainsMono_500Medium,
    JetBrainsMono_600SemiBold,
  });

  useEffect(() => {
    if (fontsLoaded || fontError) {
      SplashScreen.hideAsync().catch(() => {
        // ignore
      });
    }
  }, [fontsLoaded, fontError]);

  // Don't render until fonts resolve so text never flashes the system font.
  if (!fontsLoaded && !fontError) return null;

  return (
    <SafeAreaProvider>
      <ThemeProvider value={DarkTheme}>
        <StatusBar style="light" />
        <AuthProvider>
          <RealtimeProvider>
            <ToastProvider>
              <SessionBridge />
              <SafeAreaView style={styles.safeArea} edges={['top']}>
                <Stack
                  screenOptions={{
                    headerShown: false,
                    contentStyle: { backgroundColor: colors.bg },
                  }}
                />
              </SafeAreaView>
            </ToastProvider>
          </RealtimeProvider>
        </AuthProvider>
      </ThemeProvider>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: colors.bg,
  },
});
