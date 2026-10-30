import {
  ActivityIndicator,
  Pressable,
  StyleSheet,
  Text,
  ViewStyle,
} from 'react-native';
import { colors, fonts, radii, spacing } from '@/lib/theme';

export type ButtonVariant = 'primary' | 'ghost' | 'danger' | 'pending' | 'success';

interface Props {
  title: string;
  onPress?: () => void;
  variant?: ButtonVariant;
  loading?: boolean;
  disabled?: boolean;
  style?: ViewStyle;
}

// primary = solid violet (mirrors .glass-button), ghost = bordered surface.
export default function AppButton({
  title,
  onPress,
  variant = 'primary',
  loading = false,
  disabled = false,
  style,
}: Props) {
  const palette = {
    primary: {
      backgroundColor: colors.accent,
      color: colors.onAccent,
      borderColor: colors.accent,
    },
    ghost: {
      backgroundColor: colors.bgElevated,
      color: colors.textPrimary,
      borderColor: colors.border,
    },
    danger: {
      backgroundColor: colors.dangerSoft,
      color: colors.danger,
      borderColor: colors.dangerBorder,
    },
    pending: {
      backgroundColor: colors.warningSoft,
      color: colors.warning,
      borderColor: colors.warningBorder,
    },
    success: {
      backgroundColor: colors.successSoft,
      color: colors.success,
      borderColor: colors.successBorder,
    },
  }[variant];

  const isBusy = loading || disabled;

  return (
    <Pressable
      onPress={onPress}
      disabled={isBusy}
      style={({ pressed }) => [
        styles.button,
        palette,
        pressed && !isBusy && styles.pressed,
        isBusy && styles.disabled,
        style,
      ]}
    >
      {loading ? (
        <ActivityIndicator size="small" color={palette.color} />
      ) : (
        <Text style={[styles.label, { color: palette.color }]}>{title}</Text>
      )}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  button: {
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
    borderRadius: radii.sm,
    borderWidth: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  pressed: {
    opacity: 0.92,
    transform: [{ scale: 0.98 }],
  },
  disabled: {
    opacity: 0.5,
  },
  label: {
    fontSize: 15,
    fontFamily: fonts.sansSemiBold,
  },
});