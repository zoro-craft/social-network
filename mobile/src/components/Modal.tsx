import { ReactNode } from 'react';
import { Modal as RNModal, Pressable, ScrollView, StyleSheet, Text } from 'react-native';
import { colors, fonts, radii, spacing } from '@/lib/theme';

interface Props {
  visible: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  /**
   * Wrap the body in a ScrollView. Use this for long forms; leave it off when
   * the body is already a FlatList, since nesting scrollables breaks both.
   */
  scroll?: boolean;
}

export default function Modal({ visible, onClose, title, children, scroll = false }: Props) {
  return (
    <RNModal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={onClose}
    >
      <Pressable style={styles.overlay} onPress={onClose}>
        <Pressable style={styles.content} onPress={(e) => e.stopPropagation()}>
          {title ? <Text style={styles.title}>{title}</Text> : null}
          {scroll ? (
            <ScrollView
              style={styles.scroll}
              contentContainerStyle={styles.scrollContent}
              keyboardShouldPersistTaps="handled"
            >
              {children}
            </ScrollView>
          ) : (
            children
          )}
        </Pressable>
      </Pressable>
    </RNModal>
  );
}

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: colors.scrim,
    alignItems: 'center',
    justifyContent: 'center',
    padding: spacing.xl,
  },
  content: {
    width: '100%',
    maxWidth: 420,
    // Without this a tall body would push the action button off screen.
    maxHeight: '90%',
    backgroundColor: colors.bgCard,
    borderWidth: 1,
    borderColor: colors.borderStrong,
    borderRadius: radii.lg,
    padding: spacing.xl,
  },
  title: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 20,
    letterSpacing: -0.01,
    marginTop: 0,
    textTransform: 'capitalize',
    marginBottom: spacing.lg,
  },
  scroll: {
    flexGrow: 0,
  },
  scrollContent: {
    paddingBottom: spacing.sm,
  },
});
