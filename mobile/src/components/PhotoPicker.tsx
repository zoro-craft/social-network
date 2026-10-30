import { Image } from 'expo-image';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { PickedImage, pickImage } from '@/hooks/useImagePicker';
import { colors, fonts, radii, spacing } from '@/lib/theme';

interface Props {
  label?: string;
  value: PickedImage | null;
  onChange: (asset: PickedImage | null) => void;
}

// Mirrors the web <label>Image + <input type=file>.</label>
export default function PhotoPicker({ label = 'Image', value, onChange }: Props) {
  const choose = async () => {
    const result = await pickImage();
    if (result.ok) onChange(result.asset);
  };

  return (
    <View style={styles.group}>
      <Text style={styles.label}>{label}</Text>
      {value ? (
        <View style={styles.previewRow}>
          <Image source={{ uri: value.uri }} style={styles.preview} contentFit="cover" />
          <Pressable onPress={() => onChange(null)} style={styles.remove} hitSlop={8}>
            <Ionicons name="close-circle" size={26} color={colors.danger} />
          </Pressable>
        </View>
      ) : null}
      <Pressable onPress={choose} style={styles.pick}>
        <Ionicons name="image-outline" size={18} color={colors.textSecondary} />
        <Text style={styles.pickText}>{value ? 'Change image' : 'Choose image'}</Text>
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  group: {
    gap: spacing.sm,
  },
  label: {
    color: colors.textMuted,
    fontFamily: fonts.sansSemiBold,
    fontSize: 11,
    letterSpacing: 0.08,
    textTransform: 'uppercase',
  },
  previewRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
  },
  preview: {
    width: 88,
    height: 88,
    borderRadius: radii.md,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.bgElevated,
  },
  remove: {
    padding: spacing.xs,
  },
  pick: {
    borderWidth: 1,
    borderColor: colors.borderStrong,
    borderRadius: radii.md,
    paddingVertical: spacing.md,
    alignItems: 'center',
    justifyContent: 'center',
    gap: 6,
    flexDirection: 'row',
    backgroundColor: colors.bgElevated,
  },
  pickText: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    fontSize: 14,
  },
});