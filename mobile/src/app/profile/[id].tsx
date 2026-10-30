import { StyleSheet, View } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import HeaderBar from '@/components/HeaderBar';
import ProfileScreen from '@/components/ProfileScreen';
import { Loading } from '@/components/StateViews';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { colors } from '@/lib/theme';

export default function UserProfileScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { auth, loading } = useRequireAuth();

  if (loading || !auth.authenticated || !id) return <Loading />;

  return (
    <View style={styles.screen}>
      <HeaderBar title="Profile" onBack={() => router.back()} />
      <ProfileScreen userId={Number(id)} />
    </View>
  );
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: colors.bg,
  },
});