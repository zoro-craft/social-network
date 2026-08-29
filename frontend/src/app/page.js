export default function HomePage() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100">
      <h1 className="text-4xl font-bold mb-4 text-gray-800">Welcome to Social Network</h1>
      <p className="text-gray-600 mb-8">Authentication system is ready!</p>
      
      <div className="flex gap-4">
        <a 
          href="/login" 
          className="px-6 py-2 bg-blue-600 text-white rounded shadow hover:bg-blue-700 transition">
          Go to Login
        </a>
        <a 
          href="/register" 
          className="px-6 py-2 bg-green-600 text-white rounded shadow hover:bg-green-700 transition">
          Create Account
        </a>
      </div>
    </div>
  );
}