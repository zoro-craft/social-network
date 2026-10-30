import { useState } from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { Link, useRouter } from 'expo-router';
import Screen from '@/components/Screen';
import Card from '@/components/Card';
import AppInput from '@/components/AppInput';
import AppButton from '@/components/AppButton';
import PhotoPicker from '@/components/PhotoPicker';
import { apiRequest } from '@/lib/api';
import { appendAsset, PickedImage } from '@/hooks/useImagePicker';
import { colors, fonts, radii, spacing } from '@/lib/theme';

const MAX_AVATAR_BYTES = 10 * 1024 * 1024;

const INITIAL = {
  first_name: '',
  last_name: '',
  email: '',
  password: '',
  date_of_birth: '',
  nickname: '',
  about_me: '',
};

export default function RegisterScreen() {
  const router = useRouter();
  const [form, setForm] = useState(INITIAL);
  const [avatar, setAvatar] = useState<PickedImage | null>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const set = (key: keyof typeof INITIAL) => (value: string) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const handleSubmit = async () => {
    setError('');
    if (avatar && avatar.fileSize && avatar.fileSize > MAX_AVATAR_BYTES) {
      setError('Avatar must be less than 10MB');
      return;
    }
    setLoading(true);

    const formData = new FormData();
    Object.entries(form).forEach(([key, value]) => formData.append(key, value));
    if (avatar) appendAsset(formData, 'avatar', avatar);

    try {
      await apiRequest('/register', { method: 'POST', formData });
      router.replace('/login');
    } catch (err: any) {
      setError(err?.message || 'Registration failed');
      setLoading(false);
    }
  };

  return (
    <Screen scroll>
      <Card style={styles.card}>
        <Text style={styles.title}>Create Account</Text>
        <Text style={styles.subtitle}>Join us today</Text>

        {error ? (
          <View style={styles.errorBox}>
            <Text style={styles.errorText}>{error}</Text>
          </View>
        ) : null}

        <View style={styles.form}>
          <View style={styles.row}>
            <View style={styles.col}>
              <AppInput label="First Name *" value={form.first_name} onChangeText={set('first_name')} maxLength={30} />
            </View>
            <View style={styles.col}>
              <AppInput label="Last Name *" value={form.last_name} onChangeText={set('last_name')} maxLength={30} />
            </View>
          </View>

          <View style={styles.row}>
            <View style={styles.col}>
              <AppInput
                label="Email Address *"
                value={form.email}
                onChangeText={set('email')}
                keyboardType="email-address"
                autoCapitalize="none"
                maxLength={100}
              />
            </View>
            <View style={styles.col}>
              <AppInput
                label="Date of Birth *"
                value={form.date_of_birth}
                onChangeText={set('date_of_birth')}
                placeholder="YYYY-MM-DD"
              />
            </View>
          </View>

          <View style={styles.row}>
            <View style={styles.col}>
              <AppInput
                label="Password *"
                value={form.password}
                onChangeText={set('password')}
                secureTextEntry
                autoCapitalize="none"
                maxLength={50}
              />
            </View>
            <View style={styles.col}>
              <AppInput label="Nickname (optional)" value={form.nickname} onChangeText={set('nickname')} maxLength={30} />
            </View>
          </View>

          <AppInput
            label="About Me (optional)"
            value={form.about_me}
            onChangeText={set('about_me')}
            multiline
            maxLength={200}
          />

          <PhotoPicker label="Avatar / Image" value={avatar} onChange={setAvatar} />

          <AppButton title={loading ? 'Creating Account...' : 'Register'} onPress={handleSubmit} loading={loading} />
        </View>

        <Text style={styles.footer}>
          Already have an account? <Link href="/login" style={styles.link}>Log in here</Link>
        </Text>
      </Card>
    </Screen>
  );
}

const styles = StyleSheet.create({
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
  row: {
    flexDirection: 'row',
    gap: spacing.md,
  },
  col: {
    flex: 1,
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