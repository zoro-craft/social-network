import { useState } from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import AppInput from '@/components/AppInput';
import AppButton from '@/components/AppButton';
import HeaderBar from '@/components/HeaderBar';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useToast } from '@/context/ToastContext';
import { apiRequest } from '@/lib/api';
import { colors, fonts, spacing } from '@/lib/theme';

export default function CreateEventScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { auth, loading } = useRequireAuth();
  const showToast = useToast();

  const groupId = Number(id);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [eventDate, setEventDate] = useState('');
  const [eventTime, setEventTime] = useState('');
  const [optionA, setOptionA] = useState('');
  const [optionB, setOptionB] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async () => {
    if (!groupId) {
      showToast('Invalid group');
      return;
    }
    if (!title.trim() || !eventDate.trim() || !eventTime.trim()) {
      showToast('Title, date and time are required');
      return;
    }
    setSubmitting(true);
    try {
      await apiRequest('/events', {
        method: 'POST',
        data: {
          group_id: groupId,
          title: title.trim(),
          description: description.trim(),
          eventDate: eventDate.trim(),
          eventTime: eventTime.trim(),
          fitsOptions: optionA.trim(),
          secondOptions: optionB.trim(),
        },
      });
      showToast('Event created!', false);
      router.back();
    } catch (err: any) {
      showToast(err?.message || 'Could not create event');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading || !auth.authenticated) return null;

  return (
    <View style={styles.flex}>
      <HeaderBar title="Create Event" onBack={() => router.back()} />
      <ScrollView contentContainerStyle={styles.form} keyboardShouldPersistTaps="handled">
        <AppInput label="Title" value={title} onChangeText={setTitle} placeholder="Event title" />
        <AppInput
          label="Description"
          value={description}
          onChangeText={setDescription}
          placeholder="Event details"
          multiline
        />
        <Text style={styles.hint}>
          Use YYYY-MM-DD for the date and 24h HH:MM for the time, e.g. 2026-09-30 & 15:30.
        </Text>
        <View style={styles.row}>
          <View style={styles.col}>
            <AppInput label="Date" value={eventDate} onChangeText={setEventDate} placeholder="YYYY-MM-DD" />
          </View>
          <View style={styles.col}>
            <AppInput label="Time" value={eventTime} onChangeText={setEventTime} placeholder="HH:MM" />
          </View>
        </View>
        <AppInput
          label="Participation option 1"
          value={optionA}
          onChangeText={setOptionA}
          placeholder="e.g. Yes, available"
        />
        <AppInput
          label="Participation option 2"
          value={optionB}
          onChangeText={setOptionB}
          placeholder="e.g. No, not available"
        />
        <AppButton
          title={submitting ? 'Creating...' : 'Create Event'}
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
  row: {
    flexDirection: 'row',
    gap: spacing.md,
  },
  col: {
    flex: 1,
  },
});