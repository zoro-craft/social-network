import { useState } from 'react';
import { ScrollView, StyleSheet, View } from 'react-native';
import { useRouter } from 'expo-router';
import AppInput from '@/components/AppInput';
import AppButton from '@/components/AppButton';
import HeaderBar from '@/components/HeaderBar';
import PhotoPicker from '@/components/PhotoPicker';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { appendAsset, PickedImage } from '@/hooks/useImagePicker';
import { useToast } from '@/context/ToastContext';
import { apiRequest } from '@/lib/api';
import { colors, spacing } from '@/lib/theme';

export default function CreateGroupScreen() {
  const router = useRouter();
  const { auth, loading } = useRequireAuth();
  const showToast = useToast();

  const [titel, setTitel] = useState('');
  const [description, setDescription] = useState('');
  const [image, setImage] = useState<PickedImage | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async () => {
    if (!titel.trim() || !description.trim()) {
      showToast('Title and description are required');
      return;
    }

    setSubmitting(true);
    try {
      const form = new FormData();
      form.append('titel', titel.trim());
      form.append('description', description.trim());
      if (image) appendAsset(form, 'image', image);

      await apiRequest('/groups', { method: 'POST', formData: form });
      showToast('Group created!', false);
      router.back();
    } catch (err: any) {
      showToast(err?.message || 'Could not create group');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading || !auth.authenticated) return null;

  return (
    <View style={styles.flex}>
      <HeaderBar title="Create Group" onBack={() => router.back()} />
      <ScrollView contentContainerStyle={styles.form} keyboardShouldPersistTaps="handled">
        <AppInput
          label="Title"
          value={titel}
          onChangeText={setTitel}
          placeholder="Group title"
          maxLength={50}
        />
        <AppInput
          label="Description"
          value={description}
          onChangeText={setDescription}
          placeholder="Tell people what this group is about"
          multiline
          maxLength={200}
        />
        <PhotoPicker label="Group image" value={image} onChange={setImage} />
        <AppButton
          title={submitting ? 'Creating...' : 'Create Group'}
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
});
