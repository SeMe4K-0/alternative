import React from 'react'
import { Link, Outlet, useLocation } from 'react-router-dom'
import { useAuthStore } from '../../stores/authStore'
import { useUIStore } from '../../stores/uiStore'
import { SunIcon, MoonIcon } from '@heroicons/react/24/outline'

const Layout: React.FC = () => {
  const { user, logout } = useAuthStore()
  const { sidebarOpen, setSidebarOpen, theme, toggleTheme } = useUIStore()
  const location = useLocation()

  const navigation = [
    { name: 'Наблюдения', href: '/observations', icon: '' },
    { name: 'Расчёты', href: '/calculations', icon: '' },
    { name: 'История', href: '/history', icon: '' },
  ]

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/'
    }
    return location.pathname.startsWith(path)
  }

  return (
    <div className={`h-screen flex ${theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'}`}>
      {/* Mobile sidebar */}
      <div className={`fixed inset-0 z-50 lg:hidden ${sidebarOpen ? 'block' : 'hidden'}`}>
        <div className="fixed inset-0 bg-gray-600 bg-opacity-75" onClick={() => setSidebarOpen(false)} />
        <div className={`relative flex-1 flex flex-col max-w-xs w-full shadow-xl ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
          <div className="absolute top-0 right-0 -mr-12 pt-2">
            <button
              type="button"
              className="ml-1 flex items-center justify-center h-10 w-10 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
              onClick={() => setSidebarOpen(false)}
            >
              <span className="sr-only">Close sidebar</span>
              <span className="text-white text-xl">×</span>
            </button>
          </div>
          <div className="flex-1 h-0 pt-4 pb-4 overflow-y-auto">
               <div className="flex-shrink-0 flex items-center px-4">
                 <Link to="/dashboard" className={`p-3 rounded-xl transition-all duration-200 flex items-center ${theme === 'dark' ? 'hover:bg-gray-700' : 'hover:bg-gray-100'}`}>
                   <img src="/cometTracker.png" alt="Comet Tracker" className="w-auto h-auto" />
                 </Link>
               </div>
            <nav className="mt-3 px-2 space-y-1">
              {navigation.map((item) => (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`group flex items-center px-3 py-3 text-lg font-medium rounded-xl transition-all duration-200 ${
                    isActive(item.href)
                      ? 'bg-gradient-to-r from-blue-500 to-purple-500 text-white shadow-lg'
                      : theme === 'dark' 
                        ? 'text-gray-300 hover:bg-gray-700 hover:text-white'
                        : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
                  }`}
                  onClick={() => setSidebarOpen(false)}
                >
                  <span className="mr-3 text-xl">{item.icon}</span>
                  {item.name}
                </Link>
              ))}
            </nav>
          </div>
          <div className={`flex-shrink-0 flex border-t p-4 ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'}`}>
              <div className="flex items-center w-full">
                <Link 
                  to="/profile" 
                  className="flex-shrink-0 hover:opacity-80 transition-opacity"
                >
                  {user?.avatar ? (
                    <img
                      key={`${user.avatar}`}
                      src={`${user.avatar}?t=${Date.now()}`}
                      alt={user.name || 'User'}
                      className="h-10 w-10 rounded-full object-cover border-2 border-blue-500"
                    />
                  ) : (
                    <div className="h-10 w-10 rounded-full bg-gradient-to-r from-blue-500 to-purple-500 flex items-center justify-center">
                      <span className="text-white font-medium text-sm">
                        {user?.name?.charAt(0).toUpperCase() || 'U'}
                      </span>
                    </div>
                  )}
                </Link>
              <Link 
                to="/profile" 
                className={`ml-3 flex-1 rounded-lg p-2 transition-colors ${
                  theme === 'dark' ? 'hover:bg-gray-700' : 'hover:bg-gray-50'
                }`}
                onClick={() => setSidebarOpen(false)}
              >
                <p className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>{user?.name}</p>
                <p className={`text-xs font-medium ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>{user?.email}</p>
              </Link>
              <button
                onClick={logout}
                className={`ml-2 px-3 py-1 text-sm rounded-lg transition-colors ${
                  theme === 'dark' 
                    ? 'text-red-400 hover:text-red-300 hover:bg-red-900' 
                    : 'text-red-600 hover:text-red-700 hover:bg-red-50'
                }`}
              >
                Выйти
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Desktop sidebar */}
      <div className="hidden lg:flex lg:flex-shrink-0">
        <div className="flex flex-col w-72">
          <div className={`flex flex-col h-screen shadow-xl border-r ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-200'}`}>
            <div className="flex-1 flex flex-col pt-4 pb-4 overflow-y-auto">
               <div className="flex items-center flex-shrink-0 px-6">
                 <Link to="/dashboard" className={`p-4 rounded-2xl transition-all duration-200 flex items-center ${theme === 'dark' ? 'hover:bg-gray-700' : 'hover:bg-gray-100'}`}>
                   <img src="/cometTracker.png" alt="Comet Tracker" className="w-auto h-auto" />
                 </Link>
               </div>
              <nav className="mt-3 flex-1 px-3 space-y-1">
                {navigation.map((item) => (
                  <Link
                    key={item.name}
                    to={item.href}
                    className={`group flex items-center px-4 py-3 text-lg font-medium rounded-xl transition-all duration-200 ${
                      isActive(item.href)
                        ? 'bg-gradient-to-r from-blue-500 to-purple-500 text-white shadow-lg transform scale-105'
                        : theme === 'dark'
                          ? 'text-gray-300 hover:bg-gray-700 hover:text-white hover:transform hover:scale-105'
                          : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900 hover:transform hover:scale-105'
                    }`}
                  >
                    <span className="mr-3 text-xl">{item.icon}</span>
                    {item.name}
                  </Link>
                ))}
              </nav>
            </div>
            <div className={`flex-shrink-0 flex border-t p-4 ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'}`}>
              <div className="flex items-center w-full">
                <Link 
                  to="/profile" 
                  className="flex-shrink-0 hover:opacity-80 transition-opacity"
                  onClick={() => setSidebarOpen(false)}
                >
                  {user?.avatar ? (
                    <img
                      key={`${user.avatar}`}
                      src={`${user.avatar}?t=${Date.now()}`}
                      alt={user.name || 'User'}
                      className="h-10 w-10 rounded-full object-cover border-2 border-blue-500"
                    />
                  ) : (
                    <div className="h-10 w-10 rounded-full bg-gradient-to-r from-blue-500 to-purple-500 flex items-center justify-center">
                      <span className="text-white font-medium text-sm">
                        {user?.name?.charAt(0).toUpperCase() || 'U'}
                      </span>
                    </div>
                  )}
                </Link>
                <Link 
                  to="/profile" 
                  className={`ml-3 flex-1 rounded-lg p-2 transition-colors ${
                    theme === 'dark' ? 'hover:bg-gray-700' : 'hover:bg-gray-50'
                  }`}
                >
                  <p className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>{user?.name}</p>
                  <p className={`text-xs font-medium ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>{user?.email}</p>
                </Link>
                <button
                  onClick={logout}
                  className={`ml-2 px-3 py-1 text-sm rounded-lg transition-colors ${
                    theme === 'dark' 
                      ? 'text-red-400 hover:text-red-300 hover:bg-red-900' 
                      : 'text-red-600 hover:text-red-700 hover:bg-red-50'
                  }`}
                >
                  Выйти
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="flex flex-col flex-1 h-screen">
         {/* Mobile menu button */}
         <div className={`sticky top-0 z-10 lg:hidden pl-1 pt-1 sm:pl-3 sm:pt-3 ${theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'}`}>
           <button
             type="button"
             className={`-ml-0.5 -mt-0.5 h-12 w-12 inline-flex items-center justify-center rounded-md focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500 ${
               theme === 'dark' 
                 ? 'text-gray-400 hover:text-gray-200' 
                 : 'text-gray-500 hover:text-gray-900'
             }`}
             onClick={() => setSidebarOpen(true)}
           >
             <span className="sr-only">Open sidebar</span>
             <span className="text-xl">☰</span>
           </button>
         </div>

         {/* Theme toggle panel */}
         <div className="fixed top-4 right-4 z-50">
           <button
             onClick={toggleTheme}
             className={`p-3 rounded-xl shadow-lg transition-all duration-200 hover:scale-110 ${
               theme === 'dark' 
                 ? 'bg-gray-700 text-yellow-400 hover:bg-gray-600' 
                 : 'bg-white text-gray-600 hover:bg-gray-100'
             }`}
           >
             {theme === 'dark' ? <SunIcon className="h-6 w-6" /> : <MoonIcon className="h-6 w-6" />}
           </button>
         </div>

        {/* Page content */}
        <main className="flex-1 overflow-y-auto">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default Layout