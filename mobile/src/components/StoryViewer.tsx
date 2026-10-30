import { useCallback, useEffect, useRef, useState } from 'react';
import { Animated, Modal, Pressable, StyleSheet, Text, View } from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { Image } from 'expo-image';
import { Ionicons } from '@expo/vector-icons';
import Avatar from '@/components/Avatar';
import { useAuth } from '@/context/AuthContext';
import { useToast } from '@/context/ToastContext';
import { apiRequest, mediaSource } from '@/lib/api';
import formatDate from '@/lib/format';
import { StoryGroup } from '@/lib/types';
import { colors, fonts, radii, spacing } from '@/lib/theme';

// How long a single story stays on screen before auto-advancing.
const STORY_DURATION = 5000;

interface Props {
  groups: StoryGroup[];
  /** Index into `groups` of the story to open, or null when closed. */
  startIndex: number | null;
  onClose: () => void;
  /** Fired after a story is deleted so the rail can refetch. */
  onChanged?: () => void;
}

export default function StoryViewer({ groups, startIndex, onClose, onChanged }: Props) {
  const { auth } = useAuth();
  const showToast = useToast();
  const insets = useSafeAreaInsets();

  // The rail remounts this component (via `key`) whenever a different group is
  // opened, so the starting position only ever has to be derived once.
  const [position, setPosition] = useState({ group: startIndex ?? 0, story: 0 });
  const progress = useRef(new Animated.Value(0)).current;
  const animation = useRef<Animated.CompositeAnimation | null>(null);

  const visible = startIndex !== null;
  const group = groups[position.group];
  const story = group?.stories[position.story];
  const storySource = mediaSource(story?.image);
  const isOwn = !!story && story.user_id === auth.user?.id;

  const next = useCallback(() => {
    const current = groups[position.group];
    if (!current) {
      onClose();
      return;
    }
    if (position.story + 1 < current.stories.length) {
      setPosition({ group: position.group, story: position.story + 1 });
      return;
    }
    if (position.group + 1 < groups.length) {
      setPosition({ group: position.group + 1, story: 0 });
      return;
    }
    onClose();
  }, [groups, position, onClose]);

  const previous = useCallback(() => {
    if (position.story > 0) {
      setPosition({ group: position.group, story: position.story - 1 });
      return;
    }
    if (position.group > 0) {
      const earlier = groups[position.group - 1];
      setPosition({
        group: position.group - 1,
        story: Math.max(0, (earlier?.stories.length || 1) - 1),
      });
    }
  }, [groups, position]);

  // The animation callback needs the latest `next` without restarting the timer.
  const nextRef = useRef(next);
  useEffect(() => {
    nextRef.current = next;
  }, [next]);

  const storyId = story?.id;

  useEffect(() => {
    if (startIndex === null || storyId === undefined) return;

    progress.setValue(0);
    const current = Animated.timing(progress, {
      toValue: 1,
      duration: STORY_DURATION,
      useNativeDriver: false,
    });
    animation.current = current;
    current.start(({ finished }) => {
      // A false `finished` means we stopped it deliberately (tap / close).
      if (finished) nextRef.current();
    });

    // Let the author know who saw their story; a failure must not block playback.
    apiRequest(`/stories/${storyId}/view`, { method: 'POST' }).catch(() => {});

    return () => {
      current.stop();
    };
  }, [startIndex, storyId, progress]);

  const remove = async () => {
    if (!story || !isOwn) return;

    const remaining = group.stories.length - 1;

    try {
      await apiRequest(`/stories/${story.id}`, { method: 'DELETE' });
    } catch {
      showToast('Could not delete story');
      return;
    }

    onChanged?.();

    if (remaining > 0) {
      // The parent refetch shifts the next story into this index.
      setPosition((prev) => ({
        group: prev.group,
        story: Math.min(prev.story, remaining - 1),
      }));
    } else {
      onClose();
    }
  };

  const close = () => {
    animation.current?.stop();
    onClose();
  };

  return (
    <Modal visible={visible} transparent animationType="fade" onRequestClose={close}>
      <View style={styles.backdrop}>
        {story ? (
          <>
            {storySource ? (
              <Image source={storySource} style={styles.media} contentFit="contain" />
            ) : null}

            <View
              style={[styles.top, { paddingTop: insets.top + spacing.md }]}
              pointerEvents="box-none"
            >
              <View style={styles.bars}>
                {group.stories.map((item, index) => (
                  <View key={item.id} style={styles.barTrack}>
                    <Animated.View
                      style={[
                        styles.barFill,
                        {
                          width:
                            index < position.story
                              ? '100%'
                              : index === position.story
                                ? progress.interpolate({
                                    inputRange: [0, 1],
                                    outputRange: ['0%', '100%'],
                                  })
                                : '0%',
                        },
                      ]}
                    />
                  </View>
                ))}
              </View>

              <View style={styles.header}>
                <Avatar src={group.avatar} name={group.full_name} size={34} />
                <View style={styles.headerText}>
                  <Text style={styles.name} numberOfLines={1}>
                    {group.full_name}
                  </Text>
                  <Text style={styles.time}>{formatDate(story.created_at)}</Text>
                </View>
                {isOwn ? (
                  <Pressable onPress={remove} hitSlop={8} style={styles.headerButton}>
                    <Ionicons name="trash-outline" size={22} color={colors.textPrimary} />
                  </Pressable>
                ) : null}
                <Pressable onPress={close} hitSlop={8} style={styles.headerButton}>
                  <Ionicons name="close" size={26} color={colors.textPrimary} />
                </Pressable>
              </View>
            </View>

            {/* Tap zones: the left third steps back, the rest advances. They start
                below the header so the close/delete buttons stay tappable. */}
            <View
              style={[styles.taps, { marginTop: insets.top + 100 }]}
              pointerEvents="box-none"
            >
              <Pressable style={styles.tapLeft} onPress={previous} />
              <Pressable style={styles.tapRight} onPress={() => nextRef.current()} />
            </View>

            {story.caption ? (
              <View style={[styles.caption, { bottom: insets.bottom + spacing.xxl }]}>
                <Text style={styles.captionText}>{story.caption}</Text>
              </View>
            ) : null}
          </>
        ) : null}
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  backdrop: {
    flex: 1,
    backgroundColor: '#000',
  },
  media: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
  },
  top: {
    paddingHorizontal: spacing.md,
  },
  bars: {
    flexDirection: 'row',
    gap: 4,
  },
  barTrack: {
    flex: 1,
    height: 3,
    borderRadius: 2,
    backgroundColor: 'rgba(255,255,255,0.3)',
    overflow: 'hidden',
  },
  barFill: {
    height: '100%',
    backgroundColor: colors.textPrimary,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    paddingVertical: spacing.md,
  },
  headerText: {
    flex: 1,
  },
  name: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 15,
  },
  time: {
    color: colors.textSecondary,
    fontFamily: fonts.mono,
    fontSize: 12,
  },
  headerButton: {
    padding: spacing.xs,
  },
  taps: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    flexDirection: 'row',
    marginBottom: 90,
  },
  tapLeft: {
    flex: 1,
  },
  tapRight: {
    flex: 2,
  },
  caption: {
    position: 'absolute',
    left: spacing.lg,
    right: spacing.lg,
    backgroundColor: 'rgba(0,0,0,0.55)',
    borderRadius: radii.md,
    padding: spacing.md,
  },
  captionText: {
    color: colors.textPrimary,
    fontFamily: fonts.sans,
    fontSize: 15,
  },
});
