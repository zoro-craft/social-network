import { StyleSheet, View } from 'react-native';
import ProfileScreen from '@/components/ProfileScreen';
import { Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { colors } from '@/lib/theme';

export default function ProfileTab() {
  const { auth, loading } = useAuth();

  if (loading || !auth.user) return <Loading />;

  return (
    <View style={styles.screen}>
      <ProfileScreen userId={auth.user.id} />
    </View>
  );
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: colors.bg,
  },
});