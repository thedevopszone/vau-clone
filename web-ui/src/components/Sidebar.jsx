import React from 'react';
import { HiKey, HiUserGroup, HiDocument, HiCog } from 'react-icons/hi';

const Sidebar = () => {
  const menuItems = [
    { icon: HiKey, label: 'Secrets', path: 'secrets', active: true },
    { icon: HiUserGroup, label: 'Auth', path: 'auth', active: false },
    { icon: HiDocument, label: 'Policies', path: 'policies', active: false },
    { icon: HiCog, label: 'Settings', path: 'settings', active: false },
  ];

  return (
    <aside className="w-56 bg-vault-dark-sidebar border-r border-vault-border-dark flex flex-col">
      <nav className="py-2 flex-1">
        {menuItems.map((item, index) => {
          const Icon = item.icon;
          return (
            <div
              key={index}
              className={`
                flex items-center gap-3 px-6 py-3.5 cursor-pointer transition-all border-l-3
                ${item.active
                  ? 'bg-vault-primary/15 text-white border-l-vault-primary font-medium'
                  : item.active === false
                  ? 'text-vault-text-on-dark opacity-40 cursor-not-allowed border-l-transparent'
                  : 'text-vault-text-on-dark border-l-transparent hover:bg-white/5 hover:text-white'
                }
              `}
            >
              <Icon className={`w-5 h-5 ${item.active ? '' : 'opacity-60'}`} />
              <span className="text-sm">{item.label}</span>
            </div>
          );
        })}
      </nav>
    </aside>
  );
};

export default Sidebar;
