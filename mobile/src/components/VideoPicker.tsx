import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { VideoView, useVideoPlayer } from 'expo-video';
import { PickedImage, pickVideo } from '@/hooks/useImagePicker';
import { colors, fonts, radii, spacing } from '@/lib/theme';

interface Props {
  label?: string;
  value: PickedImage | null;
  onChange: (asset: PickedImage | null) => void;
}

// Video counterpart to PhotoPicker: same label + preview + clear affordance.
export default function VideoPicker({ label = 'Video', value, onChange }: Props) {
  const choose = async () => {
    const result = await pickVideo();
    if (result.ok) onChange(result.asset);
  };

  return (
    <View style={styles.group}>
      <Text style={styles.label}>{label}</Text>

      {value ? (
        <View style={styles.previewRow}>
          <VideoPreview uri={value.uri} />
          <Pressable onPress={() => onChange(null)} style={styles.remove} hitSlop={8}>
            <Ionicons name="close-circle" size={26} color={colors.danger} />
          </Pressable>
        </View>
      ) : null}

      <Pressable onPress={choose} style={styles.pick}>
        <Ionicons name="videocam-outline" size={18} color={colors.textSecondary} />
        <Text style={styles.pickText}>{value ? 'Change video' : 'Choose video'}</Text>
      </Pressable>
    </View>
  );
}

// Mounted only while a video is selected, so the player's lifetime matches the pick.
function VideoPreview({ uri }: { uri: string }) {
  const player = useVideoPlayer({ uri }, (instance) => {
    instance.loop = true;
  });

  return (
    <VideoView
      player={player}
      style={styles.preview}
      contentFit="contain"
      nativeControls
    />
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
    flex: 1,
    height: 200,
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
