import './globals.css';
import { AuthProvider } from '@/app/context/AuthContext';
import { ToastProvider } from '@/app/context/ToastContext';
import { RealtimeProvider } from '@/components/RealtimeProvider';
import Navbar from '@/components/Navbar';

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
              {children}
            </ToastProvider>
          </RealtimeProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
