import { Redirect } from 'expo-router';
import { Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';

export default function Index() {
  const { auth, loading } = useAuth();

  if (loading) return <Loading />;

  return <Redirect href={auth.authenticated ? '/(tabs)/feed' : '/login'} />;
}