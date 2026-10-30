import { ReactNode } from 'react';
import { ScrollView, StyleProp, StyleSheet, View, ViewStyle } from 'react-native';
import { colors, spacing } from '@/lib/theme';

interface Props {
  children: ReactNode;
  scroll?: boolean;
  style?: StyleProp<ViewStyle>;
  contentContainerStyle?: StyleProp<ViewStyle>;
  padded?: boolean;
  keyboardShouldPersistTaps?: 'always' | 'handled' | 'never';
}

// Page shell: dark background with an app-level top safe-area inset.
export default function Screen({
  children,
  scroll = true,
  style,
  contentContainerStyle,
  padded = true,
  keyboardShouldPersistTaps = 'never',
}: Props) {
  if (!scroll) {
    return <View style={[styles.base, style]}>{children}</View>;
  }

  return (
    <View style={[styles.base, style]}>
      <ScrollView
        keyboardShouldPersistTaps={keyboardShouldPersistTaps}
        contentContainerStyle={[
          styles.scrollContent,
          padded && styles.padded,
          contentContainerStyle,
        ]}
      >
        {children}
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  base: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  // flexGrow lets a short page (e.g. the centred login card) fill the viewport.
  scrollContent: {
    flexGrow: 1,
  },
  padded: {
    padding: spacing.lg,
  },
});