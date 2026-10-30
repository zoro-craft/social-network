import { useState } from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { Link, useRouter } from 'expo-router';
import Screen from '@/components/Screen';
import Card from '@/components/Card';
import AppInput from '@/components/AppInput';
import AppButton from '@/components/AppButton';
import { useAuth } from '@/context/AuthContext';
import { apiRequest } from '@/lib/api';
import { colors, fonts, radii, spacing } from '@/lib/theme';

export default function LoginScreen() {
  const router = useRouter();
  const { refreshAuth } = useAuth();
  const [identify, setIdentify] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    if (!identify.trim() || !password) {
      setError('Username/email and password are required');
      return;
    }
    setLoading(true);
    setError('');
    try {
      await apiRequest('/login', { method: 'POST', data: { identify: identify.trim(), password } });
      await refreshAuth();
      router.replace('/(tabs)/feed');
    } catch (err: any) {
      setError(err?.message || 'Login failed');
      setLoading(false);
    }
  };

  return (
    <Screen>
      <View style={styles.center}>
        <Card style={styles.card}>
          <Text style={styles.title}>Welcome Back</Text>
          <Text style={styles.subtitle}>Log in to continue</Text>

          {error ? (
            <View style={styles.errorBox}>
              <Text style={styles.errorText}>{error}</Text>
            </View>
          ) : null}

          <View style={styles.form}>
            <AppInput
              label="Username or Email"
              value={identify}
              onChangeText={setIdentify}
              placeholder="Username or Email"
              autoCapitalize="none"
            />
            <AppInput
              label="Password"
              value={password}
              onChangeText={setPassword}
              placeholder="Password"
              secureTextEntry
              autoCapitalize="none"
            />
            <AppButton title={loading ? 'Logging in...' : 'Log In'} onPress={handleSubmit} loading={loading} />
          </View>

          <Text style={styles.footer}>
            Don&apos;t have an account? <Link href="/register" style={styles.link}>Register here</Link>
          </Text>
        </Card>
      </View>
    </Screen>
  );
}

const styles = StyleSheet.create({
  center: {
    flex: 1,
    justifyContent: 'center',
    padding: spacing.lg,
  },
  card: {
    gap: spacing.md,
  },
  title: {
    color: colors.textPrimary,
    fontFamily: fonts.sansBold,
    fontSize: 26,
    letterSpacing: -0.02,
    textAlign: 'center',
  },
  subtitle: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    textAlign: 'center',
    marginBottom: spacing.sm,
  },
  errorBox: {
    backgroundColor: colors.dangerSoft,
    borderWidth: 1,
    borderColor: colors.dangerBorder,
    borderRadius: radii.md,
    padding: spacing.md,
  },
  errorText: {
    color: colors.danger,
    fontFamily: fonts.sans,
    fontSize: 14,
  },
  form: {
    gap: spacing.md,
  },
  footer: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    textAlign: 'center',
    fontSize: 14,
  },
  link: {
    color: colors.accent,
    fontFamily: fonts.sansSemiBold,
  },
});