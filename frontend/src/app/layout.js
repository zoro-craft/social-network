import './globals.css';
import { AuthProvider } from '@/app/context/AuthContext';
import { ToastProvider } from '@/app/context/ToastContext';
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
          <ToastProvider>
            <Navbar />
            {children}
          </ToastProvider>
        </AuthProvider>
      </body>
    </html>
  );
}