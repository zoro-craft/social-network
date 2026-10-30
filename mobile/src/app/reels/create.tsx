import { useState } from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import AppButton from '@/components/AppButton';
import AppInput from '@/components/AppInput';
import HeaderBar from '@/components/HeaderBar';
import VideoPicker from '@/components/VideoPicker';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useToast } from '@/context/ToastContext';
import { apiRequest } from '@/lib/api';
import { appendAsset, PickedImage } from '@/hooks/useImagePicker';
import { colors, fonts, spacing } from '@/lib/theme';

export default function CreateReelScreen() {
  const router = useRouter();
  const { auth, loading } = useRequireAuth();
  const showToast = useToast();

  const [video, setVideo] = useState<PickedImage | null>(null);
  const [caption, setCaption] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async () => {
    if (!video) {
      showToast('Pick a video for your reel');
      return;
    }

    setSubmitting(true);
    try {
      const form = new FormData();
      form.append('caption', caption.trim());
      // Videos do not always carry a mimeType, so fall back to mp4 rather than the
      // image default; the backend validates the real container from the filename.
      appendAsset(form, 'video', video, 'video/mp4');

      await apiRequest('/reels/create', { method: 'POST', formData: form });
      showToast('Reel posted!', false);
      router.back();
    } catch (err: any) {
      if (err?.status === 401) {
        router.replace('/login');
        return;
      }
      showToast(err?.message || 'Could not create reel');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading || !auth.authenticated) return null;

  return (
    <View style={styles.flex}>
      <HeaderBar title="New Reel" onBack={() => router.back()} />
      <ScrollView contentContainerStyle={styles.form} keyboardShouldPersistTaps="handled">
        <Text style={styles.hint}>
          Share a short video. MP4, MOV or WebM up to 100MB — reels stay on your profile.
        </Text>

        <VideoPicker label="Reel video" value={video} onChange={setVideo} />

        <AppInput
          label="Caption (optional)"
          value={caption}
          onChangeText={setCaption}
          placeholder="Write a caption..."
          multiline
          maxLength={300}
        />

        <AppButton
          title={submitting ? 'Posting...' : 'Share reel'}
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
});
