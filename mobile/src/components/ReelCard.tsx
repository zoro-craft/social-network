import { useEffect } from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useRouter } from 'expo-router';
import { VideoView, useVideoPlayer } from 'expo-video';
import Avatar from '@/components/Avatar';
import { mediaSource } from '@/lib/api';
import formatDate from '@/lib/format';
import { Reel } from '@/lib/types';
import { colors, fonts, radii, spacing } from '@/lib/theme';

interface Props {
  reel: Reel;
  /** Exact visible height, so each slide fills the paging viewport. */
  height: number;
  /** Whether this reel is the one currently on screen. */
  active: boolean;
  muted: boolean;
  isOwn: boolean;
  onToggleMute: () => void;
  onLike: (reel: Reel) => void;
  onSave: (reel: Reel) => void;
  onShare: (reel: Reel) => void;
  onComments: (reel: Reel) => void;
  onDelete: (reel: Reel) => void;
}

// One full-screen reel. Playback follows `active` so scrolling away pauses a video
// instead of leaving it running behind the next one.
export default function ReelCard({
  reel,
  height,
  active,
  muted,
  isOwn,
  onToggleMute,
  onLike,
  onSave,
  onShare,
  onComments,
  onDelete,
}: Props) {
  const router = useRouter();
  const source = mediaSource(reel.video);
  const player = useVideoPlayer(source, (instance) => {
    instance.loop = true;
    instance.muted = true;
  });

  useEffect(() => {
    if (active) {
      player.play();
    } else {
      player.pause();
    }
  }, [active, player]);

  useEffect(() => {
    player.muted = muted;
  }, [muted, player]);

  return (
    <View style={[styles.card, { height }]}>
      <VideoView
        player={player}
        style={styles.video}
        contentFit="contain"
        nativeControls={false}
      />

      <Pressable onPress={onToggleMute} style={styles.mute} hitSlop={8}>
        <Ionicons
          name={muted ? 'volume-mute' : 'volume-high'}
          size={18}
          color={colors.textPrimary}
        />
      </Pressable>

      <View style={styles.bottom} pointerEvents="box-none">
        <View style={styles.info}>
          <View style={styles.authorRow}>
            <Pressable
              style={styles.author}
              onPress={() => router.push(`/profile/${reel.user_id}`)}
            >
              <Avatar src={reel.avatar} name={reel.full_name} size={32} />
              <Text style={styles.authorName} numberOfLines={1}>
                {reel.full_name}
              </Text>
            </Pressable>

            {isOwn ? (
              <Pressable onPress={() => onDelete(reel)} hitSlop={8} style={styles.delete}>
                <Ionicons name="trash-outline" size={20} color={colors.textPrimary} />
              </Pressable>
            ) : null}
          </View>

          {reel.caption ? (
            <Text style={styles.caption} numberOfLines={3}>
              {reel.caption}
            </Text>
          ) : null}

          <Text style={styles.meta}>
            {formatDate(reel.created_at)} · {reel.view_count} views
          </Text>
        </View>

        <View style={styles.actions}>
          <ActionButton
            icon={reel.liked ? 'heart' : 'heart-outline'}
            count={reel.like_count}
            active={reel.liked}
            onPress={() => onLike(reel)}
          />
          <ActionButton
            icon="chatbubble-outline"
            count={reel.comment_count}
            onPress={() => onComments(reel)}
          />
          <ActionButton
            icon="paper-plane-outline"
            count={reel.share_count}
            onPress={() => onShare(reel)}
          />
          <ActionButton
            icon={reel.saved ? 'bookmark' : 'bookmark-outline'}
            count={reel.save_count}
            active={reel.saved}
            onPress={() => onSave(reel)}
          />
        </View>
      </View>
    </View>
  );
}

function ActionButton({
  icon,
  count,
  active,
  onPress,
}: {
  icon: keyof typeof Ionicons.glyphMap;
  count: number;
  active?: boolean;
  onPress: () => void;
}) {
  return (
    <Pressable onPress={onPress} style={[styles.action, active && styles.actionActive]}>
      <Ionicons name={icon} size={22} color={active ? colors.danger : colors.textPrimary} />
      <Text style={styles.actionCount}>{count}</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    width: '100%',
    backgroundColor: '#000',
  },
  // Spelled out rather than Spread from StyleSheet.absoluteFillObject, which this
  // React Native version no longer exposes in its types.
  video: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
  },
  mute: {
    position: 'absolute',
    top: spacing.md,
    right: spacing.md,
    width: 36,
    height: 36,
    borderRadius: 18,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: 'rgba(0,0,0,0.45)',
  },
  // Scrim keeps the caption and counters legible over bright footage.
  bottom: {
    position: 'absolute',
    left: 0,
    right: 0,
    bottom: 0,
    flexDirection: 'row',
    alignItems: 'flex-end',
    justifyContent: 'space-between',
    gap: spacing.md,
    padding: spacing.lg,
    backgroundColor: 'rgba(0,0,0,0.35)',
  },
  info: {
    flex: 1,
    gap: spacing.xs,
  },
  authorRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  author: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    flexShrink: 1,
  },
  authorName: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 15,
    flexShrink: 1,
  },
  delete: {
    padding: spacing.xs,
  },
  caption: {
    color: colors.textPrimary,
    fontFamily: fonts.sans,
    fontSize: 14,
    lineHeight: 20,
  },
  meta: {
    color: colors.textSecondary,
    fontFamily: fonts.mono,
    fontSize: 12,
  },
  actions: {
    gap: spacing.sm,
    alignItems: 'center',
  },
  action: {
    width: 48,
    paddingVertical: spacing.sm,
    borderRadius: radii.md,
    alignItems: 'center',
    gap: 2,
    backgroundColor: 'rgba(0,0,0,0.45)',
  },
  actionActive: {
    backgroundColor: colors.dangerSoft,
  },
  actionCount: {
    color: colors.textPrimary,
    fontFamily: fonts.mono,
    fontSize: 11,
  },
});
