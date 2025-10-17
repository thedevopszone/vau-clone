import React from 'react';
import Header from './Header';
import Sidebar from './Sidebar';

const Layout = ({ children, onLogout }) => {
  return (
    <div className="flex flex-col h-screen overflow-hidden">
      <Header onLogout={onLogout} />
      <div className="flex flex-1 overflow-hidden">
        <Sidebar />
        <main className="flex-1 overflow-y-auto bg-vault-light-bg p-8">
          {children}
        </main>
      </div>
    </div>
  );
};

export default Layout;
