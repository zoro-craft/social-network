import { useEffect, useState } from 'react';
import { StyleSheet, Text, View, Switch } from 'react-native';
import Modal from '@/components/Modal';
import AppInput from '@/components/AppInput';
import AppButton from '@/components/AppButton';
import PhotoPicker from '@/components/PhotoPicker';
import { apiRequest } from '@/lib/api';
import { appendAsset, PickedImage } from '@/hooks/useImagePicker';
import { User } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

interface Props {
  user: User;
  onClose: () => void;
  onSaved: () => void;
  onShowError: (message: string, isError?: boolean) => void;
}

const INITIAL = {
  first_name: '',
  last_name: '',
  email: '',
  date_of_birth: '',
  nickname: '',
  about_me: '',
};

export default function EditProfileModal({ user, onClose, onSaved, onShowError }: Props) {
  const [form, setForm] = useState(INITIAL);
  const [isPrivate, setIsPrivate] = useState(!!user.is_private);
  const [avatar, setAvatar] = useState<PickedImage | null>(null);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    setForm({
      first_name: user.first_name || '',
      last_name: user.last_name || '',
      email: user.email || '',
      date_of_birth: user.date_of_birth ? user.date_of_birth.split('T')[0] : '',
      nickname: user.nickname || '',
      about_me: user.about_me || '',
    });
  }, [user]);

  const set = (key: keyof typeof INITIAL) => (value: string) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const handleSubmit = async () => {
    setSaving(true);
    try {
      const formData = new FormData();
      Object.entries(form).forEach(([key, value]) => formData.append(key, value));
      formData.append('is_private', isPrivate ? 'true' : 'false');
      if (avatar) appendAsset(formData, 'avatar', avatar);
      await apiRequest('/profile/edit', { method: 'POST', formData });
      onSaved();
    } catch (error: any) {
      onShowError(error?.message || 'Failed to update profile');
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal visible onClose={onClose} title="Edit Profile" scroll>
      <View style={styles.form}>
        <AppInput label="First Name" value={form.first_name} onChangeText={set('first_name')} placeholder="First Name" />
        <AppInput label="Last Name" value={form.last_name} onChangeText={set('last_name')} placeholder="Last Name" />
        <AppInput
          label="Email"
          value={form.email}
          onChangeText={set('email')}
          placeholder="Email"
          keyboardType="email-address"
          autoCapitalize="none"
        />
        <AppInput
          label="Date of Birth"
          value={form.date_of_birth}
          onChangeText={set('date_of_birth')}
          placeholder="YYYY-MM-DD"
        />
        <AppInput label="Nickname" value={form.nickname} onChangeText={set('nickname')} placeholder="Nickname" maxLength={30} />
        <AppInput
          label="About me"
          value={form.about_me}
          onChangeText={set('about_me')}
          placeholder="About me"
          multiline
        />
        <View style={styles.toggleRow}>
          <Text style={styles.toggleLabel}>Private Profile</Text>
          <Switch
            value={isPrivate}
            onValueChange={setIsPrivate}
            trackColor={{ true: colors.accent, false: colors.bgHover }}
            thumbColor={colors.textSecondary}
          />
        </View>
        <PhotoPicker label="Avatar" value={avatar} onChange={setAvatar} />
        <AppButton title={saving ? 'Saving...' : 'Save'} onPress={handleSubmit} loading={saving} />
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  form: {
    gap: spacing.md,
  },
  toggleRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: spacing.xs,
  },
  toggleLabel: {
    color: colors.textPrimary,
    fontFamily: fonts.sans,
    fontSize: 15,
  },
});