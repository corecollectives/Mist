import { Outlet } from 'react-router-dom';
import { SidebarProvider, SidebarTrigger } from './components/ui/sidebar';
import { AppSidebar } from './components/app-sidebar';

export const Layout = () => {
  return (
    <SidebarProvider>
      <AppSidebar />
      <main className='w-full flex flex-col gap-y-4 pb-6 px-6 pt-4'>
        <SidebarTrigger />
        <Outlet />
      </main>
    </SidebarProvider>
  );
};
