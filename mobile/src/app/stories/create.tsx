import { useState } from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import AppButton from '@/components/AppButton';
import AppInput from '@/components/AppInput';
import HeaderBar from '@/components/HeaderBar';
import PhotoPicker from '@/components/PhotoPicker';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useToast } from '@/context/ToastContext';
import { apiRequest } from '@/lib/api';
import { appendAsset, PickedImage } from '@/hooks/useImagePicker';
import { colors, fonts, spacing } from '@/lib/theme';

export default function CreateStoryScreen() {
  const router = useRouter();
  const { auth, loading } = useRequireAuth();
  const showToast = useToast();

  const [image, setImage] = useState<PickedImage | null>(null);
  const [caption, setCaption] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async () => {
    if (!image) {
      showToast('Pick an image for your story');
      return;
    }

    setSubmitting(true);
    try {
      const form = new FormData();
      form.append('caption', caption.trim());
      appendAsset(form, 'image', image);

      await apiRequest('/stories/create', { method: 'POST', formData: form });
      showToast('Story added!', false);
      router.back();
    } catch (err: any) {
      if (err?.status === 401) {
        router.replace('/login');
        return;
      }
      showToast(err?.message || 'Could not create story');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading || !auth.authenticated) return null;

  return (
    <View style={styles.flex}>
      <HeaderBar title="New Story" onBack={() => router.back()} />
      <ScrollView contentContainerStyle={styles.form} keyboardShouldPersistTaps="handled">
        <Text style={styles.hint}>
          Your story disappears after 24 hours and is only visible to your followers.
        </Text>

        <PhotoPicker label="Story image" value={image} onChange={setImage} />

        <AppInput
          label="Caption (optional)"
          value={caption}
          onChangeText={setCaption}
          placeholder="Add a caption..."
          multiline
          maxLength={200}
        />

        <AppButton
          title={submitting ? 'Sharing...' : 'Share to story'}
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
