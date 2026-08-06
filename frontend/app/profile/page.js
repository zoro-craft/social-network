'use client';
import { useAuth } from '../context/AuthContext';
export default function ProfilePage() {
  const { auth, setAuth } = useAuth();
  return (
    <div
      style={{
        maxWidth: 600,
        margin: '60px auto',
        textAlign: 'center',
        color: '#f4f0fb',
      }}
    >
      <h1>Welcome to Profile {auth.first_name}</h1>
    </div>
  );
}
