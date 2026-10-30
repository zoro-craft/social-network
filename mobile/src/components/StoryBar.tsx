import { useCallback, useState } from 'react';
import { Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useFocusEffect, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import Avatar from '@/components/Avatar';
import StoryViewer from '@/components/StoryViewer';
import { useAuth } from '@/context/AuthContext';
import { apiRequest } from '@/lib/api';
import { StoryGroup } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

// The horizontal "stories" rail that sits above the feed, Instagram-style:
// your own bubble first, then everyone you follow who has an active story.
export default function StoryBar() {
  const router = useRouter();
  const { auth } = useAuth();

  const [groups, setGroups] = useState<StoryGroup[]>([]);
  const [activeIndex, setActiveIndex] = useState<number | null>(null);

  const load = useCallback(async () => {
    if (!auth.authenticated) return;
    try {
      const data = await apiRequest<StoryGroup[]>('/stories');
      setGroups(data || []);
    } catch {
      setGroups([]);
    }
  }, [auth.authenticated]);

  // Refetch on focus so a story created on the create screen shows up on return.
  useFocusEffect(
    useCallback(() => {
      load();
    }, [load]),
  );

  const ownId = auth.user?.id;
  const ownGroup = groups.find((g) => g.user_id === ownId);
  const others = groups.filter((g) => g.user_id !== ownId);

  const openOwn = () => {
    if (ownGroup) {
      setActiveIndex(0);
    } else {
      router.push('/stories/create');
    }
  };

  return (
    <View style={styles.wrap}>
      <ScrollView
        horizontal
        showsHorizontalScrollIndicator={false}
        contentContainerStyle={styles.row}
      >
        <View style={styles.tile}>
          <Pressable onPress={openOwn} style={styles.avatarButton}>
            <View style={[styles.ring, ownGroup && !ownGroup.seen && styles.ringUnseen]}>
              <Avatar src={auth.user?.avatar} name={auth.user?.first_name} size={56} />
            </View>
          </Pressable>
          {/* Sibling of the avatar button, so the two taps never overlap. */}
          <Pressable
            onPress={() => router.push('/stories/create')}
            style={styles.addBadge}
            hitSlop={6}
          >
            <Ionicons name="add" size={15} color={colors.onAccent} />
          </Pressable>
          <Text style={styles.label} numberOfLines={1}>
            Your story
          </Text>
        </View>

        {others.map((group, index) => (
          <Pressable
            key={group.user_id}
            style={styles.tile}
            onPress={() => setActiveIndex(ownGroup ? index + 1 : index)}
          >
            <View style={[styles.ring, group.seen ? styles.ringSeen : styles.ringUnseen]}>
              <Avatar src={group.avatar} name={group.full_name} size={56} />
            </View>
            <Text style={styles.label} numberOfLines={1}>
              {group.full_name.split(' ')[0]}
            </Text>
          </Pressable>
        ))}
      </ScrollView>

      <StoryViewer
        key={activeIndex ?? 'closed'}
        groups={ownGroup ? [ownGroup, ...others] : others}
        startIndex={activeIndex}
        onClose={() => setActiveIndex(null)}
        onChanged={load}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  wrap: {
    paddingVertical: spacing.sm,
  },
  // No horizontal padding: the rail inherits the feed's gutter so the bubbles
  // line up with the post cards below them.
  row: {
    gap: spacing.md,
  },
  tile: {
    width: 68,
    alignItems: 'center',
    gap: spacing.xs,
  },
  avatarButton: {
    padding: 2,
  },
  // Ring wraps the avatar; unseen stories get a bright border, seen ones fade out.
  ring: {
    padding: 2,
    borderRadius: 999,
    borderWidth: 2,
    borderColor: colors.border,
  },
  ringUnseen: {
    borderColor: colors.accent,
  },
  ringSeen: {
    borderColor: colors.borderStrong,
  },
  addBadge: {
    position: 'absolute',
    right: 2,
    top: 40,
    width: 22,
    height: 22,
    borderRadius: 11,
    backgroundColor: colors.accent,
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 2,
    borderColor: colors.bg,
  },
  label: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    fontSize: 11,
  },
});
