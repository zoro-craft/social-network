import { useEffect, useState } from 'react';
import { Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import AppInput from '@/components/AppInput';
import AppButton from '@/components/AppButton';
import Avatar from '@/components/Avatar';
import HeaderBar from '@/components/HeaderBar';
import OptionPicker, { PickOption } from '@/components/OptionPicker';
import PhotoPicker from '@/components/PhotoPicker';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useToast } from '@/context/ToastContext';
import { apiRequest } from '@/lib/api';
import { appendAsset, PickedImage } from '@/hooks/useImagePicker';
import { Group, UserSearchResult } from '@/lib/types';
import { colors, fonts, radii, spacing } from '@/lib/theme';

const PRIVACY_OPTIONS: PickOption[] = [
  { value: 'public', label: 'Public' },
  { value: 'almost_private', label: 'Almost private (followers)' },
  { value: 'private', label: 'Private (chosen followers)' },
];

export default function CreatePostScreen() {
  const router = useRouter();
  const params = useLocalSearchParams<{ group_id?: string }>();
  const preselectedGroupId = params.group_id || '';
  const { auth, loading } = useRequireAuth();
  const showToast = useToast();

  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [privacy, setPrivacy] = useState('public');
  const [image, setImage] = useState<PickedImage | null>(null);
  const [groups, setGroups] = useState<Group[]>([]);
  const [selectedGroupId, setSelectedGroupId] = useState(preselectedGroupId);
  const [followers, setFollowers] = useState<UserSearchResult[]>([]);
  const [selectedViewers, setSelectedViewers] = useState<number[]>([]);
  const [submitting, setSubmitting] = useState(false);

  // The web uses /groups/my-groups (which doesn't exist); GET /groups + is_member works.
  useEffect(() => {
    if (!auth.authenticated) return;
    apiRequest<Group[]>('/groups')
      .then((data) => setGroups((data || []).filter((g) => g.is_member)))
      .catch(() => setGroups([]));
  }, [auth.authenticated]);

  useEffect(() => {
    if (privacy !== 'private' || !auth.authenticated || selectedGroupId) {
      setFollowers([]);
      return;
    }
    apiRequest<UserSearchResult[]>(`/profile/${auth.user?.id}/followers`)
      .then((data) => setFollowers(data || []))
      .catch(() => setFollowers([]));
  }, [privacy, auth.authenticated, auth.user?.id, selectedGroupId]);

  const toggleViewer = (id: number) => {
    setSelectedViewers((prev) =>
      prev.includes(id) ? prev.filter((v) => v !== id) : [...prev, id],
    );
  };

  const handleGroupChange = (value: string) => {
    setSelectedGroupId(value);
    if (value) {
      setPrivacy('public');
      setSelectedViewers([]);
    }
  };

  const isGroupPost = !!selectedGroupId;

  const handleSubmit = async () => {
    if (!title.trim() || !content.trim()) {
      showToast('Title and content are required');
      return;
    }
    if (!selectedGroupId && privacy === 'private' && selectedViewers.length === 0) {
      showToast('Select at least one follower');
      return;
    }

    setSubmitting(true);
    try {
      const form = new FormData();
      form.append('title', title.trim());
      form.append('content', content.trim());
      form.append('privacy', selectedGroupId ? 'public' : privacy);
      if (selectedGroupId) form.append('group_id', selectedGroupId);
      if (image) appendAsset(form, 'image', image);
      selectedViewers.forEach((id) => form.append('viewer_ids', String(id)));

      await apiRequest('/posts/create', { method: 'POST', formData: form });
      showToast('Post created successfully!', false);
      if (selectedGroupId) {
        router.replace(`/groups/${selectedGroupId}`);
      } else {
        router.back();
      }
    } catch (err: any) {
      if (err?.status === 401) {
        router.replace('/login');
        return;
      }
      showToast(err?.message || 'Could not create post');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading || !auth.authenticated) return null;

  return (
    <View style={styles.flex}>
      <HeaderBar title="Create Post" onBack={() => router.back()} />
      {/* Scrolls so the Post button stays reachable once the keyboard is up. */}
      <ScrollView contentContainerStyle={styles.form} keyboardShouldPersistTaps="handled">


        <AppInput label="Title" value={title} onChangeText={setTitle} placeholder="Title" />
        <AppInput
          label="Content"
          value={content}
          onChangeText={setContent}
          placeholder="Write something..."
          multiline
          maxLength={5000}
        />

        {!isGroupPost && (
          <OptionPicker
            label="Privacy"
            value={privacy}
            onSelect={setPrivacy}
            options={PRIVACY_OPTIONS}
          />
        )}
        {isGroupPost && (
          <Text style={styles.hint}>Group posts are always public</Text>
        )}

        {privacy === 'private' && !isGroupPost && (
          <View style={styles.viewers}>
            <Text style={styles.viewerLabel}>Choose who can see this post:</Text>
            {followers.length === 0 && <Text style={styles.muted}>No followers yet</Text>}
            <View style={styles.chipRow}>
              {followers.map((f) => {
                const selected = selectedViewers.includes(f.id);
                return (
                  <Pressable
                    key={f.id}
                    onPress={() => toggleViewer(f.id)}
                    style={[styles.chip, selected && styles.chipActive]}
                  >
                    <Avatar src={f.avatar} name={f.first_name} size={28} />
                    <Text style={[styles.chipName, selected && styles.chipNameActive]}>
                      {f.first_name}
                    </Text>
                  </Pressable>
                );
              })}
            </View>
          </View>
        )}

        <PhotoPicker label="Image" value={image} onChange={setImage} />

        <AppButton
          title={submitting ? 'Posting...' : 'Post'}
          onPress={handleSubmit}
          loading={submitting}
        />
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  flex: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  form: {
    padding: spacing.lg,
    paddingBottom: spacing.xl,
    gap: spacing.md,
  },
  hint: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    fontSize: 13,
  },
  viewers: {
    gap: spacing.sm,
  },
  viewerLabel: {
    color: colors.textMuted,
    fontFamily: fonts.sansSemiBold,
    fontSize: 11,
    letterSpacing: 0.08,
    textTransform: 'uppercase',
  },
  muted: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    fontSize: 14,
  },
  chipRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.sm,
  },
  chip: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
    paddingVertical: spacing.xs,
    paddingHorizontal: spacing.sm,
    borderRadius: radii.pill,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.bgElevated,
  },
  chipActive: {
    backgroundColor: colors.accentSoft,
    borderColor: colors.accentBorder,
  },
  chipName: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    fontSize: 13,
  },
  chipNameActive: {
    color: colors.accent,
  },
});