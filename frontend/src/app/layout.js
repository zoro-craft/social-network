import './globals.css';
import { AuthProvider } from '@/app/context/AuthContext';
import { ToastProvider } from '@/app/context/ToastContext';
import { RealtimeProvider } from '@/app/components/RealtimeProvider';
import Navbar from '@/app/components/Navbar';
import AppShell from '@/app/components/AppShell';

export const metadata = {
  title: 'Social Network',
  description: 'My awesome social network project',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <AuthProvider>
          <RealtimeProvider>
            <ToastProvider>
              <Navbar />
              <AppShell>{children}</AppShell>
            </ToastProvider>
          </RealtimeProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
